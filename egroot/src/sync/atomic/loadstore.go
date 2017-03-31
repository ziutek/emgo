package atomic

import "unsafe"

func LoadUint32(addr *uint32) uint32

func LoadInt(addr *int) int

func LoadUintptr(addr *uintptr) uintptr

func LoadPointer(addr *unsafe.Pointer) unsafe.Pointer

func StoreUint32(addr *uint32, val uint32)

func StoreInt(addr *int, val int)

func StoreUintptr(addr *uintptr, val uintptr)

func StorePointer(addr *unsafe.Pointer, val unsafe.Pointer)
