package main

import (
	"delay"
	"fmt"
	"rtos"

	"arch/cortexm/bitband"
	"arch/cortexm/debug/itm"

	"stm32/hal/dma"
	"stm32/hal/exti"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/spi"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

const dbg = itm.Port(0)

var (
	nrfirq exti.Lines
	nrfspi *spi.Driver
	nrfcsn bitband.Bit
)

func init() {
	system.Setup96(8)
	systick.Setup()

	gpio.A.EnableClock(true)
	spiport, sck, miso, mosi := gpio.A, gpio.Pin5, gpio.Pin6, gpio.Pin7

	gpio.B.EnableClock(true)
	ctrport, csn, irqn, ce := gpio.B, gpio.Pin6, gpio.Pin8, gpio.Pin9
	nrfcsn = ctrport.OutPins().Bit(6)

	// nRF24 SPI.

	spiport.Setup(sck|mosi, gpio.Config{Mode: gpio.Alt, Speed: gpio.High})
	spiport.Setup(miso, gpio.Config{Mode: gpio.AltIn})
	spiport.SetAltFunc(sck|miso|mosi, gpio.SPI1)
	d := dma.DMA2
	d.EnableClock(true)
	nrfspi = spi.NewDriver(spi.SPI1, d.Channel(2, 3), d.Channel(3, 3))
	nrfspi.P.EnableClock(true)
	nrfspi.P.SetConf(
		spi.Master | spi.MSBF | spi.CPOL0 | spi.CPHA0 |
			nrfspi.P.BR(10e6) | // 10 MHz max.
			spi.SoftSS | spi.ISSHigh,
	)
	nrfspi.P.Enable()
	rtos.IRQ(irq.SPI1).Enable()
	rtos.IRQ(irq.DMA2_Stream2).Enable()
	rtos.IRQ(irq.DMA2_Stream3).Enable()

	// nRF24 control lines.

	ctrport.Setup(csn|ce, gpio.Config{Mode: gpio.Out, Speed: gpio.High})
	ctrport.SetPins(csn)
	ctrport.Setup(irqn, gpio.Config{Mode: gpio.In, Pull: gpio.PullUp})
	nrfirq = exti.Lines(irqn)
	nrfirq.Connect(ctrport)
	nrfirq.EnableFallTrig()
	nrfirq.EnableInt()
	rtos.IRQ(irq.EXTI9_5).Enable()
}

func main() {
	delay.Millisec(500) // For openocd setting SWO.
	fmt.Printf(
		"PCLK: %d Hz, SPI speed: %d Hz\n",
		nrfspi.P.Bus().Clock(), nrfspi.P.Baudrate(nrfspi.P.Conf()),
	)
	for {
		var resp [2]byte
		nrfcsn.Clear()
		n := nrfspi.WriteRead([]byte{5}, resp[:])
		nrfcsn.Set()
		fmt.Printf("n=%d data=%x err=%v\n", n, resp, nrfspi.Err())
		delay.Millisec(100)
	}
}

func exti9_5ISR() {
	lines := exti.Pending() & (exti.L9 | exti.L8 | exti.L7 | exti.L6 | exti.L5)
	lines.ClearPending()
	if lines&nrfirq != 0 {
		dbg.WriteString("nRF24 ISR\n")
	}
}

func nrfSPIISR() {
	nrfspi.ISR()
}

func nrfRxDMAISR() {
	nrfspi.DMAISR(nrfspi.RxDMA)
}

func nrfTxDMAISR() {
	nrfspi.DMAISR(nrfspi.TxDMA)
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.EXTI9_5:      exti9_5ISR,
	irq.SPI1:         nrfSPIISR,
	irq.DMA2_Stream2: nrfRxDMAISR,
	irq.DMA2_Stream3: nrfTxDMAISR,
}
