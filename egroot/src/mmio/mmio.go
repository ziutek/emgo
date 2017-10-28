// Package mmio provides data types that can be used to access memory mapped
// registers of peripherals.
package mmio

import (
	"bits"
	"sync/atomic"
	"unsafe"
)

//c:volatile
type U8 struct {
	r uint8
}

func PtrU8(addr unsafe.Pointer) *U8 {
	return (*U8)(addr)
}

func AsU8(addr *uint8) *U8 {
	return (*U8)(unsafe.Pointer(addr))
}

func (r *U8) Addr() uintptr {
	return uintptr(unsafe.Pointer(r))
}

func (r *U8) SetBit(n int) {
	r.r |= uint8(1) << uint(n)
}

func (r *U8) ClearBit(n int) {
	r.r &^= uint8(1) << uint(n)
}

func (r *U8) Bit(n int) int {
	return int(r.r>>uint(n)) & 1
}

func (r *U8) StoreBit(n, v int) {
	mask := uint8(1) << uint(n)
	r.r = r.r&^mask | uint8(v<<uint(n))&mask
}

func (r *U8) Bits(mask uint8) uint8 {
	return r.r & mask
}

func (r *U8) StoreBits(mask, bits uint8) {
	r.r = r.r&^mask | bits&mask
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

func (r *U8) Field(mask uint8) int {
	return bits.Field32(uint32(r.r), uint32(mask))
}

func (r *U8) SetField(mask uint8, v int) {
	r.StoreBits(mask, uint8(bits.Make32(v, uint32(mask))))
}

type UM8 struct {
	U    *U8
	Mask uint8
}

func (b UM8) Set()             { b.U.SetBits(b.Mask) }
func (b UM8) Clear()           { b.U.ClearBits(b.Mask) }
func (b UM8) Load() uint8      { return b.U.Bits(b.Mask) }
func (b UM8) Store(bits uint8) { b.U.StoreBits(b.Mask, bits) }
func (b UM8) LoadVal() int     { return b.U.Field(uint8(b.Mask)) }
func (b UM8) StoreVal(v int)   { b.U.SetField(b.Mask, v) }

//c:volatile
type U16 struct {
	r uint16
}

func PtrU16(addr unsafe.Pointer) *U16 {
	return (*U16)(addr)
}

func AsU16(addr *uint16) *U16 {
	return (*U16)(unsafe.Pointer(addr))
}

func (r *U16) Addr() uintptr {
	return uintptr(unsafe.Pointer(r))
}

func (r *U16) SetBit(n int) {
	r.r |= uint16(1) << uint(n)
}

func (r *U16) ClearBit(n int) {
	r.r &^= uint16(1) << uint(n)
}

func (r *U16) Bit(n int) int {
	return int(r.r>>uint(n)) & 1
}

func (r *U16) StoreBit(n, v int) {
	mask := uint16(1) << uint(n)
	r.r = r.r&^mask | uint16(v<<uint(n))&mask
}

func (r *U16) Bits(mask uint16) uint16 {
	return r.r & mask
}

func (r *U16) StoreBits(mask, bits uint16) {
	r.r = r.r&^mask | bits&mask
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

func (r *U16) Field(mask uint16) int {
	return bits.Field32(uint32(r.r), uint32(mask))
}

func (r *U16) SetField(mask uint16, v int) {
	r.StoreBits(mask, uint16(bits.Make32(v, uint32(mask))))
}

type UM16 struct {
	U    *U16
	Mask uint16
}

func (b UM16) Set()              { b.U.SetBits(b.Mask) }
func (b UM16) Clear()            { b.U.ClearBits(b.Mask) }
func (b UM16) Load() uint16      { return b.U.Bits(b.Mask) }
func (b UM16) Store(bits uint16) { b.U.StoreBits(b.Mask, bits) }
func (b UM16) LoadVal() int      { return b.U.Field(uint16(b.Mask)) }
func (b UM16) StoreVal(v int)    { b.U.SetField(b.Mask, v) }

//c:volatile
type U32 struct {
	r uint32
}

func PtrU32(addr unsafe.Pointer) *U32 {
	return (*U32)(addr)
}

func AsU32(addr *uint32) *U32 {
	return (*U32)(unsafe.Pointer(addr))
}

func (r *U32) Addr() uintptr {
	return uintptr(unsafe.Pointer(r))
}

func (r *U32) SetBit(n int) {
	r.r |= uint32(1) << uint(n)
}

func (r *U32) ClearBit(n int) {
	r.r &^= uint32(1) << uint(n)
}

func (r *U32) Bit(n int) int {
	return int(r.r>>uint(n)) & 1
}

func (r *U32) StoreBit(n, v int) {
	mask := uint32(1) << uint(n)
	r.r = r.r&^mask | uint32(v<<uint(n))&mask
}
func (r *U32) Bits(mask uint32) uint32 {
	return r.r & mask
}

func (r *U32) StoreBits(mask, bits uint32) {
	r.r = r.r&^mask | bits&mask
}

func (r *U32) SetBits(mask uint32) {
	r.r |= mask
}

func (r *U32) AtomicSetBits(mask uint32) {
	atomic.OrUint32(&r.r, mask)
}

func (r *U32) ClearBits(mask uint32) {
	r.r &^= mask
}

func (r *U32) AtomicClearBits(mask uint32) {
	atomic.XorUint32(&r.r, mask)
}

func (r *U32) Load() uint32 {
	return r.r
}

func (r *U32) Store(v uint32) {
	r.r = v
}

func (r *U32) Field(mask uint32) int {
	return bits.Field32(r.r, mask)
}

func (r *U32) SetField(mask uint32, v int) {
	r.StoreBits(mask, bits.Make32(v, mask))
}

type UM32 struct {
	U    *U32
	Mask uint32
}

func (b UM32) Set()              { b.U.SetBits(b.Mask) }
func (b UM32) Clear()            { b.U.ClearBits(b.Mask) }
func (b UM32) Load() uint32      { return b.U.Bits(b.Mask) }
func (b UM32) Store(bits uint32) { b.U.StoreBits(b.Mask, bits) }
func (b UM32) LoadVal() int      { return b.U.Field(uint32(b.Mask)) }
func (b UM32) StoreVal(v int)    { b.U.SetField(b.Mask, v) }
