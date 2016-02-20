package main

import (
	"delay"
	"math/rand"
	"rtos"
	"sync"

	"stm32/hal/gpio"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var leds *gpio.Port

const (
	Green  = gpio.Pin12
	Orange = gpio.Pin13
	Red    = gpio.Pin14
	Blue   = gpio.Pin15
)

func init() {
	system.Setup168(8)
	systick.Setup()

	gpio.D.EnableClock(false)
	leds = gpio.D

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	leds.Setup(Green|Orange|Red|Blue, &cfg)

	rnd.Seed(uint64(rtos.Nanosec()))
}

var (
	n   int
	m   sync.Mutex
	rnd rand.XorShift64
)

func loop() {
	for {
		m.Lock() // Comment this line.
		switch n & 3 {
		case 0:
			leds.ClearPins(Green)
			leds.SetPins(Orange)
		case 1:
			leds.ClearPins(Orange)
			leds.SetPins(Red)
		case 2:
			leds.ClearPins(Red)
			leds.SetPins(Blue)
		default:
			leds.ClearPins(Blue)
			leds.SetPins(Green)
		}
		// rnd.Uint32 isn't thread safe but don't care.
		delay.Millisec(50 + int(rnd.Uint32()&0xff))
		n++
		m.Unlock() // Comment this line.
	}
}

func main() {
	go loop()
	go loop()
	go loop()
	go loop()
	loop()
}
