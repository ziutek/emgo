package mmio

import "unsafe"

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

func (r *U8) StoreBits(bits, mask uint8) {
	r.r = r.r&^mask | bits
}

func (r *U8) SetBits(mask uint8) {
	r.r |= mask
}

func (r *U8) ClearBits(mask uint8) {
	r.r &^= mask
}

func (r *U8) LoadBits(mask uint8) uint8 {
	return r.r & mask
}

func (r *U8) Load() uint8 {
	return r.r
}

func (r *U8) Store(v uint8) {
	r.r = v
}

func (r *U8) Ptr() uintptr {
	return uintptr(unsafe.Pointer(r))
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

func (r *U16) StoreBits(bits, mask uint16) {
	r.r = r.r&^mask | bits
}

func (r *U16) SetBits(mask uint16) {
	r.r |= mask
}

func (r *U16) ClearBits(mask uint16) {
	r.r &^= mask
}

func (r *U16) LoadBits(mask uint16) uint16 {
	return r.r & mask
}

func (r *U16) Load() uint16 {
	return r.r
}

func (r *U16) Store(v uint16) {
	r.r = v
}

func (r *U16) Ptr() uintptr {
	return uintptr(unsafe.Pointer(r))
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

func (r *U32) LoadBits(mask uint32) uint32 {
	return r.r & mask
}

func (r *U32) StoreBits(bits, mask uint32) {
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

func (r *U32) Ptr() uintptr {
	return uintptr(unsafe.Pointer(r))
}
