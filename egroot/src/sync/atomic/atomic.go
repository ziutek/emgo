// Package atomic provides low-level atomic memory primitives.
package atomic

import "unsafe"

func CompareAndSwapUint32(addr *uint32, old, new uint32) (swapped bool) {
	return compareAndSwapUint32(addr, old, new)
}
func CompareAndSwapInt(addr *int, old, new int) (swapped bool) {
	return compareAndSwapInt(addr, old, new)
}
func CompareAndSwapUintptr(addr *uintptr, old, new uintptr) (swapped bool) {
	return compareAndSwapUintptr(addr, old, new)
}
func CompareAndSwapPointer(addr *unsafe.Pointer, old, new unsafe.Pointer) (swapped bool) {
	return compareAndSwapPointer(addr, old, new)
}

func SwapUint32(addr *uint32, new uint32) (old uint32) {
	return swapUint32(addr, new)
}
func SwapInt(addr *int, new int) (old int) {
	return swapInt(addr, new)
}
func SwapUintptr(addr *uintptr, new uintptr) (old uintptr) {
	return swapUintptr(addr, new)
}
func SwapPointer(addr *unsafe.Pointer, new unsafe.Pointer) (old unsafe.Pointer) {
	return swapPointer(addr, new)
}

func AddUint32(addr *uint32, delta uint32) (new uint32) {
	return addUint32(addr, delta)
}
func AddInt(addr *int, delta int) (new int) {
	return addInt(addr, delta)
}
func AddUintptr(addr *uintptr, delta uintptr) (new uintptr) {
	return addUintptr(addr, delta)
}

func OrUint32(addr *uint32, mask uint32) (new uint32) {
	return orUint32(addr, mask)
}
func OrUintptr(addr *uintptr, mask uintptr) (new uintptr) {
	return orUintptr(addr, mask)
}

func AndUint32(addr *uint32, mask uint32) (new uint32) {
	return andUint32(addr, mask)
}
func AndUintptr(addr *uintptr, mask uintptr) (new uintptr) {
	return andUintptr(addr, mask)
}

func XorUint32(addr *uint32, mask uint32) (new uint32) {
	return xorUint32(addr, mask)
}
func XorUintptr(addr *uintptr, mask uintptr) (new uintptr) {
	return xorUintptr(addr, mask)
}
