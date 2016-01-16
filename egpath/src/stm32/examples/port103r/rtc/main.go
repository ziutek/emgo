package main

import (
	"delay"
	"fmt"

	"stm32/hal/setup"
)

func init() {
	setup.Performance(8, 72/8, false)
	setup.UseSysTick()
	setup.UseRTC(32768)
}

func main() {
	delay.Millisec(100)
	fmt.Println("Start")

	buf := make([]int64, 256)
	var prev int64
	for i := 0; i < len(buf); {
		now := setup.RtcNow()
		if now != prev {
			buf[i] = now
			prev = now
			i++
		}
	}
	for i, v := range buf {
		if i == 0 {
			fmt.Printf("%d: %d\n", i, v)
		} else {
			fmt.Printf("%3d: %d %d\n", i, v, v-buf[i-1])
		}
		delay.Millisec(5)
	}
}
