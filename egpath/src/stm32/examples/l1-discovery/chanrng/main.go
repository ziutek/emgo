package main

import (
	"stm32/l1/setup"
	"stm32/stlink"
)

func gen(c chan<- byte, b byte) {
	for {
		c <- b
	}
}

func main() {
	setup.Performance(0)

	c0 := make(chan byte, 2)
	c1 := make(chan byte, 2)

	go gen(c0, '0')
	go gen(c1, '1')

	for {
		var b byte
		select {
		case b = <-c0:
		case b = <-c1:
		}
		stlink.Term.WriteByte(b)
	}
}
