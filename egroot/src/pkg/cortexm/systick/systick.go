package systick

import "unsafe"

type systick struct {
	ctrl  uint32 `C:"volatile"`
	load  uint32 `C:"volatile"`
	val   uint32 `C:"volatile"`
	calib uint32 `C:"volatile"`
}

const base uintptr = 0xe000e010

var stk = (*systick)(unsafe.Pointer(base))

type Flag uint32

const (
	Enable  Flag = 1 << 0  // set if SysTick counter is enabled
	TickInt Flag = 1 << 1  // set if timer asserts SysTick exception
	ClkFast Flag = 1 << 2  // set: AHB (CPU clock), unset: AHB/8
	Count   Flag = 1 << 16 // set if timer counted to 0 since last flag read
)

// Flags returns SysTick status and control flags
func Flags() Flag {
	return Flag(stk.ctrl)
}

// SetFlags sets all flags specified by f
func SetFlags(f Flag) {
	stk.ctrl |= uint32(f)
}

// ResetFlags resets flags specified by f
func ResetFlags(f Flag) {
	stk.ctrl &^= uint32(f)
}

// SetLoad sets SysTick LOAD register (v can be in the range 0 - 0x00ffffff)
func SetLoad(v uint32) {
	stk.load = v
}

// Val returns current value of SysTick counter
func Val() uint32 {
	return stk.val & 0x00ffffff
}
