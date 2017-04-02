package atomic

import "unsafe"

//c:inline
func LoadUint32(addr *uint32) uint32

//c:inline
func LoadInt(addr *int) int

//c:inline
func LoadUintptr(addr *uintptr) uintptr

//c:inline
func LoadPointer(addr *unsafe.Pointer) unsafe.Pointer

//c:inline
func StoreUint32(addr *uint32, val uint32)

//c:inline
func StoreInt(addr *int, val int)

//c:inline
func StoreUintptr(addr *uintptr, val uintptr)

//c:inline
func StorePointer(addr *unsafe.Pointer, val unsafe.Pointer)
