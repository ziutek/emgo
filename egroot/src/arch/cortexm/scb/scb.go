// BaseAddr: 0xe000ed00
//   0: CPUID CPUID Base Register
//   1: ICSR  Interrupt Control and State Register
//   2: VTOR  Vector Table Offset Register
//   3: AIRCR Application Interrupt and Reset Control Register
//   4: SCR   System Control Register
//   5: CCR   Configuration and Control Register
//   6: SHPR1 System Handler Priority Register 1
//   7: SHPR2 System Handler Priority Register 2
//   8: SHPR3 System Handler Priority Register 3
//   9: SHCSR System Handler Control and State Register
//  10: CFSR  Configurable Fault Status Register
//  11: HFSR  HardFault Status Register
//  13: MMFR  MemManage Fault Address Register
//  14: BFAR  BusFault Address Register
//  15: AFSR  Auxiliary Fault Status Register
package scb

const (
	Revision    CPUID_Field = 4<<siz + 0
	PartNo      CPUID_Field = 12<<siz + 4
	Constant    CPUID_Field = 4<<siz + 16
	Variant     CPUID_Field = 4<<siz + 20
	Implementer CPUID_Field = 8<<siz + 24
)

const (
	VECTACTIVE  ICSR_Bits  = 1<<9 - 1
	RETTOBASE   ICSR_Bits  = 1 << 11
	VECTPENDING ICSR_Field = 10<<siz + 12
	ISRPENDING  ICSR_Bits  = 1 << 22
	PENDSTCLR   ICSR_Bits  = 1 << 25
	PENDSTSET   ICSR_Bits  = 1 << 26
	PENDSVCLR   ICSR_Bits  = 1 << 27
	PENDSVSET   ICSR_Bits  = 1 << 28
	NMIPENDSET  ICSR_Bits  = 1 << 31
)

const (
	TBLOFF VTOR_Bits = (1<<25 - 1) << 7
)

const (
	VECTRESET     AIRCR_Bits  = 1 << 0
	VECTCLRACTIVE AIRCR_Bits  = 1 << 1
	SYSRESETREQ   AIRCR_Bits  = 1 << 2
	PRIGROUP      AIRCR_Field = 3<<siz + 8
	ENDIANNESS    AIRCR_Bits  = 1 << 15
	VECTKEY       AIRCR_Field = 16<<siz + 16
)

const (
	SLEEPONEXIT SCR_Bits = 1 << 1
	SLEEPDEEP   SCR_Bits = 1 << 2
	SEVONPEND   SCR_Bits = 1 << 4
)

const (
	NONBASETHRDENA CCR_Bits = 1 << 0
	USERSETMPEND   CCR_Bits = 1 << 1
	UNALIGN_TRP    CCR_Bits = 1 << 3
	DIV_0_TRP      CCR_Bits = 1 << 4
	BFHFNMIGN      CCR_Bits = 1 << 8
	STKALIGN       CCR_Bits = 1 << 9
)

const (
	PRI_MemManage  SHPR1_Field = 8<<siz + 0
	PRI_BusFault   SHPR1_Field = 8<<siz + 8
	PRI_UsageFault SHPR1_Field = 8<<siz + 16
)

const (
	PRI_SVCall SHPR2_Field = 8<<siz + 24
)

const (
	PRI_PendSV  SHPR3_Field = 8<<siz + 16
	PRI_SysTick SHPR3_Field = 8<<siz + 24
)

const (
	MEMFAULTACT    SHCSR_Bits = 1 << 0
	BUSFAULTACT    SHCSR_Bits = 1 << 1
	USGFAULTACT    SHCSR_Bits = 1 << 3
	SVCALLACT      SHCSR_Bits = 1 << 7
	MONITORACT     SHCSR_Bits = 1 << 8
	PENDSVACT      SHCSR_Bits = 1 << 10
	SYSTICKACT     SHCSR_Bits = 1 << 11
	USGFAULTPENDED SHCSR_Bits = 1 << 12
	MEMFAULTPENDED SHCSR_Bits = 1 << 13
	BUSFAULTPENDED SHCSR_Bits = 1 << 14
	SVCALLPENDED   SHCSR_Bits = 1 << 15
	MEMFAULTENA    SHCSR_Bits = 1 << 16
	BUSFAULTENA    SHCSR_Bits = 1 << 17
	USGFAULTENA    SHCSR_Bits = 1 << 18
)

const (
	// MFSR
	IACCVIOL  CFSR_Bits = 1 << 0
	DACCVIOL  CFSR_Bits = 1 << 1
	MUNSTKERR CFSR_Bits = 1 << 3
	MSTKERR   CFSR_Bits = 1 << 4
	MLSPERR   CFSR_Bits = 1 << 5
	MMARVALID CFSR_Bits = 1 << 7

	// BFSR
	IBUSERR     CFSR_Bits = 1 << 8
	PRECISERR   CFSR_Bits = 1 << 9
	IMPRECISERR CFSR_Bits = 1 << 10
	UNSTKERR    CFSR_Bits = 1 << 11
	STKERR      CFSR_Bits = 1 << 12
	LSPERR      CFSR_Bits = 1 << 13
	BFARVALID   CFSR_Bits = 1 << 15

	// UFSR
	UNDEFINSTR CFSR_Bits = 1 << 16
	INVSTATE   CFSR_Bits = 1 << 17
	INVPC      CFSR_Bits = 1 << 18
	NOCP       CFSR_Bits = 1 << 19
	UNALIGNED  CFSR_Bits = 1 << 24
	DIVBYZERO  CFSR_Bits = 1 << 25
)

const (
	VECTTBL  HFSR_Bits = 1 << 1
	FORCED   HFSR_Bits = 1 << 30
	DEBUGEVT HFSR_Bits = 1 << 31
)
