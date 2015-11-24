// Package scb gives an access to registers of System Control Block.
//
// Detailed description of all registers covered by this package can be found in
// "Cortex-M[0-4] Devices Generic User Guide", chapter 4 "Cortex-M[0-4]
// Peripherals".
//
// Notes
//
// 1. Cortex-M0 doesn't implement ACTLR, SHPR1, CFSR, HFSR, MMFR, BFAR, AFSR
// registers.
//
// 2. SHPR2, SHPR3 registers are only word accessible in Cortex-M0 so this
// package doesn't provide byte access to individual fields.
package scb

const (
	base   = 0xE000ED00
	length = 16

	base1   = 0xe000e008
	length1 = 1
)

const (
	ACTLR Reg1 = 0

	DISMCYCINT Mask = 1 << 0
	DISDEFWBUF Mask = 1 << 1
	DISFOLD    Mask = 1 << 2
	DISFPCA    Mask = 1 << 8
	DISOOFP    Mask = 1 << 9
)

const (
	CPUID Reg = 0

	Revision    Mask  = 1<<4 - 1
	PartNo      Field = 12<<x + 4
	Constant    Field = 4<<x + 16
	Variant     Field = 4<<x + 20
	Implementer Field = 8<<x + 24
)

const (
	ICSR Reg = 1

	VECTACTIVE  Mask  = 1<<9 - 1
	RETTOBASE   Mask  = 1 << 11
	VECTPENDING Field = 10<<x + 12
	ISRPENDING  Mask  = 1 << 22
	PENDSTCLR   Mask  = 1 << 25
	PENDSTSET   Mask  = 1 << 26
	PENDSVCLR   Mask  = 1 << 27
	PENDSVSET   Mask  = 1 << 28
	NMIPENDSET  Mask  = 1 << 31
)

const (
	VTOR Reg = 2

	TBLOFF Mask = (1<<25 - 1) << 7
)

const (
	AIRCR Reg = 3

	VECTRESET     Mask  = 1 << 0
	VECTCLRACTIVE Mask  = 1 << 1
	SYSRESETREQ   Mask  = 1 << 2
	PRIGROUP      Field = 3<<x + 8
	ENDIANNESS    Mask  = 1 << 15
	VECTKEY       Field = 16<<x + 16
)

const (
	SCR Reg = 4

	SLEEPONEXIT Mask = 1 << 1
	SLEEPDEEP   Mask = 1 << 2
	SEVONPEND   Mask = 1 << 4
)

const (
	CCR Reg = 5

	NONBASETHRDENA Mask = 1 << 0
	USERSETMPEND   Mask = 1 << 1
	UNALIGN_TRP    Mask = 1 << 3
	DIV_0_TRP      Mask = 1 << 4
	BFHFNMIGN      Mask = 1 << 8
	STKALIGN       Mask = 1 << 9
)

const (
	SHPR1 Reg = 6

	MEMMANAGE  Field = 8<<x + 0
	BUSFAULT   Field = 8<<x + 8
	USAGEFAULT Field = 8<<x + 16
)

const (
	SHPR2 Reg = 7

	SVCALL Field = 8<<x + 24
)

const (
	SHPR3 Reg = 8

	PENDSV  Field = 8<<x + 16
	SYSTICK Field = 8<<x + 24
)

const (
	SHCSR Reg = 9

	MEMFAULTACT    Mask = 1 << 0
	BUSFAULTACT    Mask = 1 << 1
	USGFAULTACT    Mask = 1 << 3
	SVCALLACT      Mask = 1 << 7
	MONITORACT     Mask = 1 << 8
	PENDSVACT      Mask = 1 << 10
	SYSTICKACT     Mask = 1 << 11
	USGFAULTPENDED Mask = 1 << 12
	MEMFAULTPENDED Mask = 1 << 13
	BUSFAULTPENDED Mask = 1 << 14
	SVCALLPENDED   Mask = 1 << 15
	MEMFAULTENA    Mask = 1 << 16
	BUSFAULTENA    Mask = 1 << 17
	USGFAULTENA    Mask = 1 << 18
)

const (
	CFSR Reg = 10

	// MFSR
	IACCVIOL  Mask = 1 << 0
	DACCVIOL  Mask = 1 << 1
	MUNSTKERR Mask = 1 << 3
	MSTKERR   Mask = 1 << 4
	MLSPERR   Mask = 1 << 5
	MMARVALID Mask = 1 << 7

	// BFSR
	IBUSERR     Mask = 1 << 8
	PRECISERR   Mask = 1 << 9
	IMPRECISERR Mask = 1 << 10
	UNSTKERR    Mask = 1 << 11
	STKERR      Mask = 1 << 12
	LSPERR      Mask = 1 << 13
	BFARVALID   Mask = 1 << 15

	// UFSR
	UNDEFINSTR Mask = 1 << 16
	INVSTATE   Mask = 1 << 17
	INVPC      Mask = 1 << 18
	NOCP       Mask = 1 << 19
	UNALIGNED  Mask = 1 << 24
	DIVBYZERO  Mask = 1 << 25
)

const (
	HFSR Reg = 11

	VECTTBL  Mask = 1 << 1
	FORCED   Mask = 1 << 30
	DEBUGEVT Mask = 1 << 31
)

const (
	MMFR Reg = 13
	BFAR Reg = 14
	AFSR Reg = 15
)
