// +build !cortexm0

package atomic

import "unsafe"

func compareAndSwapInt32(addr *int32, old, new int32) (swapped bool)
func compareAndSwapUint32(addr *uint32, old, new uint32) (swapped bool)
func compareAndSwapUintptr(addr *uintptr, old, new uintptr) (swapped bool)
func compareAndSwapPointer(addr *unsafe.Pointer, old, new unsafe.Pointer) (swapped bool)

func addInt32(addr *int32, delta int32) (new int32)
func addUint32(addr *uint32, delta uint32) (new uint32)
func addUintptr(addr *uintptr, delta uintptr) (new uintptr)

func orInt32(addr *int32, mask int32) (new int32)
func orUint32(addr *uint32, mask uint32) (new uint32)
func orUintptr(addr *uintptr, mask uintptr) (new uintptr)

func andInt32(addr *int32, mask int32) (new int32)
func andUint32(addr *uint32, mask uint32) (new uint32)
func andUintptr(addr *uintptr, mask uintptr) (new uintptr)

func xorInt32(addr *int32, mask int32) (new int32)
func xorUint32(addr *uint32, mask uint32) (new uint32)
func xorUintptr(addr *uintptr, mask uintptr) (new uintptr)

func swapInt32(addr *int32, new int32) (old int32)
func swapUint32(addr *uint32, new uint32) (old uint32)
func swapUintptr(addr *uintptr, new uintptr) (old uintptr)
func swapPointer(addr *unsafe.Pointer, new unsafe.Pointer) (old unsafe.Pointer)

func loadInt32(addr *int32) int32
func loadUint32(addr *uint32) uint32
func loadUintptr(addr *uintptr) uintptr
func loadPointer(addr *unsafe.Pointer) unsafe.Pointer

func storeInt32(addr *int32, val int32)
func storeUint32(addr *uint32, val uint32)
func storeUintptr(addr *uintptr, val uintptr)
func storePointer(addr *unsafe.Pointer, val unsafe.Pointer)
