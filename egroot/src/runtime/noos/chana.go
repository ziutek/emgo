package noos

import (
	"bits"
	"builtin"
	"mem"
	"sync/atomic"
	"sync/fence"
	"syscall"
	"unsafe"
)

// Asynchronous channels - lockless implementation.

func bit(w *uint32, n uintptr) bool {
	fence.Compiler()
	return *w&(1<<n) != 0
}

func atomicToggleBit(addr *uint32, n uintptr) {
	mask := uint32(1) << n
	for {
		old := atomic.LoadUint32(addr)
		if atomic.CompareAndSwapUint32(addr, old, old^mask) {
			return
		}
	}
}

type chanA struct {
	event  syscall.Event // Event must be the first field - see chanSelect.
	tosend uintptr
	torecv uintptr
	cap    uintptr
	mask   uintptr
	step   uintptr
	buf    unsafe.Pointer
	rd     *[1 << 28]uint32
	closed int32
}

func makeChanA(cap int, size, align uintptr) *chanA {
	c := new(chanA)
	c.event = syscall.AssignEvent()
	c.cap = uintptr(cap)
	c.mask = ^uintptr(0) >> (bits.LeadingZerosPtr(c.cap-1) - 1)
	c.step = mem.AlignUp(size, align)
	c.buf = builtin.Alloc(cap, size, align)
	rdlen := cap / 32
	if rdlen*32 < cap {
		rdlen++
	}
	c.rd = (*[1 << 28]uint32)(builtin.Alloc(rdlen, 4, 4))
	return c
}

func (c *chanA) Close() {
	if c == nil {
		panicCloseNil()
	}
	atomic.StoreInt32(&c.closed, 1)
	c.event.Send()
}

func (c *chanA) panicIfClosed() {
	if atomic.LoadInt32(&c.closed) != 0 {
		panicClosed()
	}
}

// BUG: TrySend, TryRecv and Len can be affected by ABA problem.
// They all rely on comparasion of two values of some counter taken at two
// points in time. If both values are equal they assume that this counter
// hasn't been modified. TrySend increments tosend couner, TryRecv increments
// torecv couner. If sizeof(uintptr) == 4 both couners can wrap after about
// 1<<31 increments in worst case. 1<<31 is quiet big value but on busy system,
// some low priority task can be blocked for enough time (say, an hour in case
// of 100 MHz CPU) to have a chanse to be affected.

// TrySend tries to reserve place in c's internal ring buffer.
// After return, if p != nil then p contains pointer to the place in internal
// buffer where data can be stored. After data store sender need to call
// c.Done(d). If p == nil then d == cagain which means that the internal buffer
// was full and TrySend can be called again.
func (c *chanA) TrySend(_, _ unsafe.Pointer) (p unsafe.Pointer, d uintptr) {
	var n uintptr
	nmask := c.mask >> 1
	abalsb := c.mask &^ nmask
	for {
		c.panicIfClosed()
		tosend := atomic.LoadUintptr(&c.tosend)
		if atomic.LoadUintptr(&c.torecv)&c.mask == (tosend^abalsb)&c.mask {
			// Channel is full.
			return nil, cagain
		}
		n = tosend & nmask
		if bit(&c.rd[n>>5], n&31) {
			// This element is still being received.
			return nil, cagain
		}
		var next uintptr
		if n+1 == c.cap {
			next = (tosend + abalsb) &^ nmask
		} else {
			next = tosend + 1
		}
		if atomic.CompareAndSwapUintptr(&c.tosend, tosend, next) {
			break
		}
	}
	return unsafe.Pointer(uintptr(c.buf) + c.step*n), n
}

// TryRecv tries to obtain pointer to the data ready to read from channel.
// After return, if p != nil then p contains pointer to the data. Receiver need
// to copy the data and after that it need to call c.Done(d). If p == nil then
// d can be equal to cagain or cclosed.
func (c *chanA) TryRecv(_, _ unsafe.Pointer) (p unsafe.Pointer, d uintptr) {
	var n uintptr
	nmask := c.mask >> 1
	abalsb := c.mask &^ nmask
	for {
		torecv := atomic.LoadUintptr(&c.torecv)
		if atomic.LoadUintptr(&c.tosend) == torecv {
			// Channel is empty.
			if atomic.LoadInt32(&c.closed) != 0 {
				return nil, cclosed
			}
			return nil, cagain
		}
		n = torecv & nmask
		if !bit(&c.rd[n>>5], n&31) {
			// This element is still being sent.
			return nil, cagain
		}
		var next uintptr
		if n+1 == c.cap {
			next = (torecv + abalsb) &^ nmask
		} else {
			next = torecv + 1
		}
		if atomic.CompareAndSwapUintptr(&c.torecv, torecv, next) {
			break
		}
	}
	return unsafe.Pointer(uintptr(c.buf) + c.step*n), n
}

// Send reserves place in c's internal buffer and returns pointer to it in p.
// Sender need to call c.Done(d) after data transfer.
func (c *chanA) Send(_ unsafe.Pointer) (p unsafe.Pointer, d uintptr) {
	for {
		p, d = c.TrySend(nil, nil)
		if p != nil {
			return
		}
		c.event.Wait()
	}
}

// Recv returns pointer to the data ready to read from channel in p or p == nil
// if channel is closed. If p == nil then d == cclosed. Receiver need to copy
// data from channel and call c.Done(d) after that
func (c *chanA) Recv(_ unsafe.Pointer) (p unsafe.Pointer, d uintptr) {
	for {
		p, d = c.TryRecv(nil, nil)
		if p != nil || d == cclosed {
			return
		}
		c.event.Wait()
	}
}

// Done must be called after data transfer that sender or reveiver need to
// perform..
func (c *chanA) Done(n uintptr) {
	fence.Memory()
	atomicToggleBit(&c.rd[n>>5], n&31)
	c.event.Send()
}

func (c *chanA) Len() int {
	torecv := atomic.LoadUintptr(&c.torecv)
	fence.Compiler()
	tosend := atomic.LoadUintptr(&c.tosend)
	fence.Compiler()
	for {
		tr := atomic.LoadUintptr(&c.torecv)
		fence.Compiler()
		ts := atomic.LoadUintptr(&c.tosend)
		fence.Compiler()
		if tr == torecv {
			break
		}
		torecv = tr
		if ts == tosend {
			break
		}
		tosend = ts
	}
	nmask := c.mask >> 1
	abalsb := c.mask &^ nmask
	nr := torecv & nmask
	ns := tosend & nmask
	if (torecv^tosend)&abalsb == 0 {
		return int(ns - nr)
	}
	return int(c.cap - nr + ns)
}

func (c *chanA) Cap() int {
	return int(c.cap)
}

type chanAMethods struct {
	Send       func(c *chanA, e unsafe.Pointer) (p unsafe.Pointer, d uintptr)
	Recv       func(c *chanA, e unsafe.Pointer) (p unsafe.Pointer, d uintptr)
	TrySend    func(c *chanA, e, _ unsafe.Pointer) (p unsafe.Pointer, d uintptr)
	TryRecv    func(c *chanA, e, _ unsafe.Pointer) (p unsafe.Pointer, d uintptr)
	CancelSend func()
	CancelRecv func()
	Done       func(c *chanA, d uintptr)
	Close      func(c *chanA)
	Len        func(c *chanA) int
	Cap        func(c *chanA) int
}

var cam = chanAMethods{
	(*chanA).Send,
	(*chanA).Recv,
	(*chanA).TrySend,
	(*chanA).TryRecv,
	nil,
	nil,
	(*chanA).Done,
	(*chanA).Close,
	(*chanA).Len,
	(*chanA).Cap,
}
