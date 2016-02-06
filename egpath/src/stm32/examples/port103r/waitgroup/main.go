package main

import (
	"delay"
	"math/rand"
	"rtos"
	"sync"

	"arch/cortexm/bitband"

	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/rtc"
)

var (
	leds [4]bitband.Bit
	key3 bitband.Bit
	rnd  rand.XorShift64
)

func init() {
	system.Setup(8, 72/8, false)
	rtc.Setup(32768)

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	port := gpio.B
	port.EnableClock(false)
	for n, pin := range []int{7, 6, 5} {
		port.SetupPin(pin, &cfg)
		leds[n] = port.OutPin(pin)
	}
	port = gpio.D
	port.EnableClock(false)
	port.SetupPin(2, &cfg)
	leds[3] = port.OutPin(2)

	cfg = gpio.Config{Mode: gpio.In, Pull: gpio.PullUp}
	port = gpio.C
	port.EnableClock(true)
	port.SetupPin(10, &cfg)
	key3 = port.InPin(10)

	rnd.Seed(uint64(rtos.Nanosec()))
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTCAlarm: rtc.ISR,
}

func waitkey(wg *sync.WaitGroup) {
	for key3.Load() != 0 {
		delay.Millisec(5)
	}
	wg.Add(-1)
}

func setled(wg *sync.WaitGroup, led bitband.Bit, v int) {
	// rnd.Uint32 isn't thread safe but don't care.
	delay.Millisec(100 + int(rnd.Uint32()&0x3ff))
	led.Store(v)
	delay.Millisec(100 + int(rnd.Uint32()&0x3ff))
	wg.Add(-1)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go waitkey(&wg)
	wg.Wait()
	for _, led := range leds[:3] {
		wg.Add(1)
		go setled(&wg, led, 1)
	}
	wg.Wait()
	for _, led := range leds[:3] {
		wg.Add(1)
		go setled(&wg, led, 0)
	}
	wg.Wait()
	leds[3].Set()
}
