package main

const CIA = 0 // Common IO Area (FN0)

// Vendor unique registers
const (
	SEP_INT_CTL = 0xF2
)

// SEP_INT_CTL bits
const (
	SEP_INTR_CTL_MASK = 1 << 0
	SEP_INTR_CTL_EN   = 1 << 1
	SEP_INTR_CTL_POL  = 1 << 2
)

const BP = 1 // Backplain function number (FN1)

// BCM43362 FN1 (Backplane) registers
const (
	GPIO_SELECT            = 0x10005
	GPIO_OUTPUT            = 0x10006
	GPIO_ENABLE            = 0x10007
	FUNCTION2_WATERMARK    = 0x10008
	DEVICE_CONTROL         = 0x10009
	BACKPLANE_ADDRESS_LOW  = 0x1000A
	BACKPLANE_ADDRESS_MID  = 0x1000B
	BACKPLANE_ADDRESS_HIGH = 0x1000C
	FRAME_CONTROL          = 0x1000D
	CHIP_CLOCK_CSR         = 0x1000E
	PULL_UP                = 0x1000F
	READ_FRAME_BC_LOW      = 0x1001B
	READ_FRAME_BC_HIGH     = 0x1001C
)

// CHIP_CLOCK_CSR bits
const (
	SBSDIO_FORCE_ALP           = 1 << 0
	SBSDIO_FORCE_HT            = 1 << 1
	SBSDIO_FORCE_ILP           = 1 << 2
	SBSDIO_ALP_AVAIL_REQ       = 1 << 3
	SBSDIO_HT_AVAIL_REQ        = 1 << 4
	SBSDIO_FORCE_HW_CLKREQ_OFF = 1 << 5
	SBSDIO_ALP_AVAIL           = 1 << 6
	SBSDIO_HT_AVAIL            = 1 << 7
)

// SDIO_FRAME_CONTROL bits
const (
	SFC_RF_TERM  = 1 << 0
	SFC_WF_TERM  = 1 << 1
	SFC_CRC4WOOS = 1 << 2
	SFC_ABORTALL = 1 << 3
)
