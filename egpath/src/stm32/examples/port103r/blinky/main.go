package main

import (
	"delay"
	"fmt"

	"stm32/hal/setup"
)

const ()

func init() {
	setup.Performance(8, 72/8, false)
}

func wait() {
	//delay.Loop(1e7)
	delay.Millisec(100)
}

func main() {
	for i := 0; ; i++ {
		wait()
		fmt.Println(i, "Debug!")
	}
}
