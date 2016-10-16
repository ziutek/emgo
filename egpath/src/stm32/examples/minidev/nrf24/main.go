package main

import (
	"delay"
	"rtos"

	"arch/cortexm/bitband"

	"stm32/hal/dma"
	"stm32/hal/exti"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/spi"
	"stm32/hal/system"
	"stm32/hal/system/timer/rtc"

	"stm32/hal/raw/rcc"
	//"stm32/hal/raw/tim"
)

var led bitband.Bit

func init() {
	system.Setup(8, 72/8, false)
	rtc.Setup(32768)

	gpio.A.EnableClock(true)
	spiport, csn, sck, miso, mosi := gpio.A, gpio.Pin4, gpio.Pin5, gpio.Pin6, gpio.Pin7

	gpio.B.EnableClock(true)
	ctrport, ce, irqn := gpio.B, gpio.Pin0, gpio.Pin1

	gpio.C.EnableClock(false)
	ledport, ledpin := gpio.C, 13

	// LED

	ledport.SetupPin(ledpin, gpio.Config{Mode: gpio.Out, Speed: gpio.Low})
	led = ledport.OutPins().Bit(ledpin)

	// nRF24 SPI

	spiport.Setup(sck|mosi, gpio.Config{Mode: gpio.Alt, Speed: gpio.High})
	spiport.Setup(miso, gpio.Config{Mode: gpio.AltIn})
	ctrport.Setup(csn, gpio.Config{Mode: gpio.Out, Speed: gpio.High})
	d := dma.DMA1
	d.EnableClock(true)
	spid := spi.NewDriver(spi.SPI1, d.Channel(2, 0), d.Channel(3, 0))
	spid.P.EnableClock(true)
	rtos.IRQ(irq.SPI1).Enable()
	rtos.IRQ(irq.DMA1_Channel2).Enable()
	rtos.IRQ(irq.DMA1_Channel3).Enable()

	// nRF24 control lines.

	ctrport.Setup(ce, gpio.Config{Mode: gpio.Alt, Speed: gpio.High})
	rcc.RCC.TIM3EN().Set()
	ctrport.SetupPin(irqn, gpio.Config{Mode: gpio.In, Pull: gpio.PullUp})
	irqline := exti.Lines(irqn)
	irqline.Connect(ctrport)
	rtos.IRQ(irq.EXTI1).Enable()

	//dci = nrfdci.NewNRFDCI(spid, spiport, csn, system.APB?.Clock(), tim.TIM3 3, irqline)
}

func main() {
	for {
		delay.Millisec(100)
		led.Set()
		delay.Millisec(100)
		led.Clear()
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTCAlarm: rtc.ISR,
}
