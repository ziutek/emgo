package main

import (
	"delay"
	"fmt"
	"rtos"

	"arch/cortexm/bitband"
	"arch/cortexm/debug/itm"

	"nrf24"

	"stm32/hal/dma"
	"stm32/hal/exti"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/spi"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

const dbg = itm.Port(0)

type nrfDCI struct {
	spi *spi.Driver
	csn bitband.Bit
	irq exti.Lines
}

func (dc *nrfDCI) WriteRead(oi ...[]byte) (n int, err error) {
	dc.csn.Clear()
	dc.spi.WriteReadMany(oi...)
	dc.csn.Set()
	return n, dc.spi.Err()
}

func (dc *nrfDCI) SetCE(v int) error {
	return nil
}

var nrfdci nrfDCI

func init() {
	system.Setup96(8)
	systick.Setup()

	gpio.A.EnableClock(true)
	spiport, sck, miso, mosi := gpio.A, gpio.Pin5, gpio.Pin6, gpio.Pin7

	gpio.B.EnableClock(true)
	ctrport, csn, irqn, ce := gpio.B, gpio.Pin6, gpio.Pin8, gpio.Pin9
	nrfdci.csn = ctrport.OutPins().Bit(6)

	// nRF24 SPI.

	spiport.Setup(sck|mosi, gpio.Config{Mode: gpio.Alt, Speed: gpio.High})
	spiport.Setup(miso, gpio.Config{Mode: gpio.AltIn})
	spiport.SetAltFunc(sck|miso|mosi, gpio.SPI1)
	d := dma.DMA2
	d.EnableClock(true)
	nrfdci.spi = spi.NewDriver(spi.SPI1, d.Channel(2, 3), d.Channel(3, 3))
	nrfdci.spi.P.EnableClock(true)
	nrfdci.spi.P.SetConf(
		spi.Master | spi.MSBF | spi.CPOL0 | spi.CPHA0 |
			nrfdci.spi.P.BR(10e6) | // 10 MHz max.
			spi.SoftSS | spi.ISSHigh,
	)
	nrfdci.spi.P.Enable()
	rtos.IRQ(irq.SPI1).Enable()
	rtos.IRQ(irq.DMA2_Stream2).Enable()
	rtos.IRQ(irq.DMA2_Stream3).Enable()

	// nRF24 control lines.

	ctrport.Setup(csn|ce, gpio.Config{Mode: gpio.Out, Speed: gpio.High})
	ctrport.SetPins(csn)
	ctrport.Setup(irqn, gpio.Config{Mode: gpio.In, Pull: gpio.PullUp})
	nrfdci.irq = exti.Lines(irqn)
	nrfdci.irq.Connect(ctrport)
	nrfdci.irq.EnableFallTrig()
	nrfdci.irq.EnableInt()
	rtos.IRQ(irq.EXTI9_5).Enable()
}

func main() {
	delay.Millisec(500) // For openocd setting SWO.
	fmt.Printf(
		"\nPCLK: %d Hz, SPI speed: %d Hz\n",
		nrfdci.spi.P.Bus().Clock(), nrfdci.spi.P.Baudrate(nrfdci.spi.P.Conf()),
	)
	nrf := nrf24.Device{DCI: &nrfdci}
	//rtos.SleepUntil(100e6) // nRF24 requires wait at least 100 ms from start.
	for {
		delay.Millisec(500)
		fmt.Printf("\n")
		config := nrf.Config()
		if nrf.Err != nil {
			fmt.Printf("error: %v\n", nrf.Err)
			continue
		}
		fmt.Printf("CONFIG: %v\n", config)
		fmt.Printf("STATUS: %v\n", nrf.Status)
	}
}

func exti9_5ISR() {
	lines := exti.Pending() & (exti.L9 | exti.L8 | exti.L7 | exti.L6 | exti.L5)
	lines.ClearPending()
	if lines&nrfdci.irq != 0 {
		dbg.WriteString("nRF24 ISR\n")
	}
}

func nrfSPIISR() {
	nrfdci.spi.ISR()
}

func nrfRxDMAISR() {
	nrfdci.spi.DMAISR(nrfdci.spi.RxDMA)
}

func nrfTxDMAISR() {
	nrfdci.spi.DMAISR(nrfdci.spi.TxDMA)
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.EXTI9_5:      exti9_5ISR,
	irq.SPI1:         nrfSPIISR,
	irq.DMA2_Stream2: nrfRxDMAISR,
	irq.DMA2_Stream3: nrfTxDMAISR,
}
