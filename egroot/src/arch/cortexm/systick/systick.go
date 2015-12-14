// Package systick gives an access to System Timer registers.
//
// Detailed description of all registers covered by this package can be found in
// "Cortex-M[0-4] Devices Generic User Guide", chapter 4 "Cortex-M[0-4]
// Peripherals".
//
// BaseAddr: 0xE000E010
//	0x00: CSR   SysTick Control and Status Register (any read clears COUNTFLAG).
//	0x04: RVR   SysTick Reload Value Register.
//	0x08: CVR   SysTick Current Value Register.
//	0x0C: CALIB SysTick Calibration Value Register.
package systick

const (
	ENABLE    CSR_Bits = 1 << 0  // Enable counter.
	TICKINT   CSR_Bits = 1 << 1  // Generate exceptions.
	CLKSOURCE CSR_Bits = 1 << 2  // Clock source: 0:external, 1:CPU.
	COUNTFLAG CSR_Bits = 1 << 16 // 1:Timer counted to 0 since last read of CSR.
)

const (
	RELOAD RVR_Bits = 1<<24 - 1
)

const (
	CURRENT CVR_Bits = 1<<24 - 1
)

const (
	TENMS CALIB_Bits = 1<<24 - 1
	SKEW  CALIB_Bits = 1 << 30
	NOREF CALIB_Bits = 1 << 31
)
