package main

import (
	"fmt"
	"time"

	"stm32/hal/setup"
)

func init() {
	setup.Performance168(8)
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
