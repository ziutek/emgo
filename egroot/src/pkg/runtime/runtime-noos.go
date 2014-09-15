// +build noos

package runtime

import (
	"builtin"
	"runtime/noos"
	"sync"
	"unsafe"
)

// This file imports noos package into runtime.
// Additonaly it contains things that can't be placed in noos package because
// of dependency loops.

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
	h := uintptr(unsafe.Pointer(&noos.Heap[0]))
	p := alignUp(h, align)
	m := size + (p - h)
	if m > uintptr(len(noos.Heap)) {
		panic("out of memory")
	}
	noos.Heap = noos.Heap[m:]
	hm.Unlock()

	builtin.Memset(unsafe.Pointer(p), 0, size)
	return unsafe.Pointer(p)
}

func init() {
	builtin.Alloc = alloc
}
