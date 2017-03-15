package main

import (
	"bufio"
	"debug/semihosting"
	"fmt"

	"nrf5/hal/clock"
	"nrf5/hal/gpio"
	"nrf5/hal/irq"
	"nrf5/hal/rtc"
	"nrf5/hal/system"
	"nrf5/hal/system/timer/rtcst"
)

var leds [5]gpio.Pin

func init() {
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC0, 1)

	for i := range leds {
		led := gpio.P0.Pin(18 + i)
		led.Setup(&gpio.Config{Mode: gpio.Out})
		leds[i] = led
	}
}

func checkErr(err error) {
	for err != nil {
	}
}

func main() {
	leds[0].Set()

	f, err := semihosting.OpenFile(":tt", semihosting.W)
	checkErr(err)
	w := bufio.NewWriterSize(f, 80)

	leds[1].Set()

	for i := 0; ; i++ {
		_, err := fmt.Fprintf(w, "%d: Hello world!\n", i)
		checkErr(err)
		checkErr(w.Flush())
		leds[2].Store(i & 1)
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC0: rtcst.ISR,
}
