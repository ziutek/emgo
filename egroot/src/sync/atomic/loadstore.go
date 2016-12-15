package atomic

import "unsafe"

//c:static inline __attribute__((always_inline))
func LoadUint32(addr *uint32) uint32

//c:static inline __attribute__((always_inline))
func LoadInt(addr *int) int

//c:static inline __attribute__((always_inline))
func LoadUintptr(addr *uintptr) uintptr

//c:static inline __attribute__((always_inline))
func LoadPointer(addr *unsafe.Pointer) unsafe.Pointer

//c:static inline __attribute__((always_inline))
func StoreUint32(addr *uint32, val uint32)

//c:static inline __attribute__((always_inline))
func StoreInt(addr *int, val int)

//c:static inline __attribute__((always_inline))
func StoreUintptr(addr *uintptr, val uintptr)

//c:static inline __attribute__((always_inline))
func StorePointer(addr *unsafe.Pointer, val unsafe.Pointer)
