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

	var rnd rand.XorShift64

	rnd.Seed(1)

	for {
		strconv.WriteUint32(dbg, rnd.Uint32(), 10, -10)
		dbg.WriteByte('\n')
	}
}
