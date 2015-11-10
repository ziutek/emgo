package systick

import "unsafe"

type systick struct {
	csr   uint32 // Control and Status Register
	rvr   uint32 // Reload Value Register
	cvr   uint32 // Current Value Register
	calib uint32 // Calibration Value Register
} //c:volatile

var stk = (*systick)(unsafe.Pointer(uintptr(0xe000e010)))

// Implemented as set/clear flags instead of enable/disable functions because
// any read of ctrl register clears Count flag.

type Flag uint32

const (
	Enable  Flag = 1 << 0  // Counter ebabled.
	TickInt Flag = 1 << 1  // Generate exceptions.
	ClkCPU  Flag = 1 << 2  // Use CPU clock as clock source.
	Count   Flag = 1 << 16 // Timer counted to 0 since last flag read.
)

// Flags returns SysTick status and control flags.
func Flags() Flag {
	return Flag(stk.csr)
}

// WriteFlags writes f to SysTick status and control register.
func WriteFlags(f Flag) {
	stk.csr = uint32(f)
}

// SetFlags sets all flags specified by f.
func SetFlags(f Flag) {
	stk.csr |= uint32(f)
}

// ClearFlags resets flags specified by f.
func ClearFlags(f Flag) {
	stk.csr &^= uint32(f)
}

// SetReload sets SysTick reload value register. v can be in the range
// [0, 0x00ffffff].
func SetReload(v uint32) {
	stk.rvr = v
}

// Reload returns value of RVR register.
func Reload() uint32 {
	return stk.rvr & 0x00ffffff
}

// Val returns current value of SysTick counter
func Val() uint32 {
	return stk.cvr & 0x00ffffff
}

// Reset counter to 0.
func Reset() {
	stk.cvr = 0
}

// Calib returns calibration properties.
func Calib() (skew, noref bool, tenms uint32) {
	c := stk.calib
	skew = c&(1<<31) != 0
	noref = c&(1<<30) != 0
	tenms = c & 0x00ffffff
	return
}
