package mmio

import (
	"bits"
	"unsafe"
)

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

func (r *U8) Field(mask uint8) int {
	return bits.Field32(uint32(r.r), uint32(mask))
}

func (r *U8) SetField(mask uint8, v int) {
	r.StoreBits(mask, uint8(bits.Make32(v, uint32(mask))))
}

type Bits8 struct {
	Reg  *U8
	Mask uint8
}

func (b Bits8) Set()             { b.Reg.SetBits(b.Mask) }
func (b Bits8) Clear()           { b.Reg.ClearBits(b.Mask) }
func (b Bits8) Load() uint8      { return b.Reg.Bits(b.Mask) }
func (b Bits8) Store(bits uint8) { b.Reg.StoreBits(b.Mask, bits) }
func (b Bits8) LoadVal() int     { return b.Reg.Field(uint8(b.Mask)) }
func (b Bits8) StoreVal(v int)   { b.Reg.SetField(b.Mask, v) }

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

func (r *U16) Field(mask uint16) int {
	return bits.Field32(uint32(r.r), uint32(mask))
}

func (r *U16) SetField(mask uint16, v int) {
	r.StoreBits(mask, uint16(bits.Make32(v, uint32(mask))))
}

type Bits16 struct {
	Reg  *U16
	Mask uint16
}

func (b Bits16) Set()              { b.Reg.SetBits(b.Mask) }
func (b Bits16) Clear()            { b.Reg.ClearBits(b.Mask) }
func (b Bits16) Load() uint16      { return b.Reg.Bits(b.Mask) }
func (b Bits16) Store(bits uint16) { b.Reg.StoreBits(b.Mask, bits) }
func (b Bits16) LoadVal() int      { return b.Reg.Field(uint16(b.Mask)) }
func (b Bits16) StoreVal(v int)    { b.Reg.SetField(b.Mask, v) }

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

func (r *U32) Field(mask uint32) int {
	return bits.Field32(r.r, mask)
}

func (r *U32) SetField(mask uint32, v int) {
	r.StoreBits(mask, bits.Make32(v, mask))
}

type Bits32 struct {
	Reg  *U32
	Mask uint32
}

func (b Bits32) Set()              { b.Reg.SetBits(b.Mask) }
func (b Bits32) Clear()            { b.Reg.ClearBits(b.Mask) }
func (b Bits32) Load() uint32      { return b.Reg.Bits(b.Mask) }
func (b Bits32) Store(bits uint32) { b.Reg.StoreBits(b.Mask, bits) }
func (b Bits32) LoadVal() int      { return b.Reg.Field(uint32(b.Mask)) }
func (b Bits32) StoreVal(v int)    { b.Reg.SetField(b.Mask, v) }
