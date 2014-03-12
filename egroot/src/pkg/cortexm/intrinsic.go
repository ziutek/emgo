package cortexm

import "unsafe"

func APSR() uint32

func SetAPSR(r uint32)

func IPSR() uint32

func SetIPSR(r uint32)

func EPSR() uint32

func SetEPSR(r uint32)

func IEPSR() uint32

func IAPSR() uint32

func SetIAPSR(r uint32)

func EAPSR() uint32

func SetEAPSR(r uint32)

func PSR() uint32

func SetPSR(r uint32)

func MSP() uintptr

func SetMSP(p unsafe.Pointer)

func PSP() uintptr

func SetPSP(p unsafe.Pointer)

type Control uint32

const (
	Unpriv Control = 1 << iota // Unprivileged thread mode
	UsePSP                     // Use PSP for in thread mode
	FPCA                       // Floating-point context active
)

func SetCtrl(c Control)

func Ctrl() Control

func SEV()

func SVC(n byte)
