package main

import (
	"delay"
	"rtos"
	"time"
	"time/tz"

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

func draw(fb []led.Color, deg int, c led.Color) {
	if deg %= 360; deg < 0 {
		deg += 360
	}
	n := deg * len(fb) / 360
	deg -= 360 * n / len(fb)
	mask := byte(255 * deg * len(fb) / 360)
	fb[n] = fb[n].Blend(c, 255-mask)
	n = (n + 1) % 24
	fb[n] = fb[n].Blend(c, mask)
}

func main() {
	time.Set(
		time.Date(2018, 4, 25, 15, 0, 0, 0, &tz.EuropeWarsaw),
		rtos.Nanosec(),
	)
	strip := make(wsuart.Strip, 24)
	fb := make([]led.Color, len(strip))
	rgb := wsuart.GRB
	for {
		h, m, s := time.Now().Clock()
		s += h%12*3600 + m*60
		draw(fb, s/120, led.RGBA(99, 0, 0, 99))
		s %= 3600
		draw(fb, s/10, led.RGBA(00, 99, 0, 99))
		s %= 60
		draw(fb, s*6, led.RGBA(0, 0, 99, 99))
		for i, c := range fb {
			strip[i] = rgb.Pixel(c)
			fb[i] = 0
		}
		tts.Write(strip.Bytes())
		delay.Millisec(500)
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
