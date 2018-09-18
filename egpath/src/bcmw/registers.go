package bcmw

// Functions
const (
	cia       = 0 // Function 0
	backplane = 1 // Function 1
	wlanData  = 2 // Function 2
)

// CCCR vendor specific registers
const (
	cccrCardCap   = 0xF0
	cccrCardCtl   = 0xF1
	cccrSepIntCtl = 0xF2
)

// cccrSepIntCtl bits
const (
	sepIntCtlMask = 1 << 0
	sepIntCtlEn   = 1 << 1
	sepIntCtlPol  = 1 << 2
)

// Backplane registers
const (
	sbsdioGPIOSel         = 0x10005
	sbsdioGPIOOut         = 0x10006
	sbsdioGPIOEn          = 0x10007
	sbsdioWatermark       = 0x10008
	sbsdioDeviceCtl       = 0x10009
	sbsdioFunc1SBAddrLow  = 0x1000A
	sbsdioFunc1SBAddrMid  = 0x1000B
	sbsdioFunc1SBAddrHigh = 0x1000C
	sbsdioFunc1FrameCtrl  = 0x1000D
	sbsdioFunc1ChipClkCSR = 0x1000E
	sbsdioFunc1SDIOPullUp = 0x1000F
)

// sbsdioFunc1ChipClkCSR bits
const (
	sbsdioForceALP         = 1 << 0
	sbsdioForceHT          = 1 << 1
	sbsdioForceILP         = 1 << 2
	sbsdioALPAvailReq      = 1 << 3
	sbsdioHTAvailReq       = 1 << 4
	sbsdioForceHwClkReqOff = 1 << 5
	sbsdioALPAvail         = 1 << 6
	sbsdioHTAvail          = 1 << 7
)

// Sonics Silicon Backplane (SSB) Core Registers
//
// Windowed access: base address (bits 15 to 31) is set in sbsdioFunc1SBAddr*.
// Less significant bits are specified in CMD52/CMD53. More info:
// http://www.gc-linux.org/wiki/Wii:WLAN

const sbsdioAccess32bit = 1 << 15

// Agent registers (common for every core).
// Source: linux/include/linux/bcma/bcma_regs.h
const (
	agentIOCtl    = 0x0408
	agentIOSt     = 0x0500
	agentResetCtl = 0x0800
	agentResetSt  = 0x0804
)

// agentIOCtl bits
const (
	ioCtlClk = 1 << 0
	ioCtlFGC = 1 << 1

	// Dot11MAC core specific control flag bits
	ioCtlDot11PhyClockEn = 1 << 2
	ioCtlDot11PhyReset   = 1 << 3
)

// Chip common registers
const (
	commonEnumBase = 0x18000000 + 0x00 // Chip ID
	commonGPIOCtl  = 0x18000000 + 0x6C
)

// SOCSRAM registers
const (
	socsramBankxIndex = 0x18004000 + 0x10
	socsramBankxPDA   = 0x18004000 + 0x44
)
