package noos

import (
	"builtin"
	"mem"
	"sync/atomic"
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
	builtin.Memset(p, 0, size)
	return p
}
