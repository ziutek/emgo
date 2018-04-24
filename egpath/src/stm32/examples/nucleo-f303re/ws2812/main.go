package main

import (
	"delay"
	"math/rand"
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
	system.SetupPLL(8, 1, 72/8)
	systick.Setup(2e6)

	// GPIO

	gpio.C.EnableClock(true)
	tx := gpio.C.Pin(10)

	tx.Setup(&gpio.Config{Mode: gpio.Alt})
	tx.SetAltFunc(gpio.UART4)
	d := dma.DMA2
	d.EnableClock(true)

	// WS2812 bit should take 1390 ns -> 463 ns for UART bit -> 2158273 bit/s.

	tts = usart.NewDriver(usart.UART4, d.Channel(5, 0), nil, nil)
	tts.Periph().EnableClock(true)
	tts.Periph().SetBaudRate(3000000000 / 1390)
	tts.Periph().SetConf1(usart.Word7b) // 9 bits: 1 start, 7 data, 1 stop.
	tts.Periph().SetConf2(usart.TxInv)  // STM32F3 need no external inverter.
	tts.Periph().Enable()
	tts.EnableTx()

	rtos.IRQ(irq.UART4).Enable()
	rtos.IRQ(irq.DMA2_Channel5).Enable()
}

func main() {
	var rnd rand.XorShift64
	rnd.Seed(1)
	strip := make(wsuart.Strip, 24)
	rgb := wsuart.GRB
	for k := 0; ; k++ {
		c := led.Color(rnd.Uint32())
		for i := range strip {
			strip[(i+k)%24] = rgb.Pixel(c.Mask(byte(255 * (i + 1) / 24)))
		}
		tts.Write(strip.Bytes())
		delay.Millisec(1000)
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
	irq.UART4:         ttsISR,
	irq.DMA2_Channel5: ttsTxDMAISR,
}
