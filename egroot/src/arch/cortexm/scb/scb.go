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
	SHPR  [3 * 4]U8
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
	VECTPENDING Field = 10<<o + 12
	ISRPENDING  Bit   = 22
	PENDSTCLR   Bit   = 25
	PENDSTSET   Bit   = 26
	PENDSVCLR   Bit   = 27
	PENDSVSET   Bit   = 28
	NMIPENDSET  Bit   = 31
)

// VTOR
const (
	TBLOFF Field = 25<<o + 7
)

// AIRCR
const (
	VECTRESET     Bit   = 0
	VECTCLRACTIVE Bit   = 1
	SYSRESETREQ   Bit   = 2
	PRIGROUP      Field = 3<<o + 8
	ENDIANNESS    Bit   = 15
	VECTKEY       Field = 16<<o + 16
)

// SCR
const (
	SLEEPONEXIT Bit = 1
	SLEEPDEEP   Bit = 2
	SEVONPEND   Bit = 4
)

// CCR
const (
	NONBASETHRDENA Bit = 0
	USERSETMPEND   Bit = 1
	UNALIGN_TRP    Bit = 3
	DIV_0_TRP      Bit = 4
	BFHFNMIGN      Bit = 8
	STKALIGN       Bit = 9
)

// SHPR
const (
	MemManage  = 0
	BusFault   = 1
	UsageFault = 2
	SVCall     = 7
	PendSV     = 10
	SysTick    = 11
)

// SHCSR
const (
	MEMFAULTACT    Bit = 0
	BUSFAULTACT    Bit = 1
	USGFAULTACT    Bit = 3
	SVCALLACT      Bit = 7
	MONITORACT     Bit = 8
	PENDSVACT      Bit = 10
	SYSTICKACT     Bit = 11
	USGFAULTPENDED Bit = 12
	MEMFAULTPENDED Bit = 13
	BUSFAULTPENDED Bit = 14
	SVCALLPENDED   Bit = 15
	MEMFAULTENA    Bit = 16
	BUSFAULTENA    Bit = 17
	USGFAULTENA    Bit = 18
)

// MMSR
const (
	IACCVIOL  Bit = 0
	DACCVIOL  Bit = 1
	MUNSTKERR Bit = 3
	MSTKERR   Bit = 4
	MLSPERR   Bit = 5
	MMARVALID Bit = 7
)

// BFSR
const (
	IBUSERR     Bit = 0
	PRECISERR   Bit = 1
	IMPRECISERR Bit = 2
	UNSTKERR    Bit = 3
	STKERR      Bit = 4
	LSPERR      Bit = 5
	BFARVALID   Bit = 7
)

// UFSR
const (
	UNDEFINSTR Bit = 0
	INVSTATE   Bit = 1
	INVPC      Bit = 2
	NOCP       Bit = 3
	UNALIGNED  Bit = 8
	DIVBYZERO  Bit = 9
)

// HFSR
const (
	VECTTBL  Bit = 1
	FORCED   Bit = 30
	DEBUGEVT Bit = 31
)
