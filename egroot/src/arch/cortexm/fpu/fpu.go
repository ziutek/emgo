// Package fpu gives an access to Floating Point Unit registers.
// Detailed description of all registers covered by this package can be found in
// "Cortex-M[0-4] Devices Generic User Guide", chapter 4 "Cortex-M[0-4]
// Peripherals".
//
// BaseAddr: 0xe000ef34
//  0: FPCCR  Floating-point Context Control Register
//  1: FPCAR  Floating-point Context Address Register
//  2: FPDSCR Floating-point Default Status Control Register
package fpu

const (
	LSPACT FPCCR_Bits = 1 << 0
	USER   FPCCR_Bits = 1 << 1
	THREAD FPCCR_Bits = 1 << 3
	HFRDY  FPCCR_Bits = 1 << 4
	MMRDY  FPCCR_Bits = 1 << 5
	BFRDY  FPCCR_Bits = 1 << 6
	MONRDY FPCCR_Bits = 1 << 8
	LSPEN  FPCCR_Bits = 1 << 30
	ASPEN  FPCCR_Bits = 1 << 31
)

const (
	ADDRESS FPCAR_Bits = 0x3fffffff << 2
)

const (
	RMode FPDSCR_Bits = 3 << 22
	FZ    FPDSCR_Bits = 1 << 24
	DN    FPDSCR_Bits = 1 << 25
	AHP   FPDSCR_Bits = 1 << 26
)
