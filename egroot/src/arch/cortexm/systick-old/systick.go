package systick

import (
	"mmio"
	"unsafe"
)

type systick struct {
	csr   mmio.U32 // Control and Status Register
	rvr   mmio.U32 // Reload Value Register
	cvr   mmio.U32 // Current Value Register
	calib mmio.U32 // Calibration Value Register
}

var stk = (*systick)(unsafe.Pointer(uintptr(0xe000e010)))

// Implemented as set/clear flags instead of enable/disable functions because
// any read of ctrl register clears Count flag.

type Flags uint32

const (
	Enable  Flags = 1 << 0  // Counter ebabled.
	TickInt Flags = 1 << 1  // Generate exceptions.
	ClkCPU  Flags = 1 << 2  // Use CPU clock as clock source.
	Count   Flags = 1 << 16 // Timer counted to 0 since last flag read.
)

// Flags returns SysTick status and control flags.
func LoadFlags() Flags {
	return Flags(stk.csr.Load())
}

func StoreFlags(f Flags) {
	stk.csr.Store(uint32(f))
}

// SetFlags sets all flags specified by f.
func SetFlags(f Flags) {
	stk.csr.SetBits(uint32(f))
}

// ClearFlags resets flags specified by f.
func ClearFlags(f Flags) {
	stk.csr.ClearBits(uint32(f))
}

// SetReload sets SysTick reload value register. v can be in the range
// [0, 0x00ffffff].
func SetReload(v uint32) {
	stk.rvr.Store(v)
}

// Reload returns value of RVR register.
func Reload() uint32 {
	return stk.rvr.LoadMask(0x00ffffff)
}

// Val returns current value of SysTick counter
func Val() uint32 {
	return stk.cvr.LoadMask(0x00ffffff)
}

// Reset counter to 0.
func Reset() {
	stk.cvr.Store(0) // Anything stored clears cvr.
}

// Calib returns calibration properties.
func Calib() (skew, noref bool, tenms uint32) {
	c := stk.calib.Load()
	skew = c&(1<<31) != 0
	noref = c&(1<<30) != 0
	tenms = c & 0x00ffffff
	return
}
