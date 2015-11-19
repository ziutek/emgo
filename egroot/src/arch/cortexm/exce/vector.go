package exce

import (
	"bits"
	"builtin"
	"mmio"
	"sync/fence"
	"unsafe"
)

// Vector represents element of interrupt table.
type Vector struct {
	handler func()
}

// VectorFor returns Vector that correspods to the handler.
func VectorFor(handler func()) Vector {
	return Vector{handler}
}

var (
	vto      = mmio.PtrU32(unsafe.Pointer(uintptr(0xe000ed08)))
	activeVT []Vector
)

// UseTable instructs CPU to use vt as vector table. vt should be properly
// aligned. Minimum alignment is 32 words which is enough for up to 16 external
// interrupts. For more interrupts, adjust the alignment by rounding up to the
// next power of two. UseTable doesn't work for Cortex-M0.
//
// This function is designed to be used by runtime. You generaly shouldn't use
// it if MaxTask > 0 in your linker script.
func UseTable(vt []Vector) {
	activeVT = vt
	fence.Memory()
	vto.Store(uint32(uintptr(unsafe.Pointer(&vt[0]))))
}

// UseHandler changes handler in currently used vector table.
func (e Exce) UseHandler(handler func()) {
	if int(e) >= len(activeVT) {
		panic("exce: vector table is too short")
	}
	activeVT[e] = VectorFor(handler)
	fence.Sync()
}

// NewTable allocates new (properly aligned) vector table for n
// interrupt vectors.
func NewTable(n int) []Vector {
	if n < 0 || n > 256 {
		panic("bad vector table length")
	}
	exp := 32 - bits.LeadingZeros32(uint32(n-1))
	if exp < 5 {
		exp = 5
	}
	m := 1 << exp
	tsize := uintptr(m) * unsafe.Sizeof(Vector{})
	vt := (*[256]Vector)(builtin.Alloc(1, tsize, tsize))
	return vt[:n]
}
