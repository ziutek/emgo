// Package rcc gives an access to STM32L15x reset and clock control registers.
//
// Peripheral: Ctrl
// Instances:
// 	RCC  0x40023800  AHB
// Registers:
//  0x00  CR         Clock control register.
//  0x04  ICSCR      Internal clock sources calibration register.
//  0x08  CFGR       Clock configuration register.
//  0x0C  CIR        Clock interrupt register.
//  0x10  AHBRSTR    AHB peripheral reset register.
//  0x14  APB2RSTR   APB2 peripheral reset register.
//  0x18  APB1RSTR   APB1 peripheral reset register.
//  0x1C  AHBENR     AHB1 peripheral clock enable register.
//  0x20  APB2ENR    APB2 peripheral clock enable register.
//  0x24  APB1ENR    APB1 peripheral clock enable register.
//  0x28  AHBLPENR   AHB1 peripheral clock enable in low power mode register.
//  0x2C  APB2LPENR  APB2 peripheral clock enabled in low power mode register.
//  0x30  APB1LPENR  APB1 peripheral clock enabled in low power mode register.
//  0x34  CSR        Clock control & status register.
package rcc

const (
	LSION    CSR_Bits = 1 << 0  // Internal low-speed oscillator enable
	LSIRDY   CSR_Bits = 1 << 1  // Internal low-speed oscillator ready
	LSEON    CSR_Bits = 1 << 8  // External low-speed oscillator enable
	LSERDY   CSR_Bits = 1 << 9  // External low-speed oscillator ready
	LSEBYP   CSR_Bits = 1 << 10 // External low-speed oscillator bypass
	LSECSSON CSR_Bits = 1 << 11 // CSS on LSE enable
	LSECSSD  CSR_Bits = 1 << 12 // CSS on LSE failure Detection
	RTCSEL   CSR_Bits = 3 << 16 // RTC and LCD clock source selection
	RTCEN    CSR_Bits = 1 << 22 // RTC clock enable
	RTCRST   CSR_Bits = 1 << 23 // RTC software reset
	RMVF     CSR_Bits = 1 << 24 // Remove reset flag
	OBLRSTF  CSR_Bits = 1 << 25 // Options bytes loading reset flag
	PINRSTF  CSR_Bits = 1 << 26 // PIN reset flag
	PORRSTF  CSR_Bits = 1 << 27 // POR/PDR reset flag
	SFTRSTF  CSR_Bits = 1 << 28 // Software reset flag
	IWDGRSTF CSR_Bits = 1 << 29 // Independent watchdog reset flag
	WWDGRSTF CSR_Bits = 1 << 30 // Window watchdog reset flag
	LPWRRSTF CSR_Bits = 1 << 31 // Low-power reset flag

	RTCSEL_LSE = 1 << 16
	RTCSEL_LSI = 2 << 16
	RTCSEL_HSE = 3 << 16
)
