package noos

import (
	"builtin"
	"sync/atomic"
	"sync/barrier"
	"unsafe"
)

// Asynchronous channels

func bit(w *uint32, n uintptr) bool {
	barrier.Compiler()
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

func alignUp(p, a uintptr) uintptr {
	a--
	return (p + a) &^ a
}

type chanA struct {
	event  Event // Event must be the first field - see chanSelect.
	tosend uint32
	torecv uint32
	cap    uintptr
	step   uintptr
	rd     *[1 << 30]uint32
	buf    unsafe.Pointer
	closed int32
}

func makeChanA(cap int, size, align uintptr) *chanA {
	c := new(chanA)
	c.event = AssignEvent()
	c.cap = uintptr(cap)
	c.step = alignUp(size, align)
	c.rd = (*[1 << 30]uint32)(builtin.Alloc(1+cap/32, 4, 4))
	c.buf = builtin.Alloc(cap, size, align)
	return c
}

const b31 = 1 << 31

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

// TrySend tries to reserve place in c's internl ring buffer.
// After return, if p != nil then p contains pointer to the place in internal
// buffer where data can be stored. After that sender need to call c.Done(d).
// If p == nil then d == cagain which means that the internal buffer was full
// and TrySend can be called again.
func (c *chanA) TrySend(_, _ unsafe.Pointer) (p unsafe.Pointer, d uintptr) {
	if c == nil {
		return nil, cagain
	}
	var n uintptr
	for {
		c.panicIfClosed()
		tosend := atomic.LoadUint32(&c.tosend)
		if atomic.LoadUint32(&c.torecv) == tosend^b31 {
			// Channel is full.
			return nil, cagain
		}
		n = uintptr(tosend &^ b31)
		if bit(&c.rd[n>>5], n&31) {
			// This element is still being received.
			return nil, cagain
		}
		var next uint32
		if n+1 == c.cap {
			next = ^tosend & b31
		} else {
			next = tosend + 1
		}
		if atomic.CompareAndSwapUint32(&c.tosend, tosend, next) {
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
	if c == nil {
		return nil, cagain
	}
	var n uintptr
	for {
		torecv := atomic.LoadUint32(&c.torecv)
		if atomic.LoadUint32(&c.tosend) == torecv {
			// Channel is empty.
			if atomic.LoadInt32(&c.closed) != 0 {
				return nil, cclosed
			}
			return nil, cagain
		}
		n = uintptr(torecv &^ b31)
		if !bit(&c.rd[n>>5], n&31) {
			// This element is still being sent.
			return nil, cagain
		}
		var next uint32
		if n+1 == c.cap {
			next = ^torecv & b31
		} else {
			next = torecv + 1
		}
		if atomic.CompareAndSwapUint32(&c.torecv, torecv, next) {
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
	barrier.Memory()
	atomicToggleBit(&c.rd[n>>5], n&31)
	c.event.Send()
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
}
