package nbl

import (
	"sync/atomic"
	"sync/fence"
)

// Uint64 implements uint64 value for one writer and multiple readers. It is
// designed for architectures that does not support 64-bit atomic operations.
type Uint64 struct {
	val [2]uint64
	aba uintptr
}

func (u *Uint64) StartLoad() uintptr {
	return atomic.LoadUintptr(&u.aba)
}

func (u *Uint64) TryLoad(aba uintptr) uint64 {
	return u.val[aba&1]
}

func (u *Uint64) CheckLoad(aba uintptr) (uintptr, bool) {
	fence.Compiler()
	aba1 := atomic.LoadUintptr(&u.aba)
	return aba1, aba1 == aba
}

func (u *Uint64) Load() uint64 {
	aba := u.StartLoad()
	for {
		v := u.TryLoad(aba)
		var ok bool
		if aba, ok = u.CheckLoad(aba); ok {
			return v
		}
	}
}

// WriterLoad is more efficient than Load but can be used only by writer.
func (u *Uint64) WriterLoad() uint64 {
	return u.val[u.aba&1]
}

// WriterStore can be used only when there is only ONE writer.
func (u *Uint64) WriterStore(v uint64) {
	aba := u.aba
	aba++
	u.val[aba&1] = v
	fence.Memory()
	atomic.StoreUintptr(&u.aba, aba)
}

func (u *Uint64) WriterAdd(v uint64) uint64 {
	v += u.WriterLoad()
	u.WriterStore(v)
	return v
}

func (u *Uint64) WriterSub(v uint64) uint64 {
	v = u.WriterLoad() - v
	u.WriterStore(v)
	return v
}
