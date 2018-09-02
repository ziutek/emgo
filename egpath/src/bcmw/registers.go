package bcmw

// Functions
const (
	cia       = 0 // Function 0
	backplane = 1 // Function 1
	wlanData  = 2 // Function 2
)

// Backplane registers
const (
	sbsdioSPROMCS         = 0x10000
	sbsdioSPROMInfo       = 0x10001
	sbsdioSPROMDataLow    = 0x10002
	sbsdioSPROMDataHigh   = 0x10003
	sbsdioSPROMAddrLow    = 0x10004
	sbsdioSPROMAddrHigh   = 0x10005
	sbsdioChipCtrlData    = 0x10006
	sbsdioChipCtrlEn      = 0x10007
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
