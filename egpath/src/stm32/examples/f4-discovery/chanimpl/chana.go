package main

import (
	"builtin"
	"runtime/noos"
	"sync/atomic"
	"sync/barrier"
	"unsafe"
)

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

type ChanA struct {
	event  noos.Event
	tosend uint32
	torecv uint32
	cap    uintptr
	step   uintptr
	rd     *[1 << 30]uint32
	buf    unsafe.Pointer
	closed int32
}

func NewChanA(cap int, size, align uintptr) *ChanA {
	c := new(ChanA)
	c.event = noos.AssignEvent()
	c.cap = uintptr(cap)
	c.step = alignUp(size, align)
	c.rd = (*[1 << 30]uint32)(builtin.Alloc(1+cap/32, 4, 4))
	c.buf = builtin.Alloc(cap, size, align)
	return c
}

const b31 = 1 << 31

func (c *ChanA) Close() {
	atomic.StoreInt32(&c.closed, 1)
	c.event.Send()
}

func (c *ChanA) panicIfClosed() {
	if atomic.LoadInt32(&c.closed) != 0 {
		panic("send: closed channel")
	}
}

// TrySend tries to send the value of variable pointed by e.
// If p != nil the data transfer need to be performed by sender and after that
// sender need to call c.Done(d). Otherwise d == 0 and any data wasn't send
// because the internal buffer was full.
func (c *ChanA) TrySend(e unsafe.Pointer) (p unsafe.Pointer, d uintptr) {
	var n uintptr
	for {
		c.panicIfClosed()
		tosend := atomic.LoadUint32(&c.tosend)
		if atomic.LoadUint32(&c.torecv) == tosend^b31 {
			// Channel is full.
			return nil, 0
		}
		n = uintptr(tosend &^ b31)
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
	for bit(&c.rd[n>>5], n&31) {
		// This element is still being received.
		c.event.Wait()
	}
	return unsafe.Pointer(uintptr(c.buf) + c.step*n), n
}

// TryRecv tries to receive a value from channel and store it into variable
// pointed by e.
// If p != nil the data transfer need to be performed by receiver and after
// that receiver need to call c.Done(d). Otherwise any data wasn't received
// beacuse the internal buffer was empty and d != 0 if channel is closed.
func (c *ChanA) TryRecv(e unsafe.Pointer) (p unsafe.Pointer, d uintptr) {
	var n uintptr
	for {
		torecv := atomic.LoadUint32(&c.torecv)
		if atomic.LoadUint32(&c.tosend) == torecv {
			// Channel is empty.
			return nil, uintptr(atomic.LoadInt32(&c.closed))
		}
		n = uintptr(torecv &^ b31)
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
	for !bit(&c.rd[n>>5], n&31) {
		// This element is still being sent.
		c.event.Wait()
	}
	return unsafe.Pointer(uintptr(c.buf) + c.step*n), n
}

// Send tries to send the value of variable pointed by e.
// If p != nil the data transfer need to be performed by sender and after that
// sender need to call c.Done(d). This function actually always returns p != nil
// but description is intended to be compatible with other channel types where
// Send can return p == nil.
func (c *ChanA) Send(e unsafe.Pointer) (p unsafe.Pointer, d uintptr) {
	for {
		p, d = c.TrySend(e)
		if p != nil {
			return
		}
		c.event.Wait()
	}
}

// Recv tries to receive a value from channel and store it into variable
// pointed by e.
// If p != nil the data transfer need to be performed by receiver and
// after that receiver need to call c.Done(d). Otherwise d != 0, channel is
// closed and there was no data to receive.
func (c *ChanA) Recv(e unsafe.Pointer) (p unsafe.Pointer, d uintptr) {
	for {
		p, d = c.TryRecv(e)
		if p != nil || d != 0 {
			return
		}
		c.event.Wait()
	}
}

// Done must be called when the data transfer that sender or reveiver need to
// perform was completed.
func (c *ChanA) Done(n uintptr) {
	barrier.Memory()
	atomicToggleBit(&c.rd[n>>5], n&31)
	c.event.Send()
}
