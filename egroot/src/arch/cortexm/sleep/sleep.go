package sleep

import "mmio"

var scr = mmio.PtrReg32(0xe000ed10)

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

//c:static inline
func WFE()

//c:static inline
func WFI()
