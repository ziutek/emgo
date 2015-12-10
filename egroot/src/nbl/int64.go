package nbl

import (
	"sync/atomic"
)

// Int64 implements int64 value for one writer and multiple readers. It is
// designed for architectures that does not support 64-bit atomic operations.
type Int64 struct {
	val [2]int64
	aba uintptr
}

func (u *Int64) StartLoad() uintptr {
	return atomic.LoadUintptr(&u.aba)
}

func (i *Int64) TryLoad(aba uintptr) int64 {
	return i.val[aba&1]
}

func (i *Int64) CheckLoad(aba uintptr) (uintptr, bool) {
	aba1 := atomic.LoadUintptr(&i.aba)
	return aba1, aba1 == aba
}

// BUG: On 32bit CPUs Load doesnt guarantee that it returns valid value. The
// probability of failure depends on the frequency of updates:
// 1 kHz: aba wraps onece per 1193 houres,
// 1 MHz: aba wraps once per 72 minutes.
// Load fails if aba wraps beetwen StartLoad and CheckLoad or between subsequent
// CheckLoads and is read twice with the same value.
func (i *Int64) Load() int64 {
	aba := i.StartLoad()
	for {
		v := i.TryLoad(aba)
		var ok bool
		if aba, ok = i.CheckLoad(aba); ok {
			return v
		}
	}
}

// WriterLoad is more efficient than Load but can be used only by writer.
func (i *Int64) WriterLoad() int64 {
	return i.val[i.aba&1]
}

// WriterStore can be used only when there is only ONE writer.
func (i *Int64) WriterStore(v int64) {
	aba := i.aba
	aba++
	i.val[aba&1] = v
	atomic.StoreUintptr(&i.aba, aba)
}

func (i *Int64) WriterAdd(v int64) int64 {
	v += i.WriterLoad()
	i.WriterStore(v)
	return v
}

func (i *Int64) WriterSub(v int64) int64 {
	v = i.WriterLoad() - v
	i.WriterStore(v)
	return v
}
