// Package cortexm provides intrinsic functions for many Cortex-M special
// instructions.
package cortexm

import "unsafe"

// SEV inserts Signal Event instruction.
//
// Implies compiler fence.
//
//c:inline
func SEV()

// DMB inserts Data Memory Barrier instruction.
//
// Implies compiler fence.
//
//c:inline
func DMB()

// DSB inserts Data Synchronization Barrier instruction.
//
// Implies compiler fence.
//
//c:inline
func DSB()

// ISB inserts Instruction Synchronization Barrier instruction.
//
// Implies compiler fence.
//
//c:inline
func ISB()

// WFE inserts Wait For Event instruction.
//
// Implies compiler fence.
//
//c:inline
func WFE()

// WFI inserts Wait For Interrupt instruction.
//
// Implies compiler fence.
//
//c:inline
func WFI()

// SVC inserts Supervisor Call instruction.
//
// Implies compiler fence.
func SVC(imm byte)

// BKPT inserts Breakpoint instruction.
//
// Implies compiler fence.
func BKPT(imm byte)

// PRIMASK returns true if all exceptions with configurable priority are
// disabled.
//
//c:inline
func PRIMASK() bool

// SetPRIMASK disables all exceptions with configurable priority. Internally it
// inserts cpsid i instruction. Atomic primitives on Cortex-M0 always enable
// exceptions after atomic operation. If you need this functions on Cortex-M0
// you should be very careful.
//
// Implies compiler fence.
//
//c:inline
func SetPRIMASK()

// ClearPRIMASK reverts SetPRIMASK. Internally it inserts cpsie i instruction.
// If you modified any data that can be used by enabled interrupt handlers you
// probably need to call fence.Memory() before use this function.
//
// Implies compiler fence.
//
//c:inline
func ClearPRIMASK()

// FAULTMASK returns true if all exceptions other than NMI are disabled.
//
//c:inline
func FAULTMASK() bool

// SetFAULTMASK disables all exceptions other than NMI. Internally it inserts
// cpsid f instruction. Not supported by Cortex-M0.
//
// Implies compiler fence.
//
//c:inline
func SetFAULTMASK()

// ClearFAULTMASK reverts SetFAULTMASK. Internally it inserts cpsie f
// instruction. If you modified any data that can be used by enabled interrupt
// handlers you probably need to call fence.Memory() before use this function.
// Not supported by Cortex-M0.
//
// Implies compiler fence.
//
//c:inline
func ClearFAULTMASK()

// BASEPRIO returns current value of BASEPRI register.
//
//c:inline
func BASEPRI() byte

// SetBASEPRI sets BASEPRI register. It prevents the activation of exceptions
// with the same or lower as p. Not supported by Cortex-M0.
//
// Implies compiler fence.
//
//c:inline
func SetBASEPRI(p byte)

//c:inline
func APSR() uint32

//c:inline
func SetAPSR(r uint32)

//c:inline
func IPSR() uint32

//c:inline
func EPSR() uint32

//c:inline
func IEPSR() uint32

//c:inline
func IAPSR() uint32

//c:inline
func EAPSR() uint32

//c:inline
func PSR() uint32

//c:inline
func SetPSR(r uint32)

//c:inline
func MSP() uintptr

//c:inline
func SetMSP(p unsafe.Pointer)

//c:inline
func PSP() uintptr

//c:inline
func SetPSP(p unsafe.Pointer)

//c:inline
func LR() uint32

//c:inline
func SetLR(r uint32)

type Cflags uint32

const (
	Unpriv Cflags = 1 << 0 // Unprivileged thread mode
	UsePSP Cflags = 1 << 1 // Use PSP for in thread mode
	FPCA   Cflags = 1 << 2 // Floating-point context active
)

//c:inline
func SetCONTROL(c Cflags)

//c:inline
func CONTROL() Cflags
