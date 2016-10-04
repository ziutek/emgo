package main

import (
	"delay"
	"fmt"
	"rtos"

	"arch/cortexm/bitband"
	"arch/cortexm/debug/itm"

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
	nrfspi = spi.NewDriver(spi.SPI1, nil, nil)
	nrfspi.EnableClock(true)
	nrfspi.SetConf(
		spi.Master | spi.MSBF | spi.CPOL0 | spi.CPHA0 |
			nrfspi.BR(10e6) | // 10 MHz max.
			spi.SoftSS | spi.ISSHigh,
	)
	nrfspi.Enable()

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

func wait(e spi.Event) bool {
	for {
		ev, err := nrfspi.Status()
		if err != 0 {
			fmt.Printf("Error: %v\n", err)
			return false
		}
		if ev&e != 0 {
			return true
		}
	}
}

func main() {
	for {
		delay.Millisec(1e3)
		nrfcsn.Clear()
		nrfspi.StoreByte(5)
		if !wait(spi.TxEmpty) {
			break
		}
		nrfspi.StoreByte(0)
		if !wait(spi.RxNotEmpty) {
			break
		}
		status := nrfspi.LoadByte()
		if !wait(spi.RxNotEmpty) {
			break
		}
		ch := nrfspi.LoadByte()
		nrfcsn.Set()
		fmt.Printf("status=%02x ch=%d\n", status, ch)
	}
}

func exti9_5ISR() {
	lines := exti.Pending() & (exti.L9 | exti.L8 | exti.L7 | exti.L6 | exti.L5)
	lines.ClearPending()
	if lines&nrfirq != 0 {
		dbg.WriteString("nRF24 ISR\n")
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.EXTI9_5: exti9_5ISR,
}
