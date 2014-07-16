package main

import (
	"math/rand"
	"strconv"

	"stm32/l1/setup"
	"stm32/stlink"
)

var st = stlink.Term

func main() {
	setup.Performance(0)

	var (
		buf [20]byte
		rnd rand.XorShift64
	)

	rnd.Seed(1)

	for {
		strconv.Utoa(buf[:], rnd.Uint32(), 10)
		st.Write(buf[:])
		st.WriteByte('\n')
	}
}
