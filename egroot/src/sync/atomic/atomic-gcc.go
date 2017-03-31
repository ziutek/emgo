// +build !cortexm0

package atomic

import "unsafe"

func compareAndSwapUint32(addr *uint32, old, new uint32) (swapped bool)

func compareAndSwapInt(addr *int, old, new int) (swapped bool)

func compareAndSwapUintptr(addr *uintptr, old, new uintptr) (swapped bool)

func compareAndSwapPointer(addr *unsafe.Pointer, old, new unsafe.Pointer) (swapped bool)

func swapUint32(addr *uint32, new uint32) (old uint32)

func swapInt(addr *int, new int) (old int)

func swapUintptr(addr *uintptr, new uintptr) (old uintptr)

func swapPointer(addr *unsafe.Pointer, new unsafe.Pointer) (old unsafe.Pointer)

func addUint32(addr *uint32, delta uint32) (new uint32)

func addInt(addr *int, delta int) (new int)

func addUintptr(addr *uintptr, delta uintptr) (new uintptr)

func orUint32(addr *uint32, mask uint32) (new uint32)

func orUintptr(addr *uintptr, mask uintptr) (new uintptr)

func andUint32(addr *uint32, mask uint32) (new uint32)

func andUintptr(addr *uintptr, mask uintptr) (new uintptr)

func xorUint32(addr *uint32, mask uint32) (new uint32)

func xorUintptr(addr *uintptr, mask uintptr) (new uintptr)
