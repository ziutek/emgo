package scb

import "unsafe"

var acr = (*wordReg)(unsafe.Pointer(uintptr(0xe000e008)))

func DisableOOFP() {
	acr.setBit(9)
}

func EnableOOFP() {
	acr.clearBit(9)
}
