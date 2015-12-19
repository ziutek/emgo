// Package pwr provides interface to STM32L15x power controller registers.
//
// Peripheral: Ctrl
// Instances:
//  PWR  0x40007000  APB1
// Registers:
//  0x0  CR   Power control register.
//  0x4  CSR  Power control/status register.
package pwr

const (
	LPSDSR CR_Bits = 1 << 0  // Low-power deepsleep/sleep/low power run
	PDDS   CR_Bits = 1 << 1  // Power-down deepsleep
	CWUF   CR_Bits = 1 << 2  // Clear wakeup flag
	CSBF   CR_Bits = 1 << 3  // Clear standby flag
	PVDE   CR_Bits = 1 << 4  // Power voltage detector enable
	PLS    CR_Bits = 7 << 5  // PVD level selection
	DBP    CR_Bits = 1 << 8  // Disable backup domain write protection
	ULP    CR_Bits = 1 << 9  // Ultralow power mode
	FWU    CR_Bits = 1 << 10 // Fast wakeup
	VOS    CR_Bits = 3 << 11 // Voltage scaling range selection
	LPRUN  CR_Bits = 1 << 14 // Low power run mode
)

const (
	WUF         CSR_Bits = 1 << 0  // Wakeup flag
	SBF         CSR_Bits = 1 << 1  // Standby flag
	PVDO        CSR_Bits = 1 << 2  // PVD output
	VREFINTRDYF CSR_Bits = 1 << 3  // Internal voltage reference (VREFINT) ready flag
	VOSF        CSR_Bits = 1 << 4  // Voltage Scaling select flag
	REGLPF      CSR_Bits = 1 << 5  // Regulator LP flag
	EWUP1       CSR_Bits = 1 << 8  // Enable WKUP pin 1 bit
	EWUP2       CSR_Bits = 1 << 9  // Enable WKUP pin 2
	EWUP3       CSR_Bits = 1 << 10 // Enable WKUP pin 3
)
