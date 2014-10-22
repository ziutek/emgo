// Package mpu allows configure Cortex-M memory protection unit.
package mpu

import "unsafe"

type registers struct {
	typ   uint32 `C:"volatile"`
	ctrl  uint32 `C:"volatile"`
	rn    uint32 `C:"volatile"`
	rba   uint32 `C:"volatile"`
	ras   uint32 `C:"volatile"`
	rbaa1 uint32 `C:"volatile"`
	rasa1 uint32 `C:"volatile"`
	rbaa2 uint32 `C:"volatile"`
	rasa2 uint32 `C:"volatile"`
	rbaa3 uint32 `C:"volatile"`
	rasa3 uint32 `C:"volatile"`
}

var r = (*registers)(unsafe.Pointer(uintptr(0xE000ED90)))

// Type returns information about MPU unit:
// i - number of supported instruction regions,
// d - number of supported data regions.
// s - true if separate instruction and data regions are supported.
func Type() (i, d int, s bool) {
	i = int(r.typ>>16) & 0xf
	d = int(r.typ>>8) & 0xf
	s = (r.typ&1 != 0)
	return
}

type Flags uint32

const (
	// If HFNM is not set the MPU will be disabled during HardFault, NMI and
	// FAULTMASK handlers.
	HFNM Flags = 1 << (iota + 1)
	// If PRIVDEF is set the default memory map is used as background region for
	// privileged software access.
	PRIVDEF
)

// SetMode sets flags that globally determine the behavior of the MPU.
func SetMode(fl Flags) {
	r.ctrl = r.ctrl&1 | uint32(fl)
}

// Mode returns current flags.
func Mode() Flags {
	return Flags(r.ctrl &^ 1)
}

// Enable enables MPU.
func Enable() {
	r.ctrl |= 1
}

// Disable disables MPU.
func Disable() {
	r.ctrl &^= 1
}

type Attr uint32

const (
	B Attr = 1 << 16 // Bufferable.
	C Attr = 1 << 17 // Cacheable.
	S Attr = 1 << 18 // Shareable.

	// Access permissons.
	P____ Attr = 0 << 24 // No access.
	Pr___ Attr = 5 << 24 // Priv-RO.
	Prw__ Attr = 1 << 24 // Priv-RW.
	Pr_r_ Attr = 6 << 24 // Priv-RO, Unpriv-RO.
	Prwr_ Attr = 2 << 24 // Priv-RW, Unpriv-RO.
	Prwrw Attr = 3 << 24 // Priv-RW, Unpriv-RW.

	Xn Attr = 1 << 28 // Instruction Access Disable.
)

// SetTex sets type extension in a.
func (a *Attr) SetTex(tex byte) {
	*a |= Attr(tex&7) << 19
}

// Tex extracts type extension from a.
func (a Attr) Tex() byte {
	return byte(a>>19) & 7
}

// SetRegion setups region rn at address addr of size 1<<sizeExp.
// Any bit set in disable excludes 1/8 of memory (subregion) from region rn.
// Only regions of size >= 256B can be divided to subregions. The least
// significant bit of disable controls the first subregion. attr specifies
// attributes for region rn.
func SetRegion(rn int, addr uintptr, sizeExp int, disable byte, attr Attr) {

}
