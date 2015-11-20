package scb

import "unsafe"

type SCB struct {
	ACTLR U32
	_     [829]U32
	CPUID U32
	ICSR  U32
	VTOR  U32
	AIRCR U32
	SCR   U32
	CCR   U32
	SHPR1 U32
	SHPR2 U32
	SHPR3 U32
	SHCRS U32
	MMSR  U8
	BFSR  U8
	UFSR  U16
	HFSR  U32
	_     U32
	MMAR  U32
	BFAR  U32
	AFSR  U32
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
	Revision    Field = 4<<o + 0
	PartNo      Field = 12<<o + 4
	Constant    Field = 4<<o + 16
	Variant     Field = 4<<o + 20
	Implementer Field = 8<<o + 24
)

// ICSR
const (
	VECTACTIVE  Field = 9<<o + 0
	RETTOBASE   Bit   = 11
	VECTPENDING Field = 11<<o + 12
	ISRPENDING  Bit   = 22
	PENDSTCLR   Bit   = 25
	PENDSTSET   Bit   = 26
	PENDSVCLR   Bit   = 27
	PENDSVSET   Bit   = 28
	NMIPENDSET  Bit   = 31
)

// AIRCR
const (
	VECTRESET     Bit   = 0
	VECTCLRACTIVE Bit   = 1
	SYSRESETREQ   Bit   = 2
	PRIGROUP      Field = 2<<o + 8
	ENDIANNESS    Bit   = 15
	VECTKEY       Field = 16<<o + 16
)
