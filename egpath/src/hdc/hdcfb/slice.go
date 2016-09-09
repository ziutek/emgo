package hdcfb

import (
	"errors"
	"sync/atomic"
)

var ErrTooLarge = errors.New("hdcfb: too large")

// Slice represents slice of the framebuffer. Multiple gorutines can use one
// framebuffer concurently using nonoverlaping slices of it. One slice cannot
// be used concurently by multiple gorutines.
//
// Starting writing to the slice marks current internal buffer as used and
// blocks transfer of its content to the display (that is FB.Draw() will block).
// Flush methods of all slices that use locked buffer must be called to unlock
// it.
//
// Any nonzero byte written to the slice is transfered to the display and
// replaces the corresponding byte in its DDRAM. Zero byte preserves previous
// content of corresponding byte in DDRAM.
//
// Write, WriteString and Flush methods are all lockless.
type Slice struct {
	fb   *FB
	buf  *buffer
	m, n byte
	p    byte
}

// Slice returns slice of framebuffer started at byte m and ended just before
// byte n.
func (fb *FB) Slice(m, n int) Slice {
	return Slice{fb: fb, m: byte(m), n: byte(n), p: byte(m)}
}

// NewSlice is like Slice but returns pointer to the heap allocated Slice.
func (fb *FB) NewSlice(m, n int) *Slice {
	s := new(Slice)
	*s = fb.Slice(m, n)
	return s
}

func (sl *Slice) start() {
	if sl.buf != nil {
		return
	}
	for {
		sl.buf = sl.fb.buf0
		if atomic.AddInt32(&sl.buf.used, 1) > 0 {
			sl.fb.buf0.mod.Set()
			return
		}
		if atomic.AddInt32(&sl.buf.used, -1) == minInt32 {
			sl.fb.e.Send()
		}
	}
}

func (sl *Slice) Pos() int {
	return int(sl.p - sl.m)
}

func (sl *Slice) SetPos(p int) {
	if p >= int(sl.n-sl.m) {
		sl.p = sl.n
	}
	sl.p = sl.m + byte(p)
}

func (sl *Slice) Flush(p int) {
	if atomic.AddInt32(&sl.buf.used, -1) == minInt32 {
		sl.fb.e.Send()
	}
	sl.buf = nil
	sl.SetPos(p)
}

func (sl *Slice) WriteString(s string) (int, error) {
	sl.start()
	n := copy(sl.buf.data[sl.p:sl.n], s)
	sl.p += byte(n)
	if n != len(s) {
		return n, ErrTooLarge
	}
	return n, nil
}

// SyncSlice is simplified version of Slice. It is intended to be used by
// gorutine that performs drawing framebuffer to the LCD display. Its Write and
// WriteString methods can be only called between FB.Swap and FB.Draw calls.
// Typically it is used to draw something periodically, eg:drawing current time.
type SyncSlice struct {
	fb   *FB
	m, n byte
	p    byte
}

func (fb *FB) SyncSlice(m, n int) SyncSlice {
	return SyncSlice{fb: fb, m: byte(m), n: byte(n), p: byte(m)}
}

func (fb *FB) NewSyncSlice(m, n int) *SyncSlice {
	s := new(SyncSlice)
	*s = fb.SyncSlice(m, n)
	return s
}

func (sl *SyncSlice) WriteString(s string) (int, error) {
	n := copy(sl.fb.buf1.data[sl.p:sl.n], s)
	sl.p += byte(n)
	if n != len(s) {
		return n, ErrTooLarge
	}
	return n, nil
}

func (sl *SyncSlice) Pos() int {
	return int(sl.p - sl.m)
}

func (sl *SyncSlice) SetPos(p int) {
	if p >= int(sl.n-sl.m) {
		sl.p = sl.n
	}
	sl.p = sl.m + byte(p)
}
