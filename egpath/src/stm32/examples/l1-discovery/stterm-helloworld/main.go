package main

import (
	"stm32/l1/setup"
	"stm32/stlink"
)

var st = stlink.Term

func main() {
	setup.Performance(0)
	st.WriteString("Hello world!\n")
}
