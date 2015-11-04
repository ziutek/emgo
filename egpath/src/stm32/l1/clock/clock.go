package clock

import "unsafe"

type regs struct {
	cr    uint32
	icscr uint32
	cfgr  uint32
	cir   uint32
} //C:volatile

var c = (*regs)(unsafe.Pointer(uintptr(0x40023800)))

func ResetCR() {
	c.cr = 0x00000300
}

func EnableHSI() {
	c.cr |= 1
}

func DisableHSI() {
	c.cr &^= 1
}

func HSIReady() bool {
	return c.cr&(1<<1) != 0
}

func EnableMSI() {
	c.cr |= 1 << 8
}

func DisableMSI() {
	c.cr &^= 1 << 8
}

func MSIReady() bool {
	return c.cr&(1<<9) != 0
}

func EnableHSE() {
	c.cr |= 1 << 16
}

func DisableHSE() {
	c.cr &^= 1 << 16
}

func HSEReady() bool {
	return c.cr&(1<<17) != 0
}

func SetHSEBypass() {
	c.cr |= 1 << 18
}

func ResetHSEBypass() {
	c.cr &^= 1 << 18
}

func EnableMainPLL() {
	c.cr |= 1 << 24
}

func DisableMainPLL() {
	c.cr &^= 1 << 24
}

func MainPLLReady() bool {
	return c.cr&(1<<25) != 0
}

func EnableSecurity() {
	c.cr |= 1 << 28
}

func DisableSecurity() {
	c.cr &^= 1 << 28
}

func ResetICSCR() {
	c.icscr = 0xB000
}

type MSIRange byte

const (
	Range65k5 MSIRange = iota
	Range131k
	Range262k
	Range524k
	Range1M05
	Range2M10 // default
	Range4M19
)

func SetMSIRange(rang MSIRange) {
	c.icscr = c.icscr&^(7<<13) | uint32(rang)<<13
}

func ResetCFGR() {
	c.cfgr = 0
}

type PLLSrc byte

const (
	SrcHSI PLLSrc = 0
	SrcHSE PLLSrc = 1
)

func PLLClock() PLLSrc {
	return PLLSrc((c.cfgr >> 16) & 1)
}

func SetPLLClock(src PLLSrc) {
	c.cfgr = c.cfgr&^(1<<16) | uint32(src)<<16
}

type PLLMul byte

const (
	PLLMul3 PLLMul = iota
	PLLMul4
	PLLMul6
	PLLMul8
	PLLMul12
	PLLMul16
	PLLMul24
	PLLMul32
	PLLMul48
)

// SetPLLMul sets multipler for PLL VCO.
// If USB or SDIO interface is used VCO must output 96 MHz.
// PLL VCO should avoid exceeding:
// 96 MHz for product voltage range 1
// 48 MHz for product voltage range 2
// 24 MHz for product voltage range 3.
func SetPLLMul(mul PLLMul) {
	c.cfgr = c.cfgr&^(0xf<<18) | uint32(mul)<<18
}

// SetPLLDiv sets divisor for PLL SysClk output.
// Allowed values: 2, 3, 4.
// SysClock frequency should avoid exceeding 32 MHz.
func SetPLLDiv(div int) {
	c.cfgr = c.cfgr&^(3<<22) | uint32(div-1)<<22
}

type AHBDiv byte

const (
	AHBDiv1   AHBDiv = 0
	AHBDiv2   AHBDiv = 8
	AHBDiv4   AHBDiv = 9
	AHBDiv8   AHBDiv = 10
	AHBDiv16  AHBDiv = 11
	AHBDiv64  AHBDiv = 12
	AHBDiv128 AHBDiv = 13
	AHBDiv256 AHBDiv = 14
	AHBDiv512 AHBDiv = 15
)

// SetPrescalerAHB sets prescaler for AHB bus
func SetPrescalerAHB(div AHBDiv) {
	c.cfgr = c.cfgr&^(0xf<<4) | uint32(div)<<4
}

type APBDiv byte

const (
	APBDiv1  APBDiv = 0
	APBDiv2  APBDiv = 4
	APBDiv4  APBDiv = 5
	APBDiv8  APBDiv = 6
	APBDiv16 APBDiv = 7
)

// SetPrescalerAPB1 sets prescaler for APB low-speed bus
func SetPrescalerAPB1(div APBDiv) {
	c.cfgr = c.cfgr&^(7<<8) | uint32(div)<<8
}

// SetPrescalerAPB2 sets prescaler for APB high-speed bus
func SetPrescalerAPB2(div APBDiv) {
	c.cfgr = c.cfgr&^(7<<11) | uint32(div)<<11
}

type Clock byte

const (
	MSI Clock = iota
	HSI
	HSE
	PLL
)

func SysClock() Clock {
	return Clock((c.cfgr >> 2) & 3)
}

func SetSysClock(src Clock) {
	c.cfgr = c.cfgr&^3 | uint32(src)
}

func ResetCIR() {
	c.cir = 0
}
