// Package acc gives an access to the Access Control registers.
// Detailed description of all registers covered by this package can be found in
// "Cortex-M7 Devices Generic User Guide", chapter 4 "Cortex-M7
// Peripherals".
//
// Peripheral: ACC_Periph  Access Control
// Instances:
//  ACC  0xE000EF90
// Registers:
//  0x00 32  ITCMCR  Instruction Tightly-Coupled Memory Control Register
//  0x04 32  DTCMCR  Data Tightly-Coupled Memory Control Register
//	0x08 32  AHBPCR  AHBP Control Register
//  0x0C 32  CACR    L1 Cache Control Register
//  0x10 32  AHBSCR  AHB Slave Control Register
//  0x18 32  ABFSR   Auxiliary Bus Fault Status Register
package acc

const (
	ITCMEN    ITCMCR = 1 << 0   //+ TCM enable.
	ITCMRMW   ITCMCR = 1 << 1   //+ Read-Modify-Write (RMW) enable.
	ITCMRETEN ITCMCR = 1 << 2   //+ Retry phase enable.
	ITCMSZ    ITCMCR = 0xF << 3 //+ TCM size: 0:0K, 3:4K, 4:8K, ..., 16:16M.
)

const (
	ITCMENn    = 0
	ITCMRMWn   = 1
	ITCMRETENn = 2
	ITCMSZn    = 3
)

const (
	DTCMEN    DTCMCR = 1 << 0   //+ TCM enable.
	DTCMRMW   DTCMCR = 1 << 1   //+ Read-Modify-Write (RMW) enable.
	DTCMRETEN DTCMCR = 1 << 2   //+ Retry phase enable.
	DTCMSZ    DTCMCR = 0xF << 3 //+ TCM size. 0:0K, 3:4K, 4:8K, ..., 16:16M.
)

const (
	DTCMENn    = 0
	DTCMRMWn   = 1
	DTCMRETENn = 2
	DTCMSZn    = 3
)

const (
	AHBPEN AHBPCR = 1 << 0   //+ AHBP enable.
	AHBPSZ AHBPCR = 0x7 << 1 //+ AHBP size. 1:64M, 2:128M, 3:256M, 4:512M.
)

const (
	AHBPENn = 0
	AHBPSZn = 1
)

const (
	SIWT    CACR = 1 << 0 //+ Shared cacheable-is-WT for data cache.
	ECCDIS  CACR = 1 << 1 //+ ECC in the instruction and data cache.
	FORCEWT CACR = 1 << 2 //+ Force Write-Through in the data cache.
)

const (
	SIWTn    = 0
	ECCDISn  = 1
	FORCEWTn = 2
)

const (
	CTL       AHBSCR = 0x3 << 0   //+ AHBS prioritization control.
	TPRI      AHBSCR = 0x1FF << 2 //+ Thresh. exec. prio. for traffic demotion.
	INITCOUNT AHBSCR = 0x1F << 11 //+ Fairness counter initialization value.
)

const (
	CTLn       = 0
	TPRIn      = 2
	INITCOUNTn = 11
)

const (
	ITCM     ABFSR = 1 << 0   //+ Asynchronous fault on ITCM interface
	DTCM     ABFSR = 1 << 1   //+ Asynchronous fault on DTCM interface.
	AHBP     ABFSR = 1 << 2   //+ Asynchronous fault on AHBP interface.
	AXIM     ABFSR = 1 << 3   //+ Asynchronous fault on AXIM interface.
	EPPB     ABFSR = 1 << 4   //+ Asynchronous fault on EPPB interface.
	AXIMTYPE ABFSR = 0x3 << 8 //+ The type of fault on the AXIM interface.
)

const (
	ITCMn     = 0
	DTCMn     = 1
	AHBPn     = 2
	AXIMn     = 3
	EPPBn     = 4
	AXIMTYPEn = 8
)
