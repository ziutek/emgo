package main

import (
	"strconv"

	"stm32/l1/setup"
	"stm32/stlink"
)

var st = stlink.Term

func main() {
	setup.Performance(0)

	st.WriteString("Press any key (Ctrl-C to exit)\n")

	var buf [stlink.TermBufLen]byte
	num := buf[:7]

	for i := 0; ; i++ {
		n := strconv.Itoa(num, int32(i), 10)
		st.Write(num[n:])
		st.WriteString(": ")

		n, _ = st.Read(buf[:])

		st.Write(buf[:n])
		st.WriteByte('\n')
	}
}
