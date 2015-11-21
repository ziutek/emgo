package scb

import (
	"mmio"
	"unsafe"
)

type Mask uint32
type Field uint16

const x = 8

type Reg int

func reg(r Reg) *mmio.U32 {
	return &(*[num]mmio.U32)(unsafe.Pointer(uintptr(base)))[r]
}

func (r Reg) Bits(m Mask) uint32         { return reg(r).Bits(uint32(m)) }
func (r Reg) StoreBits(m Mask, b uint32) { reg(r).StoreBits(uint32(m), b) }
func (r Reg) SetBits(m Mask)             { reg(r).SetBits(uint32(m)) }
func (r Reg) ClearBits(m Mask)           { reg(r).ClearBits(uint32(m)) }
func (r Reg) Field(f Field) uint32       { return reg(r).Field(uint16(f)) }
func (r Reg) SetField(f Field, v uint32) { reg(r).SetField(uint16(f), v) }
func (r Reg) Load() uint32               { return reg(r).Load() }
func (r Reg) Store(v uint32)             { reg(r).Store(v) }
