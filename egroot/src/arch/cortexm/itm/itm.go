// Package itm provides interface to Instrumentation Trace Macrocell
package itm

import (
	"mmio"
	"unsafe"
)

type regs struct {
	stim [256]mmio.Reg32
	_    [640]uint32
	te   [8]mmio.Reg32
	_    [8]uint32
	tp   mmio.Reg32
	_    [15]uint32
	tc   mmio.Reg32
}

var p = (*regs)(unsafe.Pointer(uintptr(0xe0000000)))

// PrivMask returns conten of Trace Privilege Register. Every bit in returned
// value corresponds to eight stimulus ports. If bit is set then the
// corresponding ports can be accessed by privileged code only.
//
// Bits for unimplemented ports are always returned as 0. To determine the
// number of implemnted ports call SetPrivMask(0xffffffff) and next call
// PrivMask().
func PrivMask() uint32 {
	return p.tp.Load()
}

// SetPrivMask writes mask to Trace Privilege Register.
func StimNum(mask uint32) {
	p.tp.Store(mask)
}
