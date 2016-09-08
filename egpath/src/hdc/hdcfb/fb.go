package hdcfb

import (
	"errors"
	"sync/atomic"
	"syscall"

	"hdc"
)

type buffer struct {
	data []byte
	used int32
}

// FB represents double buffered text framebuffer.
type FB struct {
	d    *hdc.Display
	e    syscall.Event
	buf0 *buffer
	buf1 *buffer
}

func NewFB(d *hdc.Display) *FB {
	fb := new(FB)
	fb.d = d
	fb.e = syscall.AssignEvent()
	fb.buf0 = new(buffer)
	fb.buf1 = new(buffer)
	fb.buf0.data = make([]byte, d.Cols*d.Rows)
	fb.buf1.data = make([]byte, d.Cols*d.Rows)
	return fb
}

const minInt32 = -2147483648

// Swap swaps internal buffers.
func (fb *FB) Swap() {
	buf := fb.buf0
	fb.buf0 = fb.buf1
	atomic.AddInt32(&buf.used, minInt32)
	fb.buf1 = buf
}

func (fb *FB) draw(data []byte) (err error) {
	for i, b := range data {
		if b != 0 {
			err = fb.d.WriteByte(b)
			data[i] = 0
		} else {
			err = fb.d.Shift(hdc.Cursor | hdc.Right)
		}
		if err != nil {
			break
		}
	}
	return
}

// Draw draws content of the previously swapped internal buffer to the display.
// It can block if there are still writers of swapped buffer that did not called
// Slice.Flush method.
func (fb *FB) Draw() error {
	buf := fb.buf1
	for atomic.LoadInt32(&buf.used) != minInt32 {
		fb.e.Wait()
	}
	buf.used = 0
	data := buf.data
	if fb.d.Rows != 4 {
		return fb.draw(data)
	}
	c := fb.d.Cols
	if err := fb.draw(data[:c]); err != nil {
		return err
	}
	if err := fb.draw(data[2*c : 3*c]); err != nil {
		return err
	}
	if err := fb.draw(data[c : 2*c]); err != nil {
		return err
	}
	return fb.draw(data[3*c:])
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

func (sl *Slice) start() {
	if sl.buf != nil {
		return
	}
	for {
		sl.buf = sl.fb.buf0
		if atomic.AddInt32(&sl.buf.used, 1) > 0 {
			return
		}
		if atomic.AddInt32(&sl.buf.used, -1) == minInt32 {
			sl.fb.e.Send()
		}
	}
}

func (sl *Slice) Flush() {
	if atomic.AddInt32(&sl.buf.used, -1) == minInt32 {
		sl.fb.e.Send()
	}
	sl.buf = nil
	sl.p = sl.m
}

var ErrTooLarge = errors.New("hdcfb: too large")

func (sl *Slice) WriteString(s string) (int, error) {
	sl.start()
	n := copy(sl.buf.data[sl.p:sl.n], s)
	sl.p += byte(n)
	if n != len(s) {
		return n, ErrTooLarge
	}
	return n, nil
}

func (sl *Slice) Write(s []byte) (int, error) {
	sl.start()
	n := copy(sl.buf.data[sl.p:sl.n], s)
	sl.p += byte(n)
	if n != len(s) {
		return n, ErrTooLarge
	}
	return n, nil
}
