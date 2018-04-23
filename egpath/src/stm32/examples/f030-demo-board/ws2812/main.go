package main

import (
	"delay"
	"rtos"

	"led"
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
	system.SetupPLL(8, 1, 48/8)
	systick.Setup(2e6)

	gpio.A.EnableClock(true)
	tx := gpio.A.Pin(9)

	tx.Setup(&gpio.Config{Mode: gpio.Alt})
	tx.SetAltFunc(gpio.USART1_AF1)
	d := dma.DMA1
	d.EnableClock(true)

	// WS2812 bit should take 1390 ns -> 463 ns for UART bit -> 2158273 bit/s.

	tts = usart.NewDriver(usart.USART1, d.Channel(2, 0), nil, nil)
	tts.Periph().EnableClock(true)
	tts.Periph().SetBaudRate(3000000000 / 1390)
	//tts.Periph().SetConf1(usart.Word7b) // Not-supported by STM32F0.
	tts.Periph().SetConf2(usart.TxInv) // STM32F0 need no external inverter.
	tts.Periph().Enable()
	tts.EnableTx()

	rtos.IRQ(irq.USART1).Enable()
	rtos.IRQ(irq.DMA1_Channel2_3).Enable()
}

func main() {
	rgb := wsuart.GRB
	strip := make(wsuart.Strip, 24)
	for i := uint32(0); ; i++ {
		strip.Clear()
		h := (i * 2 / 3) % 24
		m := (i / 2) % 24
		s := (1<<32 - 1 - i*2/5) % 24
		strip[h] = rgb.Pixel(led.RGB(99, 0, 0))
		if h != m {
			strip[m] = rgb.Pixel(led.RGB(0, 0, 99))
		} else {
			strip[m] = rgb.Pixel(led.RGB(99, 0, 99))
		}
		if s != m && s != h {
			strip[s] = rgb.Pixel(led.RGB(0, 99, 0))
		} else if s == m && s == h {
			strip[s] = rgb.Pixel(led.RGB(99, 99, 99))
		} else if s == m {
			strip[s] = rgb.Pixel(led.RGB(0, 99, 99))
		} else {
			strip[s] = rgb.Pixel(led.RGB(99, 99, 0))
		}
		tts.Write(strip.Bytes())
		delay.Millisec(50)
	}
}

func ttsISR() {
	tts.ISR()
}

func ttsDMAISR() {
	tts.TxDMAISR()
}

//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.USART1:          ttsISR,
	irq.DMA1_Channel2_3: ttsDMAISR,
}
