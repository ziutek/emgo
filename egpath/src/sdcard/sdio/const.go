package sdio

// CIA registers

const (
	CCCR_REV       = 0x00 // CCCR/SDIO Revision
	CCCR_SDREV     = 0x01 // SD Revision
	CCCR_IOEN      = 0x02 // I/O Enable
	CCCR_IORDY     = 0x03 // I/O Ready
	CCCR_INTEN     = 0x04 // Int Enable
	CCCR_INTPEND   = 0x05 // Int Pending
	CCCR_IOABORT   = 0x06 // I/O Abort
	CCCR_BUSICTRL  = 0x07 // Bus Interface Control
	CCCR_CARDCAP   = 0x08 // Card Capabilities
	CCCR_CCISPTR0  = 0x09 // Common CIS Pointer 0 LSB
	CCCR_CCISPTR1  = 0x0A // Common CIS Pointer 1
	CCCR_CCISPTR2  = 0x0B // Common CIS Pointer 2 MSB
	CCCR_BUSSUSP   = 0x0C // Bus Suspend
	CCCR_FUNCSEL   = 0x0D // Function select
	CCCR_EXECFLAGS = 0x0E // Exec Flags
	CCCR_RDYFLAGS  = 0x0F // Ready Flags
	CCCR_BLKSIZE0  = 0x10 // FN0 Block Size 0 (LSB)
	CCCR_BLKSIZE1  = 0x11 // FN0 Block Size 1 (MSB)
	CCCR_POWERCTRL = 0x12 // Power Control
	CCCR_SPEEDSEL  = 0x13 // Bus Speed Select
	CCCR_UHSI      = 0x14 // UHS-I Support
	CCCR_DRIVE     = 0x15 // Drive Strength
	CCCR_INTEXT    = 0x16 // Interrupt Extension
)

const (
	FBR1 = 0x100 // Function Basic Registers for FN1
	FBR2 = 0x200 // Function Basic Registers for FN2
	FBR3 = 0x300 // Function Basic Registers for FN3
	FBR4 = 0x400 // Function Basic Registers for FN4
	FBR5 = 0x500 // Function Basic Registers for FN5
	FBR6 = 0x600 // Function Basic Registers for FN6
	FBR7 = 0x700 // Function Basic Registers for FN7
)

const (
	FBR_CSASFIC  = 0x00 // CSA, Standard SDIO Function Interface Code
	FBR_ESFIC    = 0x01 // Extended Standard SDIO Function Interface Code
	FBR_PS       = 0x02 // Power Selection/state
	FBR_CISPTR0  = 0x09 // Address pointer to Function CIS 0 (LSB)
	FBR_CISPTR1  = 0x0A // Address pointer to Function CIS 1
	FBR_CISPTR2  = 0x0B // Address pointer to Function CIS 2 (MSB)
	FBR_CSAPTR0  = 0x0C // Address pointer to Function CSA 0 (LSB)
	FBR_CSAPTR1  = 0x0D // Address pointer to Function CSA 1
	FBR_CSAPTR2  = 0x0E // Address pointer to Function CSA 2 (MSB)
	FBR_CSADAW   = 0x0F // Data access window to CSA
	FBR_BLKSIZE0 = 0x10 // Block Size 0 (LSB)
	FBR_BLKSIZE1 = 0x11 // Block Size 0 (MSB)
)
