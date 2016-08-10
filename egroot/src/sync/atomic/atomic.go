// Package atomic provides low-level atomic memory primitives.
//
// These functions works in sequentially consistent memory model so they
// provide "happens-before" edges for all load/store operations.
package atomic

import "unsafe"

func CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool) {
	return compareAndSwapInt32(addr, old, new)
}
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

func AddInt32(addr *int32, delta int32) (new int32) {
	return addInt32(addr, delta)
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

func OrInt32(addr *int32, mask int32) (new int32) {
	return orInt32(addr, mask)
}
func OrUint32(addr *uint32, mask uint32) (new uint32) {
	return orUint32(addr, mask)
}
func OrUintptr(addr *uintptr, mask uintptr) (new uintptr) {
	return orUintptr(addr, mask)
}

func AndInt32(addr *int32, mask int32) (new int32) {
	return andInt32(addr, mask)
}
func AndUint32(addr *uint32, mask uint32) (new uint32) {
	return andUint32(addr, mask)
}
func AndUintptr(addr *uintptr, mask uintptr) (new uintptr) {
	return andUintptr(addr, mask)
}

func XorInt32(addr *int32, mask int32) (new int32) {
	return xorInt32(addr, mask)
}
func XorUint32(addr *uint32, mask uint32) (new uint32) {
	return xorUint32(addr, mask)
}
func XorUintptr(addr *uintptr, mask uintptr) (new uintptr) {
	return xorUintptr(addr, mask)
}

func SwapInt32(addr *int32, new int32) (old int32) {
	return swapInt32(addr, new)
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

func LoadInt32(addr *int32) (val int32) {
	return loadInt32(addr)
}
func LoadUint32(addr *uint32) (val uint32) {
	return loadUint32(addr)
}
func LoadInt(addr *int) (val int) {
	return loadInt(addr)
}
func LoadUintptr(addr *uintptr) (val uintptr) {
	return loadUintptr(addr)
}
func LoadPointer(addr *unsafe.Pointer) (val unsafe.Pointer) {
	return loadPointer(addr)
}

func StoreInt32(addr *int32, val int32) {
	storeInt32(addr, val)
}
func StoreUint32(addr *uint32, val uint32) {
	storeUint32(addr, val)
}
func StoreInt(addr *int, val int) {
	storeInt(addr, val)
}
func StoreUintptr(addr *uintptr, val uintptr) {
	storeUintptr(addr, val)
}
func StorePointer(addr *unsafe.Pointer, val unsafe.Pointer) {
	storePointer(addr, val)
}
