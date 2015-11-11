package clock

import "unsafe"

type regs struct {
	cr      uint32
	pllcfgr uint32
	cfgr    uint32
	cir     uint32
} //c:volatile

var c = (*regs)(unsafe.Pointer(uintptr(0x40023800)))

func ResetCR() {
	c.cr = 0x000083
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
	c.cr |= 1 << 19
}

func DisableSecurity() {
	c.cr &^= 1 << 19
}

func ResetPLLCFGR() {
	c.pllcfgr = 0x24003010
}

type PLLSrc byte

const (
	SrcHSI PLLSrc = 0
	SrcHSE PLLSrc = 1
)

func SetPLLClock(src PLLSrc) {
	c.pllcfgr = c.pllcfgr&^(1<<22) | uint32(src)<<22
}

// SetPLLInputDiv sets common input divisor for both PLLs.
// PLL input freq should be in range from 1 MHz to 2 MHz (recomended).
// Allowed values: 2 <= div <= 63
func SetPLLInputDiv(div int) {
	c.pllcfgr = c.pllcfgr&^0x3f | uint32(div&0x3f)
}

// SetMainPLLMul sets multipler for main PLL VCO.
// VCO output should be in range from 192 to 432 MHz.
// Allowed values: 2 <= mul <= 432,
func SetMainPLLMul(mul int) {
	c.pllcfgr =
		c.pllcfgr&^(0x1ff<<6) | uint32(mul&0x1ff)<<6
}

// SetMainPLLSysDiv sets divisor for main PLL SysClk output.
// SysClk should be <= 168 MHz.
// Allowed values: 2, 4, 6, 8.
func SetMainPLLSysDiv(div int) {
	div = (div >> 1) - 1
	c.pllcfgr = c.pllcfgr&^(3<<16) | uint32(div)<<16
}

// SetMainPLLPeriphDiv sets divisor for main PLL output used by USBFS,
// SDIO and RNG. USB OTG requires 48 MHz clock (SDIO and RNG <= 48 MHz).
// Allowed values: 2 <= div <= 15
func SetMainPLLPeriphDiv(div int) {
	c.pllcfgr = c.pllcfgr&^(0xf<<24) | uint32(div&0xf)<<24
}

func ResetCFGR() {
	c.cfgr = 0
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
	c.cfgr = c.cfgr&^(7<<10) | uint32(div)<<10
}

// SetPrescalerAPB2 sets prescaler for APB high-speed bus
func SetPrescalerAPB2(div APBDiv) {
	c.cfgr = c.cfgr&^(7<<13) | uint32(div)<<13
}

type Clock byte

const (
	HSI Clock = iota
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
