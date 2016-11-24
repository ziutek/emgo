package noos

import (
	"internal"
	"mem"
	"sync/atomic"
	"sync/fence"
	"unsafe"
)

func panicMemory() {
	panic("out of memory")
}

var (
	aHeapBegin = heapBegin()
	aHeapEnd   = heapEnd()
)

func allocBottom(n int, size, align uintptr) unsafe.Pointer {
	for {
		hb := atomic.LoadUintptr(&aHeapBegin)
		he := atomic.LoadUintptr(&aHeapEnd)
		if hb < heapBegin() || he > heapEnd() || hb > he {
			panicMemory()
		}
		p := mem.AlignUp(hb, align)
		newhb := p + size
		if atomic.CompareAndSwapUintptr(&aHeapBegin, hb, newhb) {
			fence.RW_SMP() // Ensure another load(aHeapEnd) after CAS.
			he := atomic.LoadUintptr(&aHeapEnd)
			if newhb < heapBegin() || newhb > he {
				panicMemory()
			}
			return unsafe.Pointer(p)
		}
	}
}

func allocTop(n int, size, align uintptr) unsafe.Pointer {
	for {
		hb := atomic.LoadUintptr(&aHeapBegin)
		he := atomic.LoadUintptr(&aHeapEnd)
		if hb < heapBegin() || he > heapEnd() || hb > he {
			panicMemory()
		}
		p := mem.AlignDown(he-size, align)
		if atomic.CompareAndSwapUintptr(&aHeapEnd, he, p) {
			fence.RW_SMP()  // Ensure another load(aHeapBegin) after CAS.
			hb := atomic.LoadUintptr(&aHeapBegin)
			if p > heapEnd() || hb > p {
				panicMemory()
			}
			return unsafe.Pointer(p)
		}
	}
}

// alloc is trivial, non-blocking memory allocator.
// For now there is no way to deallocate memory allocated by alloc.
func alloc(n int, size, align uintptr) unsafe.Pointer {
	size = mem.AlignUp(size, align) * uintptr(n)
	var p unsafe.Pointer
	if unpriv() {
		p = allocBottom(n, size, align)
	} else {
		p = allocTop(n, size, align)
	}
	internal.Memset(p, 0, size)
	return p
}
