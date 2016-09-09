package hdcfb

import (
	"rtos"
	"sync/atomic"
	"syscall"

	"hdc"
)

type buffer struct {
	data []byte
	used int32
	mod  rtos.EventFlag
}

// FB represents double buffered text framebuffer. It contains two internal
// buffers. At any time one bufer is in active state and one in sync state.
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

// WaitAndSwap waits until any write to the active buffer or deadline occurs. It
// returns false in case of timeout, otherwise it swaps internal buffers and
// returns true. Deadline == 0 means no deadline.
func (fb *FB) WaitAndSwap(deadline int64) bool {
	ok := fb.buf0.mod.Wait(deadline)
	if ok {
		fb.buf0, fb.buf1 = fb.buf1, fb.buf0
	}
	atomic.AddInt32(&fb.buf1.used, minInt32)
	return ok
}

// Swap checks if there was any write to the active buffer. If it occured Swap
// swaps internal buffers and returns true. Otherwise it returns false.
func (fb *FB) Swap() bool {
	ok := fb.buf0.mod.Val() != 0
	if ok {
		fb.buf0, fb.buf1 = fb.buf1, fb.buf0
	}
	atomic.AddInt32(&fb.buf1.used, minInt32)
	return ok
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
	buf.mod.Clear()
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
