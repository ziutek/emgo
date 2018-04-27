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

var (
	tts *usart.Driver
	btn gpio.Pin
)

func init() {
	system.SetupPLL(8, 1, 48/8)
	systick.Setup(2e6)

	gpio.A.EnableClock(true)
	btn = gpio.A.Pin(4)
	tx := gpio.A.Pin(9)

	btn.Setup(&gpio.Config{Mode: gpio.In, Pull: gpio.PullUp})

	tx.Setup(&gpio.Config{Mode: gpio.Alt})
	tx.SetAltFunc(gpio.USART1_AF1)
	d := dma.DMA1
	d.EnableClock(true)

	// 1390 ns/WS2812bit = 3 * 463 ns/UARTbit -> BR = 3 * 1e9 ns/s / 1390 ns/bit

	tts = usart.NewDriver(usart.USART1, d.Channel(2, 0), nil, nil)
	tts.Periph().EnableClock(true)
	tts.Periph().SetBaudRate(3000000000 / 1390)
	tts.Periph().SetConf2(usart.TxInv) // STM32F0 need no external inverter.
	tts.Periph().Enable()
	tts.EnableTx()

	rtos.IRQ(irq.USART1).Enable()
	rtos.IRQ(irq.DMA1_Channel2_3).Enable()
}

func main() {
	var setClock, setSpeed int
	rgb := wsuart.GRB
	strip := make(wsuart.Strip, 24)
	for {
		hs := int(rtos.Nanosec() / 5e8) // Half-seconds elapsed since reset.
		hs += setClock

		hs %= 12 * 3600 * 2 // Half-seconds since the last 0:00 or 12:00.
		h := len(strip) * hs / (12 * 3600 * 2)

		hs %= 3600 * 2 // Half-second since the beginning of the current hour.
		m := len(strip) * hs / (3600 * 2)

		hs %= 60 * 2 // Half-second since the beginning of the current minute.
		s := len(strip) * hs / (60 * 2)

		hc := led.Color(0x550000)
		mc := led.Color(0x005500)
		sc := led.Color(0x000055)

		// Blend colors if the hands of the clock overlap.
		if h == m {
			hc |= mc
			mc = hc
		}
		if m == s {
			mc |= sc
			sc = mc
		}
		if s == h {
			sc |= hc
			hc = sc
		}

		// Draw the clock and send to the ring.
		strip.Clear()
		strip[h] = rgb.Pixel(hc)
		strip[m] = rgb.Pixel(mc)
		strip[s] = rgb.Pixel(sc)
		tts.Write(strip.Bytes())

		// Adjust the clock.
		if btn.Load() == 0 {
			setClock += setSpeed
			i, n := 0, 10
			for btn.Load() == 0 && i < n {
				delay.Millisec(20)
				i++
			}
			if i == n && setSpeed < 10*60*2 {
				setSpeed += 10
			}
			continue
		}
		setSpeed = 5
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
