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
	tts   *usart.Driver
	btn   gpio.Pin
	btnev rtos.EventFlag
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
	ei.EnableRiseTrig()
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
	rgb := wsuart.GRB
	strip := wsuart.Make(24)
	ds := 4 * 60 / len(strip) // Interval between LEDs (quarter-seconds).
	adjust := 0
	adjspeed := ds
	for {
		qs := int(rtos.Nanosec() / 25e7) // Quarter-seconds elapsed since reset.
		qa := qs + adjust

		qa %= 12 * 3600 * 4 // Quarter-seconds since the last 0:00 or 12:00.
		hi := len(strip) * qa / (12 * 3600 * 4)

		qa %= 3600 * 4 // Quarter-seconds in the current hour.
		mi := len(strip) * qa / (3600 * 4)

		qa %= 60 * 4 // Quarter-seconds in the current minute.
		si := len(strip) * qa / (60 * 4)

		hc := led.Color(0x550000)
		mc := led.Color(0x005500)
		sc := led.Color(0x000055)

		// Blend the colors if the hands of the clock overlap.
		if hi == mi {
			hc |= mc
			mc = hc
		}
		if mi == si {
			mc |= sc
			sc = mc
		}
		if si == hi {
			sc |= hc
			hc = sc
		}

		// Draw the clock and write to the ring.
		strip.Clear()
		strip[hi] = rgb.Pixel(hc)
		strip[mi] = rgb.Pixel(mc)
		strip[si] = rgb.Pixel(sc)
		tts.Write(strip.Bytes())

		// Sleep until the button pressed or the second hand should be moved.
		if btnWait(0, int64(qs+ds)*25e7) {
			adjust += adjspeed
			// Sleep until the button is released or timeout.
			if !btnWait(1, rtos.Nanosec()+100e6) {
				if adjspeed < 5*60*4 {
					adjspeed += 2 * ds
				}
				continue
			}
			adjspeed = ds
		}
	}
}

func btnWait(state int, deadline int64) bool {
	for btn.Load() != state {
		if !btnev.Wait(1, deadline) {
			return false // timeout
		}
		btnev.Reset(0)
	}
	delay.Millisec(50) // debouncing
	return true
}

func exti4_15ISR() {
	pending := exti.Pending() & 0xFFF0
	pending.ClearPending()
	if pending&exti.Lines(btn.Mask()) != 0 {
		btnev.Signal(1)
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
