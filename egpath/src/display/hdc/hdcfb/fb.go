package hdcfb

import (
	"rtos"
	"sync/atomic"
	"syscall"

	"display/hdc"
)

type buffer struct {
	data []byte
	used int
	mod  rtos.EventFlag
}

// FB represents double buffered text framebuffer. It contains two internal
// buffers. At any time one bufer is in active state and one in sync state.
type FB struct {
	d      *hdc.Display
	bufrel syscall.Event
	buf0   *buffer
	buf1   *buffer
}

// NewFB creates new framebuffer. If double is false only SyncSlices can be used
// to draw to it.
func NewFB(d *hdc.Display, double bool) *FB {
	fb := new(FB)
	fb.d = d
	if double {
		fb.bufrel = syscall.AssignEvent()
		fb.buf0 = new(buffer)
		fb.buf0.data = make([]byte, d.Cols*d.Rows)
	}
	fb.buf1 = new(buffer)
	fb.buf1.data = make([]byte, d.Cols*d.Rows)
	return fb
}

func (fb *FB) Display() *hdc.Display {
	return fb.d
}

const minInt = ^int(^uint(0) >> 1)

func panicSingleBuffer() {
	panic("hdcfb: single buffer only")
}

// Swap check whather the active internal buffer was modified. If yes then it
// swaps internal buffers and returns true. Otherwise it returns false.
func (fb *FB) Swap() bool {
	if fb.buf0 == nil {
		panicSingleBuffer()
	}
	if fb.buf1.used < 0 {
		panic("hcdfb: Draw not called")
	}
	if fb.buf0.mod.Value() == 0 {
		return false
	}
	fb.buf0, fb.buf1 = fb.buf1, fb.buf0
	atomic.AddInt(&fb.buf1.used, minInt) // Prevent use of buf1 by new writes.
	return true
}

// WaitAndSwap waits until any write to the active buffer or deadline occurs.
// After that it calls Swap and forwards its return value.
func (fb *FB) WaitAndSwap(deadline int64) bool {
	if fb.buf0 == nil {
		panicSingleBuffer()
	}
	fb.buf0.mod.Wait(1, deadline)
	return fb.Swap()
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
// It blocks until all writes to swapped buffer were finished.
func (fb *FB) Draw() error {
	buf1 := fb.buf1
	if fb.buf0 != nil && buf1.mod.Value() != 0 {
		for atomic.LoadInt(&buf1.used) != minInt {
			fb.bufrel.Wait()
		}
		buf1.used = 0
		buf1.mod.Reset(0)
	}
	data := buf1.data
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
