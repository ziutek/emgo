package linux

import (
	"mem"
	"sync/atomic"
	"sync/fence"
	"syscall"
	"unsafe"
)

var (
	aHeapBegin uintptr
	aHeapEnd   uintptr
)

func init() {
	const length = 100 * 1024 * 1024
	// Preallocate length bytes of (zeroed) virtual memory.
	aHeapBegin, _ = syscall.Mmap(
		^uintptr(0), length,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_PRIVATE|syscall.MAP_ANONYMOUS,
		0, 0,
	)
	aHeapEnd = aHeapBegin + length
}

// alloc is primitive non-blocking memory allocator.
func alloc(n int, size, align uintptr) unsafe.Pointer {
	size = mem.AlignUp(size, align) * uintptr(n)
	for {
		hb := atomic.LoadUintptr(&aHeapBegin)
		p := mem.AlignUp(hb, align)
		newhb := p + size
		if atomic.CompareAndSwapUintptr(&aHeapBegin, hb, newhb) {
			if newhb > aHeapEnd {
				panic("out of memory")
			}
			// Returned memory is zeroed by OS.
			return unsafe.Pointer(p)
		}
	}
}
