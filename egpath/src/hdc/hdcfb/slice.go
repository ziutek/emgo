package hdcfb

import (
	"bytes"
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
// All gorutines that use locked buffer must call Flush to unlock it.
//
// Any nonzero byte written to the slice is transfered to the display and
// replaces the corresponding byte in its DDRAM. Zero byte preserves previous
// content of corresponding byte in DDRAM.
//
// Write, WriteString and Flush methods are all lockless so even interrupt
// handler is permitted to use them.
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
			sl.buf.mod.Set()
			return
		}
		if atomic.AddInt32(&sl.buf.used, -1) == minInt32 {
			sl.fb.e.Send()
		}
	}
}

// Pos returns current write position in slice.
func (sl *Slice) Pos() int {
	return int(sl.p - sl.m)
}

// SetPos sets write position in slice.
func (sl *Slice) SetPos(p int) {
	if p >= int(sl.n-sl.m) {
		sl.p = sl.n
	}
	sl.p = sl.m + byte(p)
}

func (sl *Slice) Remain() int {
	return int(sl.n - sl.p)
}

// Flush marks slice as ready to draw and sets write position to p. There is no
// guarantee that Slice will be drawed to the display before subsequent write
// TODO: Implement Wait method that waits for buffers swap to give such
// guarantee.
func (sl *Slice) Flush(p int) {
	if atomic.AddInt32(&sl.buf.used, -1) == minInt32 {
		sl.fb.e.Send()
	}
	sl.buf = nil
	sl.SetPos(p)
}

// WriteString writes s to the slice. It returns number of bytes written. Only
// possible error is ErrTooLarge which means truncated write.
func (sl *Slice) WriteString(s string) (int, error) {
	sl.start()
	n := copy(sl.buf.data[sl.p:sl.n], s)
	sl.p += byte(n)
	if n != len(s) {
		return n, ErrTooLarge
	}
	return n, nil
}

func (sl *Slice) WriteByte(b byte) error {
	sl.start()
	if sl.p >= sl.n {
		return ErrTooLarge
	}
	sl.buf.data[sl.p] = b
	sl.p++
	return nil
}

func (sl *Slice) Fill(n int, b byte) {
	if n <= 0 {
		return
	}
	sl.start()
	n += int(sl.p)
	if n > int(sl.n) {
		n = int(sl.n)
	}
	bytes.Fill(sl.buf.data[sl.p:n], b)
	sl.p = byte(n)
}

// SyncSlice is simplified version of Slice. It is intended to be used by
// gorutine that performs drawing framebuffer to the LCD display. Its Write and
// WriteString methods can be only called between FB.Swap and FB.Draw calls.
// Typically it is used to draw something periodically (eg: current time).
type SyncSlice struct {
	fb   *FB
	m, n byte
	p    byte
}

// SyncSlice returns slice of framebuffer started at byte m and ended just
// before byte n.
func (fb *FB) SyncSlice(m, n int) SyncSlice {
	return SyncSlice{fb: fb, m: byte(m), n: byte(n), p: byte(m)}
}

// NewSyncSlice is like SyncSlice but returns pointer to the heap allocated
// SyncSlice.
func (fb *FB) NewSyncSlice(m, n int) *SyncSlice {
	s := new(SyncSlice)
	*s = fb.SyncSlice(m, n)
	return s
}

// WriteString: see Slice.WriteString.
func (sl *SyncSlice) WriteString(s string) (int, error) {
	n := copy(sl.fb.buf1.data[sl.p:sl.n], s)
	sl.p += byte(n)
	if n != len(s) {
		return n, ErrTooLarge
	}
	return n, nil
}

func (sl *SyncSlice) WriteByte(b byte) error {
	if sl.p >= sl.n {
		return ErrTooLarge
	}
	sl.fb.buf1.data[sl.p] = b
	sl.p++
	return nil
}

// Pos: see Slice.Pos.
func (sl *SyncSlice) Pos() int {
	return int(sl.p - sl.m)
}

// SetPos: see Slice.SetPos.
func (sl *SyncSlice) SetPos(p int) {
	if p >= int(sl.n-sl.m) {
		sl.p = sl.n
	}
	sl.p = sl.m + byte(p)
}

func (sl *SyncSlice) Remain() int {
	return int(sl.n - sl.p)
}

func (sl *SyncSlice) Fill(n int, b byte) {
	if n <= 0 {
		return
	}
	n += int(sl.p)
	if n > int(sl.n) {
		n = int(sl.n)
	}
	bytes.Fill(sl.fb.buf1.data[sl.p:n], b)
	sl.p = byte(n)
}
