package internal

import (
	"mmio"
)

type AtomicBit struct {
	R *mmio.U32
	N int
}

func (b AtomicBit) Load() int {
	return int(b.R.Load()>>uint(b.N)) & 1
}

func (b AtomicBit) Store(v int) {
	b.R.Store(uint32(v&1) << uint(b.N))
}

func (b AtomicBit) Set() {
	b.R.AtomicSetBits(1 << uint(b.N))
}

func (b AtomicBit) Clear() {
	b.R.AtomicClearBits(1 << uint(b.N))
}
