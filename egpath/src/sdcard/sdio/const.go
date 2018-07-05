package sdio

const CIA = 0 // Function 0

// CIA registers

const (
	CCCR_REV        = 0x00 // CCCR/SDIO Revision
	CCCR_SDREV      = 0x01 // SD Revision
	CCCR_IOEN       = 0x02 // I/O Enable
	CCCR_IORDY      = 0x03 // I/O Ready
	CCCR_INTEN      = 0x04 // Int Enable1
	CCCR_INTPEND    = 0x05 // Int Pending
	CCCR_IOABORT    = 0x06 // I/O Abort
	CCCR_BUSICTRL   = 0x07 // Bus Interface Control
	CCCR_CARDCAP    = 0x08 // Card Capabilities
	CCCR_CCISPTR0   = 0x09 // Common CIS Pointer 0 LSB
	CCCR_CCISPTR1   = 0x0A // Common CIS Pointer 1
	CCCR_CCISPTR2   = 0x0B // Common CIS Pointer 2 MSB
	CCCR_BUSSUSP    = 0x0C // Bus Suspend
	CCCR_FUNCSEL    = 0x0D // Function select
	CCCR_EXECFLAGS  = 0x0E // Exec Flags
	CCCR_RDYFLAGS   = 0x0F // Ready Flags
	CCCR_BLOCKSIZE0 = 0x10 // FN0 Block Size LSB
	CCCR_BLOCKSIZE1 = 0x11 // FN0 Block Size MSB
	CCCR_POWERCTRL  = 0x12 // Power Control
	CCCR_SPEEDSEL   = 0x13 // Bus Speed Select
	CCCR_UHSI       = 0x14 // UHS-I Support
	CCCR_DRIVE      = 0x15 // Drive Strength
	CCCR_INTEXT     = 0x16 // Interrupt Extension
)

const (
	FN1 = 1 << 1
	FN2 = 1 << 2
	FN3 = 1 << 3
	FN4 = 1 << 4
	FN5 = 1 << 5
	FN6 = 1 << 6
	FN7 = 1 << 7
)
