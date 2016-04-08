package linux

import (
	"internal"
	"mem"
	"sync/atomic"
	"syscall"
	"unsafe"
)

func end() uintptr

var aHeapEnd = end()

// alloc is very naive memory allocator that tries to be thread-safe without
// using mutexes.
func alloc(n int, size, align uintptr) unsafe.Pointer {
	size = mem.AlignUp(size, align) * uintptr(n)
	for {
		he := atomic.LoadUintptr(&aHeapEnd)
		p := mem.AlignUp(he, align)
		newhe := p + size
		if syscall.Brk(unsafe.Pointer(newhe)) < newhe {
			panic("out of memory")
		}
		if atomic.CompareAndSwapUintptr(&aHeapEnd, he, newhe) {
			internal.Memset(unsafe.Pointer(p), 0, size)
			return unsafe.Pointer(p)
		}
	}
}
