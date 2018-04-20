package main

import (
	"delay"
	"rtos"

	"led/ws281x/wsuart"

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
	"stm32/hal/usart"
)

var tts *usart.Driver

func init() {
	system.SetupPLL(8, 1, 72/8)
	systick.Setup(2e6)

	// GPIO

	gpio.A.EnableClock(true)
	tx := gpio.A.Pin(9)

	tx.Setup(&gpio.Config{Mode: gpio.Alt})
	tx.SetAltFunc(gpio.USART1)
	d := dma.DMA1
	d.EnableClock(true)

	// WS2812 bit should take 1390 ns -> 463 ns for UART bit -> 2158273 bit/s.

	tts = usart.NewDriver(usart.USART1, d.Channel(4, 0), nil, nil)
	tts.Periph().EnableClock(true)
	tts.Periph().SetBaudRate(3000000000 / 1390)
	tts.Periph().SetConf1(usart.Word7b)
	tts.Periph().SetConf2(usart.TxInv) // STM32F3 need no external inverter.
	tts.Periph().Enable()
	tts.EnableTx()

	rtos.IRQ(irq.USART1).Enable()
	rtos.IRQ(irq.DMA1_Channel4).Enable()
}

func main() {
	rgb := wsuart.GRB
	strip := wsuart.Strip{
		rgb.Pixel(0, 0, 99),
		rgb.Pixel(0, 99, 0),
		rgb.Pixel(0, 99, 99),
		rgb.Pixel(99, 0, 0),
		rgb.Pixel(99, 0, 99),
		rgb.Pixel(99, 99, 0),
		rgb.Pixel(99, 99, 99),
	}
	for {
		tts.Write(strip.Bytes())
		delay.Millisec(1e3)
		p := strip[0]
		copy(strip, strip[1:])
		strip[len(strip)-1] = p
	}
}

func ttsISR() {
	tts.ISR()
}

func ttsTxDMAISR() {
	tts.TxDMAISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.USART1:        ttsISR,
	irq.DMA1_Channel4: ttsTxDMAISR,
}
