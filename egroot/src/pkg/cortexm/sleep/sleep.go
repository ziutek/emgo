package sleep

import (
	"mmio"
	"unsafe"
)

var scr = (*mmio.Reg32)(unsafe.Pointer(uintptr(0xe000ed10)))

func EventOnPend() bool {
	return scr.Bit(4)
}

func EnableEventOnPend() {
	scr.SetBit(4)
}

func DisableEventOnPend() {
	scr.ClearBit(4)
}
func SleepDeep() bool {
	return scr.Bit(2)
}

func EnableSleepDeep() {
	scr.SetBit(2)
}

func DisableSleepDeep() {
	scr.ClearBit(2)
}
func SleepOnExit() bool {
	return scr.Bit(1)
}
func EnableSleepOnExit() {
	scr.SetBit(1)
}

func DisableSleepOnExit() {
	scr.SetBit(1)
}

func WFE()

func WFI()
