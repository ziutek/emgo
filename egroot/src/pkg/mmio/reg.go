package mmio

import "unsafe"

type Reg8 struct {
	r uint8 `C:"volatile"`
}

func NewReg8(addr uintptr) *Reg8 {
	return (*Reg8)(unsafe.Pointer(addr))
}

func (r *Reg8) SetBit(n int) {
	r.r |= uint8(1) << uint(n)
}

func (r *Reg8) ClearBit(n int) {
	r.r &^= uint8(1) << uint(n)
}

func (r *Reg8) Bit(n int) bool {
	return r.r&(uint8(1)<<uint(n)) != 0
}

func (r *Reg8) SetBits(bits, mask uint8) {
	r.r = r.r &^ mask | bits
}

func (r *Reg8) Bits(mask uint8) uint8 {
	return r.r & mask
}

func (r *Reg8) Read() uint8 {
	return r.r
}

func (r *Reg8) Write(v uint8) {
	r.r = v
}

type Reg16 struct {
	r uint16 `C:"volatile"`
}

func NewReg16(addr uintptr) *Reg16 {
	return (*Reg16)(unsafe.Pointer(addr))
}

func (r *Reg16) SetBit(n int) {
	r.r |= uint16(1) << uint(n)
}

func (r *Reg16) ClearBit(n int) {
	r.r &^= uint16(1) << uint(n)
}

func (r *Reg16) Bit(n int) bool {
	return r.r&(uint16(1)<<uint(n)) != 0
}

func (r *Reg16) SetBits(bits, mask uint16) {
	r.r = r.r &^ mask | bits
}

func (r *Reg16) Bits(mask uint16) uint16 {
	return r.r & mask
}

func (r *Reg16) Read() uint16 {
	return r.r
}

func (r *Reg16) Write(v uint16) {
	r.r = v
}

type Reg32 struct {
	r uint32 `C:"volatile"`
}

func NewReg32(addr uintptr) *Reg32 {
	return (*Reg32)(unsafe.Pointer(addr))
}

func (r *Reg32) SetBit(n int) {
	r.r |= uint32(1) << uint(n)
}

func (r *Reg32) ClearBit(n int) {
	r.r &^= uint32(1) << uint(n)
}

func (r *Reg32) Bit(n int) bool {
	return r.r&(uint32(1)<<uint(n)) != 0
}

func (r *Reg32) SetBits(bits, mask uint32) {
	r.r = r.r &^ mask | bits
}

func (r *Reg32) Bits(mask uint32) uint32 {
	return r.r & mask
}

func (r *Reg32) Read() uint32 {
	return r.r
}

func (r *Reg32) Write(v uint32) {
	r.r = v
}
