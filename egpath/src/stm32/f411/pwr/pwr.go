// Package pwr gives an access to STM32F411xC/E power controller registers.
//
// BaseAddr: 0x40007000  APB1
//  0x0: CR  Power control register.
//  0x4: CSR Power control/status register.
package pwr

const (
	LPDS   CR_Bits = 1 << 0  // Low-power deepsleep
	PDDS   CR_Bits = 1 << 1  // Power-down deepsleep
	CWUF   CR_Bits = 1 << 2  // Clear wakeup flag
	CSBF   CR_Bits = 1 << 3  // Clear standby flag
	PVDE   CR_Bits = 1 << 4  // Power voltage detector enable
	PLS    CR_Bits = 7 << 5  // PVD level selection
	DBP    CR_Bits = 1 << 8  // Disable backup domain write protection
	FPDS   CR_Bits = 1 << 9  // Flash power-down in Stop mode
	LPLVDS CR_Bits = 1 << 10 // Low-power regulator Low Voltage in Deep Sleep
	MRLVDS CR_Bits = 1 << 11 // Main regulator Low Voltage in Deep Sleep
	ADCDC1 CR_Bits = 1 << 13 //
	VOS    CR_Bits = 3 << 14 // Regulator voltage scaling output selection
	FMSSR  CR_Bits = 1 << 20 // Flash Memory Sleep System Run.
	FISSR  CR_Bits = 1 << 21 // Flash Interface Stop while System Run
)

const (
	WUF    CSR_Bits = 1 << 0 // Wakeup flag
	SBF    CSR_Bits = 1 << 0 // Standby flag
	PVDO   CSR_Bits = 1 << 0 // PVD output
	BRR    CSR_Bits = 1 << 0 // Backup regulator ready
	EWUP   CSR_Bits = 1 << 0 // Enable WKUP pin
	BRE    CSR_Bits = 1 << 0 // Backup regulator enable
	VOSRDY CSR_Bits = 1 << 0 // Regulator voltage scaling output selection ready bit
)
