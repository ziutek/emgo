// +build cortexm4f
package fpu

import (
	"mmio"
	"unsafe"
)

type Perm uint32

const (
	Deny Perm = 0 << 20
	Priv Perm = 1 << 20
	Full Perm = 3 << 20
)

const cpacaddr uintptr = 0xe000ed88

func SetAccess(p Perm) {
	cpac := mmio.PtrU32(unsafe.Pointer(cpacaddr))
	cpac.StoreBits(uint32(p), 3<<20)
}

func Access() Perm {
	cpac := mmio.PtrU32(unsafe.Pointer(cpacaddr))
	return Perm(cpac.LoadBits(3 << 20))
}

// SPFlags control/describe FPU state preservation behavior during exception
// handling.
type SPFlags uint32

const (
	DeferSP SPFlags = 1 << iota
	User
	_
	Thread
	HardFault
	MemFault
	BusFault
	_
	DebugMon

	LazySP SPFlags = 1 << 30
	AutoSP SPFlags = 1 << 31
)

const fpccaddr uintptr = 0xe000ef34

func SetSP(f SPFlags) {
	fpcc := mmio.PtrU32(unsafe.Pointer(fpccaddr))
	fpcc.Store(uint32(f))
}

func SP() SPFlags {
	fpcc := mmio.PtrU32(unsafe.Pointer(fpccaddr))
	return SPFlags(fpcc.Load())
}

const fpcaaddr uintptr = 0xe000ef38
