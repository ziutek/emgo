// Package systick gives an access to System Timer registers.
//
// Detailed description of all registers covered by this package can be found in
// "Cortex-M[0-4] Devices Generic User Guide", chapter 4 "Cortex-M[0-4]
// Peripherals".
//
// Peripheral: SYSTICK_Periph  Cortex-M 24-bit system timer.
// Instances:
//  SYSTICK  0xE000E010
// Registers:
//	0x00 32  CSR   Control and Status Register (any read clears COUNTFLAG)
//	0x04 32  RVR   Reload Value Register.
//	0x08 32  CVR   Current Value Register.
//	0x0C 32  CALIB Calibration Value Register.
package systick

const (
	ENABLE    CSR_Bits = 1 << 0  //+ Enable counter.
	TICKINT   CSR_Bits = 1 << 1  //+ Generate exceptions.
	CLKSOURCE CSR_Bits = 1 << 2  //+ Clock source: 0:external, 1:CPU.
	COUNTFLAG CSR_Bits = 1 << 16 //+ 1:Timer counted to 0 since last read of CSR
)

const (
	RELOAD RVR_Bits = 1<<24 - 1 //+ Loaded into CVR when the counter reaches 0.
)

const (
	CURRENT CVR_Bits = 1<<24 - 1 //+ Read: couter value. Write: clears to zero.
)

const (
	TENMS CALIB_Bits = 1<<24 - 1 //+
	SKEW  CALIB_Bits = 1 << 30   //+
	NOREF CALIB_Bits = 1 << 31   //+
)