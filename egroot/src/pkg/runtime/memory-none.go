// +build none

package runtime

import "unsafe"

func freeStart() uintptr

func freeEnd() uintptr

func freeSize() uintptr

func HeapSize() uintptr

func panicMemory() {
	panic("not enough memory")
}

func setSlice(sptr, addr unsafe.Pointer, len, cap uint)