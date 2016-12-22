// Peripheral: SCB_Periph  System Control Block
// Instances:
//  SCB  0xE000ED00
// Registers:
//  0x00 32  CPUID CPUID Base Register
//  0x04 32  ICSR  Interrupt Control and State Register
//  0x08 32  VTOR  Vector Table Offset Register
//  0x0C 32  AIRCR Application Interrupt and Reset Control Register
//  0x10 32  SCR   System Control Register
//  0x14 32  CCR   Configuration and Control Register
//  0x18 32  SHPR1 System Handler Priority Register 1
//  0x1C 32  SHPR2 System Handler Priority Register 2
//  0x20 32  SHPR3 System Handler Priority Register 3
//  0x24 32  SHCSR System Handler Control and State Register
//  0x28 32  CFSR  Configurable Fault Status Register
//  0x2C 32  HFSR  HardFault Status Register
//  0x34 32  MMFR  MemManage Fault Address Register
//  0x38 32  BFAR  BusFault Address Register
//  0x3C 32  AFSR  Auxiliary Fault Status Register
package scb

const (
	Revision    CPUID_Bits = 0xf << 0   //+
	PartNo      CPUID_Bits = 0xfff << 4 //+
	Constant    CPUID_Bits = 0xf << 16  //+
	Variant     CPUID_Bits = 0xf << 20  //+
	Implementer CPUID_Bits = 0xff << 24 //+
)

const (
	Revision_n    = 0
	PartNo_n      = 4
	Constant_n    = 16
	Variant_n     = 20
	Implementer_n = 24
)

const (
	VECTACTIVE  ICSR_Bits = 0x1ff << 0  //+
	RETTOBASE   ICSR_Bits = 1 << 11     //+
	VECTPENDING ICSR_Bits = 0x3ff << 12 //+
	ISRPENDING  ICSR_Bits = 1 << 22     //+
	PENDSTCLR   ICSR_Bits = 1 << 25     //+
	PENDSTSET   ICSR_Bits = 1 << 26     //+
	PENDSVCLR   ICSR_Bits = 1 << 27     //+
	PENDSVSET   ICSR_Bits = 1 << 28     //+
	NMIPENDSET  ICSR_Bits = 1 << 31     //+
)

const (
	VECTACTIVEn  = 0
	VECTPENDINGn = 12
)

const (
	TBLOFF VTOR_Bits = 0x1ffffff << 7 //+
)

const (
	VECTRESET     AIRCR_Bits = 1 << 0       //+
	VECTCLRACTIVE AIRCR_Bits = 1 << 1       //+
	SYSRESETREQ   AIRCR_Bits = 1 << 2       //+
	PRIGROUP      AIRCR_Bits = 7 << 8       //+
	ENDIANNESS    AIRCR_Bits = 1 << 15      //+
	VECTKEY       AIRCR_Bits = 0xffff << 16 //+
)

const (
	VECTKEYn = 16
)

const (
	SLEEPONEXIT SCR_Bits = 1 << 1 //+
	SLEEPDEEP   SCR_Bits = 1 << 2 //+
	SEVONPEND   SCR_Bits = 1 << 4 //+
)

const (
	NONBASETHRDENA CCR_Bits = 1 << 0  //+
	USERSETMPEND   CCR_Bits = 1 << 1  //+
	UNALIGN_TRP    CCR_Bits = 1 << 3  //+
	DIV_0_TRP      CCR_Bits = 1 << 4  //+
	BFHFNMIGN      CCR_Bits = 1 << 8  //+
	STKALIGN       CCR_Bits = 1 << 9  //+
	DC             CCR_Bits = 1 << 16 //+
	IC             CCR_Bits = 1 << 17 //+
	BP             CCR_Bits = 1 << 18 //+
)

const (
	PRI_MemManage  SHPR1_Bits = 0xff << 0  //+
	PRI_BusFault   SHPR1_Bits = 0xff << 8  //+
	PRI_UsageFault SHPR1_Bits = 0xff << 16 //+
)

const (
	PRI_MemManage_n  = 0
	PRI_BusFault_n   = 8
	PRI_UsageFault_n = 16
)

const (
	PRI_SVCall SHPR2_Bits = 0xff << 24 //+
)

const (
	PRI_SVCall_n = 24
)

const (
	PRI_PendSV  SHPR3_Bits = 0xff << 16 //+
	PRI_SysTick SHPR3_Bits = 0xff << 24 //+
)

const (
	PRI_PendSV_n  = 16
	PRI_SysTick_n = 24
)

const (
	MEMFAULTACT    SHCSR_Bits = 1 << 0  //+
	BUSFAULTACT    SHCSR_Bits = 1 << 1  //+
	USGFAULTACT    SHCSR_Bits = 1 << 3  //+
	SVCALLACT      SHCSR_Bits = 1 << 7  //+
	MONITORACT     SHCSR_Bits = 1 << 8  //+
	PENDSVACT      SHCSR_Bits = 1 << 10 //+
	SYSTICKACT     SHCSR_Bits = 1 << 11 //+
	USGFAULTPENDED SHCSR_Bits = 1 << 12 //+
	MEMFAULTPENDED SHCSR_Bits = 1 << 13 //+
	BUSFAULTPENDED SHCSR_Bits = 1 << 14 //+
	SVCALLPENDED   SHCSR_Bits = 1 << 15 //+
	MEMFAULTENA    SHCSR_Bits = 1 << 16 //+
	BUSFAULTENA    SHCSR_Bits = 1 << 17 //+
	USGFAULTENA    SHCSR_Bits = 1 << 18 //+
)

const (
	// MFSR
	IACCVIOL  CFSR_Bits = 1 << 0 //+
	DACCVIOL  CFSR_Bits = 1 << 1 //+
	MUNSTKERR CFSR_Bits = 1 << 3 //+
	MSTKERR   CFSR_Bits = 1 << 4 //+
	MLSPERR   CFSR_Bits = 1 << 5 //+
	MMARVALID CFSR_Bits = 1 << 7 //+

	// BFSR
	IBUSERR     CFSR_Bits = 1 << 8  //+
	PRECISERR   CFSR_Bits = 1 << 9  //+
	IMPRECISERR CFSR_Bits = 1 << 10 //+
	UNSTKERR    CFSR_Bits = 1 << 11 //+
	STKERR      CFSR_Bits = 1 << 12 //+
	LSPERR      CFSR_Bits = 1 << 13 //+
	BFARVALID   CFSR_Bits = 1 << 15 //+

	// UFSR
	UNDEFINSTR CFSR_Bits = 1 << 16 //+
	INVSTATE   CFSR_Bits = 1 << 17 //+
	INVPC      CFSR_Bits = 1 << 18 //+
	NOCP       CFSR_Bits = 1 << 19 //+
	UNALIGNED  CFSR_Bits = 1 << 24 //+
	DIVBYZERO  CFSR_Bits = 1 << 25 //+
)

const (
	VECTTBL  HFSR_Bits = 1 << 1  //+
	FORCED   HFSR_Bits = 1 << 30 //+
	DEBUGEVT HFSR_Bits = 1 << 31 //+
)
