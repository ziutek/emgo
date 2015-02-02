// +build cortexm0

package atomic

import (
	"unsafe"

	"arch/cortexm/exce"
)

func compareAndSwapInt32(addr *int32, old, new int32) (swapped bool) {
	exce.DisablePri()
	if swapped = (*addr == old); swapped {
		*addr = new
	}
	exce.EnablePri()
	return
}

func compareAndSwapUint32(addr *uint32, old, new uint32) (swapped bool) {
	exce.DisablePri()
	if swapped = (*addr == old); swapped {
		*addr = new
	}
	exce.EnablePri()
	return
}

func compareAndSwapUintptr(addr *uintptr, old, new uintptr) (swapped bool) {
	exce.DisablePri()
	if swapped = (*addr == old); swapped {
		*addr = new
	}
	exce.EnablePri()
	return
}

func compareAndSwapPointer(addr *unsafe.Pointer, old, new unsafe.Pointer) (swapped bool) {
	exce.DisablePri()
	if swapped = (*addr == old); swapped {
		*addr = new
	}
	exce.EnablePri()
	return
}

func addInt32(addr *int32, delta int32) (new int32) {
	exce.DisablePri()
	new = *addr + delta
	*addr = new
	exce.EnablePri()
	return
}

func addUint32(addr *uint32, delta uint32) (new uint32) {
	exce.DisablePri()
	new = *addr + delta
	*addr = new
	exce.EnablePri()
	return
}

func addUintptr(addr *uintptr, delta uintptr) (new uintptr) {
	exce.DisablePri()
	new = *addr + delta
	*addr = new
	exce.EnablePri()
	return
}

func orInt32(addr *int32, mask int32) (new int32) {
	exce.DisablePri()
	new = *addr | mask
	*addr = new
	exce.EnablePri()
	return
}

func orUint32(addr *uint32, mask uint32) (new uint32) {
	exce.DisablePri()
	new = *addr | mask
	*addr = new
	exce.EnablePri()
	return
}

func orUintptr(addr *uintptr, mask uintptr) (new uintptr) {
	exce.DisablePri()
	new = *addr | mask
	*addr = new
	exce.EnablePri()
	return
}

func andInt32(addr *int32, mask int32) (new int32) {
	exce.DisablePri()
	new = *addr & mask
	*addr = new
	exce.EnablePri()
	return
}

func andUint32(addr *uint32, mask uint32) (new uint32) {
	exce.DisablePri()
	new = *addr & mask
	*addr = new
	exce.EnablePri()
	return
}

func andUintptr(addr *uintptr, mask uintptr) (new uintptr) {
	exce.DisablePri()
	new = *addr & mask
	*addr = new
	exce.EnablePri()
	return
}

func xorInt32(addr *int32, mask int32) (new int32) {
	exce.DisablePri()
	new = *addr ^ mask
	*addr = new
	exce.EnablePri()
	return
}

func xorUint32(addr *uint32, mask uint32) (new uint32) {
	exce.DisablePri()
	new = *addr ^ mask
	*addr = new
	exce.EnablePri()
	return
}

func xorUintptr(addr *uintptr, mask uintptr) (new uintptr) {
	exce.DisablePri()
	new = *addr ^ mask
	*addr = new
	exce.EnablePri()
	return
}

func swapInt32(addr *int32, new int32) (old int32) {
	exce.DisablePri()
	old = *addr
	*addr = new
	exce.EnablePri()
	return
}

func swapUint32(addr *uint32, new uint32) (old uint32) {
	exce.DisablePri()
	old = *addr
	*addr = new
	exce.EnablePri()
	return
}

func swapUintptr(addr *uintptr, new uintptr) (old uintptr) {
	exce.DisablePri()
	old = *addr
	*addr = new
	exce.EnablePri()
	return
}

func swapPointer(addr *unsafe.Pointer, new unsafe.Pointer) (old unsafe.Pointer) {
	exce.DisablePri()
	old = *addr
	*addr = new
	exce.EnablePri()
	return
}

func loadInt32(addr *int32) int32 {
	return *addr
}

func loadUint32(addr *uint32) uint32 {
	return *addr
}

func loadUintptr(addr *uintptr) uintptr {
	return *addr
}

func loadPointer(addr *unsafe.Pointer) unsafe.Pointer {
	return *addr
}

func storeInt32(addr *int32, val int32) {
	*addr = val
}

func storeUint32(addr *uint32, val uint32) {
	*addr = val
}

func storeUintptr(addr *uintptr, val uintptr) {
	*addr = val
}

func storePointer(addr *unsafe.Pointer, val unsafe.Pointer) {
	*addr = val
}
