package scb

import (
	"mmio"
	"unsafe"
)

type Mask uint32
type Bit uint

type SCB struct {
	ACTLR mmio.U32
	_     [829]mmio.U32
	CPUID mmio.U32
	ICSR  mmio.U32
	VTOR  mmio.U32
	AIRCR mmio.U32
	SCR   mmio.U32
	CCR   mmio.U32
	SHPR1 mmio.U32
	SHPR2 mmio.U32
	SHPR3 mmio.U32
	SHCRS mmio.U32
	MMSR  mmio.U8
	BFSR  mmio.U8
	UFSR  mmio.U16
	HFSR  mmio.U32
	_     mmio.U32
	MMAR  mmio.U32
	BFAR  mmio.U32
	AFSR  mmio.U32
}

var R = (*SCB)(unsafe.Pointer(uintptr(0xe000e008)))

// ACTLR
const (
	DISMCYCINT Bit = 0
	DISDEFWBUF Bit = 1
	DISFOLD    Bit = 2
	DISFPCA    Bit = 8
	DISOOFP    Bit = 9
)

// CPUID
const (
	Revision    Mask = 0xf
	PartNo      Mask = 0xfff << 4
	Constant    Mask = 0xf << 16
	Variant     Mask = 0xf << 20
	Implementer Mask = 0xff << 24
)

// ICSR
const (
	VECTACTIVE  Mask = 0x1ff
	RETTOBASE   Bit  = 11
	VECTPENDING Mask = 0x3ff << 12
	ISRPENDING  Bit  = 22
	PENDSTCLR   Bit  = 25
	PENDSTSET   Bit  = 26
	PENDSVCLR   Bit  = 27
	PENDSVSET   Bit  = 28
	NMIPENDSET  Bit  = 31
)

// AIRCR
const (
	VECTRESET     Bit  = 0
	VECTCLRACTIVE Bit  = 1
	SYSRESETREQ   Bit  = 2
	PRIGROUP      Mask = 0x3 << 8
	ENDIANNESS    Bit  = 15
	VECTKEY       Mask = 0xffff << 16
)
