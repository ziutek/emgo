package main

import (
	"math/rand"
	"rtos"
	"strconv"

	"stm32/l1/setup"
)

var dbg = rtos.Debug(0)

func main() {
	setup.Performance(0)

	var (
		buf [20]byte
		rnd rand.XorShift64
	)

	rnd.Seed(1)

	for {
		strconv.Utoa(buf[:], uint(rnd.Uint32()), 10)
		dbg.Write(buf[:])
		dbg.WriteByte('\n')
	}
}
