package main

import (
	"debug/semihosting"
	"io"

	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

func init() {
	system.SetupPLL(8, 1, 48/8)
	systick.Setup(2e6)
}

func checkErr(err error) {
	for err != nil {
	}
}

func main() {
	f, err := semihosting.OpenFile(":tt", semihosting.W)
	checkErr(err)
	io.WriteString(f, "Hello, World!\n")
}
