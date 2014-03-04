package main

import (
	_ "cortexm/startup"

	"stm32/l1/setup"
	"stm32/stlink"
)

func main() {
	setup.Performance(0)

	var buf [40]byte
	for {
		n, _ := stlink.Con.Read(buf[:])
		stlink.Con.Write(buf[:n])
	}
}
