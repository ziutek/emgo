package flash

import "unsafe"

type regs struct {
	acr      uint32
	keyr     uint32
	optkeyr  uint32
	sr       uint32
	cr       uint32
	ar       uint32
	reserved uint32
	obr      uint32
	wrpr     uint32
} //C:volatile

func (r *regs) acrSetBit(n uint, b bool) {
	if b {
		r.acr |= 1 << n
	} else {
		r.acr &^= 1 << n
	}
}

var f = (*regs)(unsafe.Pointer(uintptr(0x40023c00)))

func ResetACR() {
	f.acr = 0
}

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
	return int(f.acr & 0xf)
}

// SetLatency sets number of waitStates used for
// Flash memory access. Allowed values:
// 0-7  for STM32F405xx/07xx, STM32F415xx/17xx,
// 0-15 for STM32F42xxx and STM32F43xxx.
func SetLatency(waitStates int) {
	f.acr = f.acr&^0xf | uint32(waitStates)&0xf
}
