package flash

import "unsafe"

type regs struct {
	acr      uint32 `C:"volatile"`
	keyr     uint32 `C:"volatile"`
	optkeyr  uint32 `C:"volatile"`
	sr       uint32 `C:"volatile"`
	cr       uint32 `C:"volatile"`
	ar       uint32 `C:"volatile"`
	reserved uint32 `C:"volatile"`
	obr      uint32 `C:"volatile"`
	wrpr     uint32 `C:"volatile"`
}

func (r *regs) acrSetBit(n uint, b bool) {
	if b {
		r.acr |= 1 << n
	} else {
		r.acr &^= 1 << n
	}
}

const base uintptr = 0x40023c00

var f = (*regs)(unsafe.Pointer(base))

func Prefetch() bool {
	return f.acr&(1<<8) != 0
}

func SetPrefetch(enable bool) {
	f.acrSetBit(8, enable)
}

func ICache() bool {
	return f.acr&(1<<9) != 0
}

func SetICache(enable bool) {
	f.acrSetBit(9, enable)
}

func DCache() bool {
	return f.acr&(1<<10) != 0
}

func SetDCache(enable bool) {
	f.acrSetBit(10, enable)
}

func Latency() int {
	return int(f.acr & 0x7)
}

func SetLatency(waitStates int) {
	f.acr = f.acr&^0x7 | uint32(waitStates)&0x7
}