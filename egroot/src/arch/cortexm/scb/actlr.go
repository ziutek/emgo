package scb

import (
	"mmio"
	"unsafe"
)

type Reg1 int

func reg1(r Reg1) *mmio.U32 {
	return &(*[num]mmio.U32)(unsafe.Pointer(uintptr(0xe000e008)))[r]
}

func (r Reg1) Bits(m Mask) uint32         { return reg1(r).Bits(uint32(m)) }
func (r Reg1) StoreBits(m Mask, b uint32) { reg1(r).StoreBits(uint32(m), b) }
func (r Reg1) SetBits(m Mask)             { reg1(r).SetBits(uint32(m)) }
func (r Reg1) ClearBits(m Mask)           { reg1(r).ClearBits(uint32(m)) }
func (r Reg1) Field(f Field) uint32       { return reg1(r).Field(uint16(f)) }
func (r Reg1) SetField(f Field, v uint32) { reg1(r).SetField(uint16(f), v) }
func (r Reg1) Load() uint32               { return reg1(r).Load() }
func (r Reg1) Store(v uint32)             { reg1(r).Store(v) }

const (
	ACTLR Reg1 = 0

	DISMCYCINT Mask = 1 << 0
	DISDEFWBUF Mask = 1 << 1
	DISFOLD    Mask = 1 << 2
	DISFPCA    Mask = 1 << 8
	DISOOFP    Mask = 1 << 9
)
