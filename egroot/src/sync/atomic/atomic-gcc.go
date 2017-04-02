// +build !cortexm0

package atomic

import "unsafe"

//c:inline
func compareAndSwapUint32(addr *uint32, old, new uint32) (swapped bool)

//c:inline
func compareAndSwapInt(addr *int, old, new int) (swapped bool)

//c:inline
func compareAndSwapUintptr(addr *uintptr, old, new uintptr) (swapped bool)

//c:inline
func compareAndSwapPointer(addr *unsafe.Pointer, old, new unsafe.Pointer) (swapped bool)

//c:inline
func swapUint32(addr *uint32, new uint32) (old uint32)

//c:inline
func swapInt(addr *int, new int) (old int)

//c:inline
func swapUintptr(addr *uintptr, new uintptr) (old uintptr)

//c:inline
func swapPointer(addr *unsafe.Pointer, new unsafe.Pointer) (old unsafe.Pointer)

//c:inline
func addUint32(addr *uint32, delta uint32) (new uint32)

//c:inline
func addInt(addr *int, delta int) (new int)

//c:inline
func addUintptr(addr *uintptr, delta uintptr) (new uintptr)

//c:inline
func orUint32(addr *uint32, mask uint32) (new uint32)

//c:inline
func orUintptr(addr *uintptr, mask uintptr) (new uintptr)

//c:inline
func andUint32(addr *uint32, mask uint32) (new uint32)

//c:inline
func andUintptr(addr *uintptr, mask uintptr) (new uintptr)

//c:inline
func xorUint32(addr *uint32, mask uint32) (new uint32)

//c:inline
func xorUintptr(addr *uintptr, mask uintptr) (new uintptr)
