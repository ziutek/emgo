// Package rcc gives an access to STM32F405xx/07xx reset and clock control
// registers.
//
// BaseAddr: 0x40023800  AHB1
//  0x00: CR         Clock control register.
//  0x04: PLLCFGR    PLL configuration register.
//  0x08: CFGR       Clock configuration register.
//  0x0C: CIR        Clock interrupt register.
//  0x10: AHB1RSTR   AHB1 peripheral reset register.
//  0x14: AHB2RSTR   AHB2 peripheral reset register.
//  0x18: AHB3RSTR   AHB3 peripheral reset register.
//  0x20: APB1RSTR   APB1 peripheral reset register.
//  0x24: APB2RSTR   APB2 peripheral reset register.
//  0x30: AHB1ENR    AHB1 peripheral clock enable register.
//  0x34: AHB2ENR    AHB2 peripheral clock enable register.
//  0x38: AHB3ENR    AHB3 peripheral clock enable register.
//  0x40: APB1ENR    APB1 peripheral clock enable register.
//  0x44: APB2ENR    APB2 peripheral clock enable register.
//  0x50: AHB1LPENR  AHB1 peripheral clock enable in low power mode register.
//  0x54: AHB2LPENR  AHB2 peripheral clock enable in low power mode register.
//  0x58: AHB3LPENR  AHB3 peripheral clock enable in low power mode register.
//  0x60: APB1LPENR  APB1 peripheral clock enabled in low power mode register.
//  0x64: APB2LPENR  APB2 peripheral clock enabled in low power mode register.
//  0x70: BDCR       Backup domain control register.
//  0x74: CSR        Clock control & status register.
//  0x80: SSCGR      Spread spectrum clock generation register.
//  0x84: PLLI2SCFGR PLLI2S configuration register.
package rcc

const (
	LSEON  BDCR_Bits = 1 << 0  // External low-speed oscillator enable.
	LSERDY BDCR_Bits = 1 << 1  // External low-speed oscillator ready.
	LSEBYP BDCR_Bits = 1 << 2  // External low-speed oscillator bypass.
	RTCSEL BDCR_Bits = 3 << 8  // RTC clock source selection.
	RTCEN  BDCR_Bits = 1 << 15 // RTC clock enable.
	BDRST  BDCR_Bits = 1 << 16 // Backup domain software reset.
)
