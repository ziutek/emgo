package atomic

import "unsafe"

//c:static inline
func LoadUint32(addr *uint32) uint32

//c:static inline
func LoadInt(addr *int) int

//c:static inline
func LoadUintptr(addr *uintptr) uintptr

//c:static inline
func LoadPointer(addr *unsafe.Pointer) unsafe.Pointer

//c:static inline
func StoreUint32(addr *uint32, val uint32)

//c:static inline
func StoreInt(addr *int, val int)

//c:static inline
func StoreUintptr(addr *uintptr, val uintptr)

//c:static inline
func StorePointer(addr *unsafe.Pointer, val unsafe.Pointer)
