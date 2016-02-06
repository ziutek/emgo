package main

import (
	"fmt"
	"time"

	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

func init() {
	system.Setup168(8)
	systick.Setup()
	initConsole()
}

func main() {
	t := time.Now()
	fmt.Println(t)
	fmt.Println(true, false)
	fmt.Println(10, -10, 1234567890, -123456789)
	fmt.Println(int64(1234567890123), int64(-1234567890123))
	fmt.Println(123.456e-20, -123.456e2)
}
