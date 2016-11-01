// +build !cortexm0

package atomic

import "unsafe"

//c:static inline
func compareAndSwapInt32(addr *int32, old, new int32) (swapped bool)

//c:static inline
func compareAndSwapUint32(addr *uint32, old, new uint32) (swapped bool)

//c:static inline
func compareAndSwapInt(addr *int, old, new int) (swapped bool)

//c:static inline
func compareAndSwapUintptr(addr *uintptr, old, new uintptr) (swapped bool)

//c:static inline
func compareAndSwapPointer(addr *unsafe.Pointer, old, new unsafe.Pointer) (swapped bool)

//c:static inline
func addInt32(addr *int32, delta int32) (new int32)

//c:static inline
func addUint32(addr *uint32, delta uint32) (new uint32)

//c:static inline
func addInt(addr *int, delta int) (new int)

//c:static inline
func addUintptr(addr *uintptr, delta uintptr) (new uintptr)

//c:static inline
func orInt32(addr *int32, mask int32) (new int32)

//c:static inline
func orUint32(addr *uint32, mask uint32) (new uint32)

//c:static inline
func orUintptr(addr *uintptr, mask uintptr) (new uintptr)

//c:static inline
func andInt32(addr *int32, mask int32) (new int32)

//c:static inline
func andUint32(addr *uint32, mask uint32) (new uint32)

//c:static inline
func andUintptr(addr *uintptr, mask uintptr) (new uintptr)

//c:static inline
func xorInt32(addr *int32, mask int32) (new int32)

//c:static inline
func xorUint32(addr *uint32, mask uint32) (new uint32)

//c:static inline
func xorUintptr(addr *uintptr, mask uintptr) (new uintptr)

//c:static inline
func swapInt32(addr *int32, new int32) (old int32)

//c:static inline
func swapUint32(addr *uint32, new uint32) (old uint32)

//c:static inline
func swapInt(addr *int, new int) (old int)

//c:static inline
func swapUintptr(addr *uintptr, new uintptr) (old uintptr)

//c:static inline
func swapPointer(addr *unsafe.Pointer, new unsafe.Pointer) (old unsafe.Pointer)
