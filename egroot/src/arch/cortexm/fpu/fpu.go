// +build cortexm4f
package fpu

import "mmio"

var cpac = mmio.PtrReg32(0xe000ed88)

type Perm uint32

const (
	Deny Perm = 0 << 20
	Priv Perm = 1 << 20
	Full Perm = 3 << 20
)

func SetAccess(p Perm) {
	cpac.StoreBits(uint32(p), 3<<20)
}

func Access() Perm {
	return Perm(cpac.LoadBits(3 << 20))
}

var fpcc = mmio.PtrReg32(0xe000ef34)

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

func SetSP(f SPFlags) {
	fpcc.Store(uint32(f))
}

func SP() SPFlags {
	return SPFlags(fpcc.Load())
}

var fpca = mmio.PtrReg32(0xe000ef38)
