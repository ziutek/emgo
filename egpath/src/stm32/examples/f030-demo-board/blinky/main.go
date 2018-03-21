package main

import (
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

func init() {
	system.SetupPLL(8, 1, 48/8)
	systick.Setup(2e6)
}

func main() {
	for {

	}
}
