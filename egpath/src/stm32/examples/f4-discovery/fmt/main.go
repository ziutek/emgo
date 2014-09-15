package main

import (
	"fmt"
	"stm32/f4/setup"
)

func init() {
	setup.Performance168(8)
	initConsole()
}

func main() {
	con.WriteString("\nKasia ma kota!\n")
	fmt.Fprint(con, fmt.Int64(15))
}
