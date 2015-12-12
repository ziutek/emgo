package mmio

import "unsafe"

// Fsiz is used by Field, SetField methods. Field descriptor
// is uint16 value: f = FieldFsize<<Fsiz + FieldPos.
const Fsiz = 8

type U8 struct {
	r uint8
} //c:volatile

func PtrU8(addr unsafe.Pointer) *U8 {
	return (*U8)(addr)
}

func AsU8(addr *uint8) *U8 {
	return (*U8)(unsafe.Pointer(addr))
}

func (r *U8) SetBit(n int) {
	r.r |= uint8(1 << uint(n))
}

func (r *U8) ClearBit(n int) {
	r.r &^= uint8(1 << uint(n))
}

func (r *U8) Bit(n int) bool {
	return r.r&uint8(1<<uint(n)) != 0
}

func (r *U8) Bits(mask uint8) uint8 {
	return r.r & mask
}

func (r *U8) StoreBits(mask, bits uint8) {
	r.r = r.r&^mask | bits
}

func (r *U8) SetBits(mask uint8) {
	r.r |= mask
}

func (r *U8) ClearBits(mask uint8) {
	r.r &^= mask
}

func (r *U8) Load() uint8 {
	return r.r
}

func (r *U8) Store(v uint8) {
	r.r = v
}

func (r *U8) Addr() uintptr {
	return uintptr(unsafe.Pointer(r))
}

func (u *U8) Field(f uint16) int {
	o := byte(f)
	m := uint8(1)<<(f>>Fsiz) - 1
	return int(u.Bits(m<<o) >> o)
}

func (u *U8) SetField(f uint16, v int) {
	o := byte(f)
	m := uint8(1)<<(f>>Fsiz) - 1
	u.StoreBits(m<<o, uint8(v)<<o)
}

type U16 struct {
	r uint16
} //c:volatile

func PtrU16(addr unsafe.Pointer) *U16 {
	return (*U16)(addr)
}

func AsU16(addr *uint16) *U16 {
	return (*U16)(unsafe.Pointer(addr))
}

func (r *U16) SetBit(n int) {
	r.r |= uint16(1 << uint(n))
}

func (r *U16) ClearBit(n int) {
	r.r &^= uint16(1 << uint(n))
}

func (r *U16) Bit(n int) bool {
	return r.r&uint16(1<<uint(n)) != 0
}

func (r *U16) Bits(mask uint16) uint16 {
	return r.r & mask
}

func (r *U16) StoreBits(mask, bits uint16) {
	r.r = r.r&^mask | bits
}

func (r *U16) SetBits(mask uint16) {
	r.r |= mask
}

func (r *U16) ClearBits(mask uint16) {
	r.r &^= mask
}

func (r *U16) Load() uint16 {
	return r.r
}

func (r *U16) Store(v uint16) {
	r.r = v
}

func (r *U16) Addr() uintptr {
	return uintptr(unsafe.Pointer(r))
}

func (u *U16) Field(f uint16) int {
	o := byte(f)
	m := uint16(1)<<(f>>Fsiz) - 1
	return int(u.Bits(m<<o) >> o)
}

func (u *U16) SetField(f uint16, v int) {
	o := byte(f)
	m := uint16(1)<<(f>>Fsiz) - 1
	u.StoreBits(m<<o, uint16(v)<<o)
}

type U32 struct {
	r uint32
} //c:volatile

func PtrU32(addr unsafe.Pointer) *U32 {
	return (*U32)(addr)
}

func AsU32(addr *uint32) *U32 {
	return (*U32)(unsafe.Pointer(addr))
}

func (r *U32) SetBit(n int) {
	r.r |= uint32(1 << uint(n))
}

func (r *U32) ClearBit(n int) {
	r.r &^= uint32(1 << uint(n))
}

func (r *U32) Bit(n int) bool {
	return r.r&uint32(1<<uint(n)) != 0
}

func (r *U32) Bits(mask uint32) uint32 {
	return r.r & mask
}

func (r *U32) StoreBits(mask, bits uint32) {
	r.r = r.r&^mask | bits
}

func (r *U32) SetBits(mask uint32) {
	r.r |= mask
}

func (r *U32) ClearBits(mask uint32) {
	r.r &^= mask
}

func (r *U32) Load() uint32 {
	return r.r
}

func (r *U32) Store(v uint32) {
	r.r = v
}

func (r *U32) Addr() uintptr {
	return uintptr(unsafe.Pointer(r))
}

func (u *U32) Field(f uint16) int {
	o := byte(f)
	m := uint32(1)<<(f>>Fsiz) - 1
	return int(u.Bits(m<<o) >> o)
}

func (u *U32) SetField(f uint16, v int) {
	o := byte(f)
	m := uint32(1)<<(f>>Fsiz) - 1
	u.StoreBits(m<<o, uint32(v)<<o)
}

type Bits32 struct {
	Reg  *U32
	Mask uint32
}

func (b Bits32) Load() uint32      { return b.Reg.Bits(b.Mask) }
func (b Bits32) Store(bits uint32) { b.Reg.StoreBits(b.Mask, bits) }
func (b Bits32) Set()              { b.Reg.SetBits(b.Mask) }
func (b Bits32) Clear()            { b.Reg.ClearBits(b.Mask) }

type Field32 struct {
	Reg *U32
	Sel uint16
}

func (f Field32) Load() int   { return f.Reg.Field(f.Sel) }
func (f Field32) Store(v int) { f.Reg.SetField(f.Sel, v) }
