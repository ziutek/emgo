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

var LED *gpio.Port

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
	LED = gpio.D

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	LED.Setup(Green|Orange|Red|Blue, &cfg)

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
			LED.ClearPins(Green)
			LED.SetPins(Orange)
		case 1:
			LED.ClearPins(Orange)
			LED.SetPins(Red)
		case 2:
			LED.ClearPins(Red)
			LED.SetPins(Blue)
		default:
			LED.ClearPins(Blue)
			LED.SetPins(Green)
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
