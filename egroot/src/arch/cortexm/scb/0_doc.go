// Package scb gives an access to registers of System Control Block.
//
// Detailed description of all registers covered by this package can be found in
// "Cortex-M[0-4] Devices Generic User Guide", chapter 4 "Cortex-M[0-4]
// Peripherals".
//
// Notes
//
// 1. Cortex-M0 doesn't implement ACTLR, SHPR1, CFSR, HFSR, MMFR, BFAR, AFSR
// registers.
//
// 2. Cortex-M0 supports only word access to SHPR2, SHPR3 so this package does
// not provide byte access to individual fields.
package scb