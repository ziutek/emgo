package flash

import "unsafe"

type regs struct {
	acr uint32
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

// SetLatency sets number of waitStates used for
// Flash memory access. Allowed values: 0, 1.
func Latency() int {
	return int(f.acr & 1)
}

func SetLatency(waitStates int) {
	f.acr = f.acr&^1 | uint32(waitStates)&1
}

func Prefetch() bool {
	return f.acr&(1<<1) != 0
}

func SetPrefetch(enable bool) {
	f.acrSetBit(1, enable)
}

func Acc64() bool {
	return f.acr&(1<<2) != 0
}

func SetAcc64(enable bool) {
	f.acrSetBit(2, enable)
}

func SleepPD() bool {
	return f.acr&(1<<3) != 0
}

func SetSleepPD(powerDown bool) {
	f.acrSetBit(3, powerDown)
}

func RunPD() bool {
	return f.acr&(1<<4) != 0
}

func SetRunPD(powerDown bool) {
	f.acrSetBit(4, powerDown)
}
