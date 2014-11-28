package noos

import (
	"builtin"
	"sync"
	"unsafe"
)

func alignUp(p, a uintptr) uintptr {
	a--
	return (p + a) &^ a
}

var hm sync.Mutex

// alloc is trivial memory allocator.
// For now there is no way to deallocate memory allocated by alloc.
func alloc(n int, size, align uintptr) unsafe.Pointer {
	size = alignUp(size, align) * uintptr(n)

	hm.Lock()
	h := uintptr(unsafe.Pointer(&Heap[0]))
	p := alignUp(h, align)
	m := size + (p - h)
	if m > uintptr(len(Heap)) {
		panic("out of memory")
	}
	Heap = Heap[m:]
	hm.Unlock()

	builtin.Memset(unsafe.Pointer(p), 0, size)
	return unsafe.Pointer(p)
}

func init() {
	builtin.Alloc = alloc
}
