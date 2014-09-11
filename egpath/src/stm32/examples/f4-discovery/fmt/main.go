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
	con.WriteString("\n\nStart:\n")

	fmt.Int32(17).Format(con, 10)
	con.WriteByte('\n')
}
