package main

import "stm32/stlink"

func main() {
	for {
		stlink.Term.WriteString("Hello world!\n")
	}
}
