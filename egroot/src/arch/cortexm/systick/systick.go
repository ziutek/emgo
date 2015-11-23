// Package systick gives an access to System Timer registers.
//
// Detailed description of all registers covered by this package can be found in
// "Cortex-M[0-4] Devices Generic User Guide", chapter 4 "Cortex-M[0-4]
// Peripherals".
package systick

const (
	base   = 0xe000e010
	length = 4
)

const (
	CSR Reg = 0 // Any read of CSR clears COUNTFLAG.

	ENABLE    Mask = 1 << 0  // Enable counter.
	TICKINT   Mask = 1 << 1  // Generate exceptions.
	CLKSOURCE Mask = 1 << 2  // Clock source: 0:external, 1:CPU.
	COUNTFLAG Mask = 1 << 16 // 1:Timer counted to 0 since last register read.
)

const (
	RVR Reg = 1

	RELOAD Mask = 1<<24 - 1
)

const (
	CVR Reg = 2

	CURRENT Mask = 1<<24 - 1
)

const (
	CALIB Reg = 3

	TENMS Mask = 1<<24 - 1
	SKEW  Mask = 1 << 30
	NOREF Mask = 1 << 31
)
