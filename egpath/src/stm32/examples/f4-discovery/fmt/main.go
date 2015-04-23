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
	con.WriteString("fmt test:\n")

	fmt.Fprint(con, true, " ", false, "\n")
	fmt.Fprint(con, 10, " ", -10, " ", 1234567890, " ", -1234567890, "\n")
	fmt.Fprint(con, int64(1234567890123), " ", int64(-1234567890123), "\n")
	fmt.Fprint(con, 123.456e-20, " ", -123.456e20, "\n")
}
