// Package fpu gives an access to Floating Point Unit registers.
// Detailed description of all registers covered by this package can be found in
// "Cortex-M[0-4] Devices Generic User Guide", chapter 4 "Cortex-M[0-4]
// Peripherals".
//
// Peripheral: FPU_Periph  Floating Point Unit
// Instances:
//  FPU  0xe000ED88
// Registers:
//  0x000 32  CPACR   Coprocessor Access Control Register
//  0x1AC 32  FPCCR   Floating-point Context Control Register
//  0x1B0 32  FPCAR   Floating-point Context Address Register
//  0x1B4 32  FPDSCR  Floating-point Default Status Control Register
package fpu

const (
	CP10 CPACR = 3 << 20 //+ Access privileges for coprocessor 10.
	CP11 CPACR = 3 << 22 //+ Access privileges for coprocessor 11.

	CPACDENY CPACR = 0
	CPACPRIV CPACR = 1
	CPACFULL CPACR = 3
)

const (
	CP10n = 20
	CP11n = 22
)

const (
	LSPACT FPCCR = 1 << 0  //+
	USER   FPCCR = 1 << 1  //+
	THREAD FPCCR = 1 << 3  //+
	HFRDY  FPCCR = 1 << 4  //+
	MMRDY  FPCCR = 1 << 5  //+
	BFRDY  FPCCR = 1 << 6  //+
	MONRDY FPCCR = 1 << 8  //+
	LSPEN  FPCCR = 1 << 30 //+
	ASPEN  FPCCR = 1 << 31 //+
)

const (
	ADDRESS FPCAR = 0x3fffffff << 2 //+
)

const (
	RMode FPDSCR = 3 << 22 //+
	FZ    FPDSCR = 1 << 24 //+
	DN    FPDSCR = 1 << 25 //+
	AHP   FPDSCR = 1 << 26 //+
)
