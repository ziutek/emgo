package cortexm

import "unsafe"

//c:static inline
func APSR() uint32

//c:static inline
func SetAPSR(r uint32)

//c:static inline
func IPSR() uint32

//c:static inline
func SetIPSR(r uint32)

//c:static inline
func EPSR() uint32

//c:static inline
func SetEPSR(r uint32)

//c:static inline
func IEPSR() uint32

//c:static inline
func IAPSR() uint32

//c:static inline
func SetIAPSR(r uint32)

//c:static inline
func EAPSR() uint32

//c:static inline
func SetEAPSR(r uint32)

//c:static inline
func PSR() uint32

//c:static inline
func SetPSR(r uint32)

//c:static inline
func MSP() uintptr

//c:static inline
func SetMSP(p unsafe.Pointer)

//c:static inline
func PSP() uintptr

//c:static inline
func SetPSP(p unsafe.Pointer)

//c:static inline
func LR() uint32

//c:static inline
func SetLR(r uint32)

type Control uint32

const (
	Unpriv Control = 1 << iota // Unprivileged thread mode
	UsePSP                     // Use PSP for in thread mode
	FPCA                       // Floating-point context active
)

//c:static inline
func SetCtrl(c Control)

//c:static inline
func Ctrl() Control

//c:static inline
func SEV()

//c:static inline
func ISB()

func SVC(imm byte)

func BKPT(imm byte)
