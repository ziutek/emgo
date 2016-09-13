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

func (fb *FB) Display() *hdc.Display {
	return fb.d
}

const minInt32 = -2147483648

// Swap swaps internal buffers. It reports whather active buffer was modified
// and need to be drawed.
func (fb *FB) Swap() bool {
	if fb.buf1.used < 0 {
		panic("hcdfb: Draw not called")
	}
	fb.buf0, fb.buf1 = fb.buf1, fb.buf0
	buf1 := fb.buf1
	atomic.AddInt32(&buf1.used, minInt32)
	mod := buf1.mod.Val() != 0
	if !mod {
		buf1.used = 0
	}
	return mod
}

// WaitAndSwap waits until any write to the active buffer or deadline occurs.
// After that it calls Swap and forwards its return value.
func (fb *FB) WaitAndSwap(deadline int64) bool {
	fb.buf0.mod.Wait(deadline)
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
	if buf1.mod.Val() != 0 {
		for atomic.LoadInt32(&buf1.used) != minInt32 {
			fb.e.Wait()
		}
		buf1.used = 0
		buf1.mod.Clear()
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
