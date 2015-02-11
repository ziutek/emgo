package debug

import (
	"unsafe"
)

type regs struct {
	dhcs uint32 `C:"volatile"`
	dcrs uint32 `C:"volatile"`
	dcrd uint32 `C:"volatile"`
	demc uint32 `C:"volatile"`
}

var drs = (*regs)(unsafe.Pointer(uintptr(0xe000edf0)))

type DEMC uint32

const (
	CoreResetVC DEMC = 1 << iota
	_
	_
	_
	MMErrVC
	NoCPErrVC
	ChkErrVC
	StatErrVC

	BusErrVC
	IntErrVC
	HardErrVC
	_
	_
	_
	_
	_

	MonEna
	MonPend
	MonStep
	MonReq
	_
	_
	_
	_

	TrcEna
)

// DEMCR returns value of DEMC register.
func DEMCR() DEMC {
	return DEMC(drs.demc)
}

// SetDEMCR returns value of DEMC register.
func SetDEMCR(v DEMC) {
	drs.demc = uint32(v)
}