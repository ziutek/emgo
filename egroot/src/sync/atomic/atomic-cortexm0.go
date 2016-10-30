// +build cortexm0

package atomic

import (
	"sync/fence"
	"unsafe"

	"arch/cortexm"
)

func compareAndSwapInt32(addr *int32, old, new int32) (swapped bool) {
	cortexm.SetPRIMASK()
	if swapped = (*addr == old); swapped {
		*addr = new
	}
	cortexm.ClearPRIMASK()
	return
}

func compareAndSwapUint32(addr *uint32, old, new uint32) (swapped bool) {
	cortexm.SetPRIMASK()
	if swapped = (*addr == old); swapped {
		*addr = new
	}
	cortexm.ClearPRIMASK()
	return
}

func compareAndSwapInt(addr *int, old, new int) (swapped bool) {
	cortexm.SetPRIMASK()
	if swapped = (*addr == old); swapped {
		*addr = new
	}
	cortexm.ClearPRIMASK()
	return
}

func compareAndSwapUintptr(addr *uintptr, old, new uintptr) (swapped bool) {
	cortexm.SetPRIMASK()
	if swapped = (*addr == old); swapped {
		*addr = new
	}
	cortexm.ClearPRIMASK()
	return
}

func compareAndSwapPointer(addr *unsafe.Pointer, old, new unsafe.Pointer) (swapped bool) {
	cortexm.SetPRIMASK()
	if swapped = (*addr == old); swapped {
		*addr = new
	}
	cortexm.ClearPRIMASK()
	return
}

func addInt32(addr *int32, delta int32) (new int32) {
	cortexm.SetPRIMASK()
	new = *addr + delta
	*addr = new
	cortexm.ClearPRIMASK()
	return
}

func addUint32(addr *uint32, delta uint32) (new uint32) {
	cortexm.SetPRIMASK()
	new = *addr + delta
	*addr = new
	cortexm.ClearPRIMASK()
	return
}

func addInt(addr *int, delta int) (new int) {
	cortexm.SetPRIMASK()
	new = *addr + delta
	*addr = new
	cortexm.ClearPRIMASK()
	return
}

func addUintptr(addr *uintptr, delta uintptr) (new uintptr) {
	cortexm.SetPRIMASK()
	new = *addr + delta
	*addr = new
	cortexm.ClearPRIMASK()
	return
}

func orInt32(addr *int32, mask int32) (new int32) {
	cortexm.SetPRIMASK()
	new = *addr | mask
	*addr = new
	cortexm.ClearPRIMASK()
	return
}

func orUint32(addr *uint32, mask uint32) (new uint32) {
	cortexm.SetPRIMASK()
	new = *addr | mask
	*addr = new
	cortexm.ClearPRIMASK()
	return
}

func orUintptr(addr *uintptr, mask uintptr) (new uintptr) {
	cortexm.SetPRIMASK()
	new = *addr | mask
	*addr = new
	cortexm.ClearPRIMASK()
	return
}

func andInt32(addr *int32, mask int32) (new int32) {
	cortexm.SetPRIMASK()
	new = *addr & mask
	*addr = new
	cortexm.ClearPRIMASK()
	return
}

func andUint32(addr *uint32, mask uint32) (new uint32) {
	cortexm.SetPRIMASK()
	new = *addr & mask
	*addr = new
	cortexm.ClearPRIMASK()
	return
}

func andUintptr(addr *uintptr, mask uintptr) (new uintptr) {
	cortexm.SetPRIMASK()
	new = *addr & mask
	*addr = new
	cortexm.ClearPRIMASK()
	return
}

func xorInt32(addr *int32, mask int32) (new int32) {
	cortexm.SetPRIMASK()
	new = *addr ^ mask
	*addr = new
	cortexm.ClearPRIMASK()
	return
}

func xorUint32(addr *uint32, mask uint32) (new uint32) {
	cortexm.SetPRIMASK()
	new = *addr ^ mask
	*addr = new
	cortexm.ClearPRIMASK()
	return
}

func xorUintptr(addr *uintptr, mask uintptr) (new uintptr) {
	cortexm.SetPRIMASK()
	new = *addr ^ mask
	*addr = new
	cortexm.ClearPRIMASK()
	return
}

func swapInt32(addr *int32, new int32) (old int32) {
	cortexm.SetPRIMASK()
	old = *addr
	*addr = new
	cortexm.ClearPRIMASK()
	return
}

func swapInt(addr *int, new int) (old int) {
	cortexm.SetPRIMASK()
	old = *addr
	*addr = new
	cortexm.ClearPRIMASK()
	return
}

func swapUint32(addr *uint32, new uint32) (old uint32) {
	cortexm.SetPRIMASK()
	old = *addr
	*addr = new
	cortexm.ClearPRIMASK()
	return
}

func swapUintptr(addr *uintptr, new uintptr) (old uintptr) {
	cortexm.SetPRIMASK()
	old = *addr
	*addr = new
	cortexm.ClearPRIMASK()
	return
}

func swapPointer(addr *unsafe.Pointer, new unsafe.Pointer) (old unsafe.Pointer) {
	cortexm.SetPRIMASK()
	old = *addr
	*addr = new
	cortexm.ClearPRIMASK()
	return
}

func loadInt32(addr *int32) int32 {
	fence.Compiler()
	return *addr
}

func loadUint32(addr *uint32) uint32 {
	fence.Compiler()
	return *addr
}

func loadInt(addr *int) int {
	fence.Compiler()
	return *addr
}

func loadUintptr(addr *uintptr) uintptr {
	fence.Compiler()
	return *addr
}

func loadPointer(addr *unsafe.Pointer) unsafe.Pointer {
	fence.Compiler()
	return *addr
}

func storeInt32(addr *int32, val int32) {
	*addr = val
	fence.Compiler()
}

func storeUint32(addr *uint32, val uint32) {
	*addr = val
	fence.Compiler()
}

func storeInt(addr *int, val int) {
	*addr = val
	fence.Compiler()
}

func storeUintptr(addr *uintptr, val uintptr) {
	*addr = val
	fence.Compiler()
}

func storePointer(addr *unsafe.Pointer, val unsafe.Pointer) {
	*addr = val
	fence.Compiler()
}
