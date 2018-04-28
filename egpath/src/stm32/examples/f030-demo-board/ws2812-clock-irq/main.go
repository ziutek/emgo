package main

import (
	"delay"
	"rtos"

	"led"
	"led/ws281x/wsuart"

	"stm32/hal/dma"
	"stm32/hal/exti"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
	"stm32/hal/usart"
)

var (
	tts      *usart.Driver
	btn      gpio.Pin
	btnEvent rtos.EventFlag
)

func init() {
	system.SetupPLL(8, 1, 48/8)
	systick.Setup(2e6)

	gpio.A.EnableClock(true)
	btn = gpio.A.Pin(4)
	tx := gpio.A.Pin(9)

	btn.Setup(&gpio.Config{Mode: gpio.In, Pull: gpio.PullUp})
	ei := exti.Lines(btn.Mask())
	ei.Connect(btn.Port())
	ei.EnableFallTrig()
	ei.EnableRisiTrig()
	ei.EnableIRQ()
	rtos.IRQ(irq.EXTI4_15).Enable()

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

		hs %= 3600 * 2 // Half-seconds since the beginning of the current hour.
		m := len(strip) * hs / (3600 * 2)

		hs %= 60 * 2 // Half-seconds since the beginning of the current minute.
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
		if btnWait(0, 250) {
			setClock += setSpeed
			delay.Millisec(100)
			if !btnWait(1, 100) && setSpeed < 10*60*2 {
				setSpeed += 10
			}
			continue
		}
		setSpeed = 5
	}
}

func btnWait(state, ms int) bool {
	deadline := rtos.Nanosec() + int64(ms)*1e6
	for btn.Load() != state {
		if !btnEvent.Wait(1, deadline) {
			return false // Timeout
		}
		btnEvent.Reset(0)
	}
	return true
}

func exti4_15ISR() {
	pending := exti.Pending()
	pending &= exti.L4 | exti.L5 | exti.L6 | exti.L7 | exti.L8 | exti.L9 |
		exti.L10 | exti.L11 | exti.L12 | exti.L13 | exti.L14 | exti.L15
	pending.ClearPending()
	if pending&exti.Lines(btn.Mask()) != 0 {
		btnEvent.Signal(1)
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
	irq.EXTI4_15:        exti4_15ISR,
}
