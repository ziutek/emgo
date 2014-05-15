package noos

import (
	"builtin"
	"sync/atomic"
	"sync/barrier"
	"unsafe"
)

func makeChan(cap int, size, align uintptr) (c builtin.Chan) {
	if cap == 0 {
		c.P = unsafe.Pointer(makeChanS())
		c.M = (*builtin.ChanMethods)(unsafe.Pointer(&csm))
	} else {
		c.P = unsafe.Pointer(makeChanA(cap, size, align))
		c.M = (*builtin.ChanMethods)(unsafe.Pointer(&cam))
	}
	return
}

func init() {
	builtin.MakeChan = makeChan
}

// Synchronous channels

const (
	ready int32 = iota
	recvWait
	sendWait
	closed
)

type chanS struct {
	event Event
	src   unsafe.Pointer
	dst   unsafe.Pointer
	state int32
}

func makeChanS() *chanS {
	c := new(chanS)
	c.event = AssignEvent()
	return c
}

func (c *chanS) Close() {
	atomic.StoreInt32(&c.state, closed)
	c.event.Send()
}

func panicIfClosed(state int32) {
	if state == closed {
		panic("send: closed channel")
	}
}

func (c *chanS) getDstAndUnlock() (unsafe.Pointer, uintptr) {
	// Save c.dst. It is only thing that is need for sender and receiver to
	// complete this communication.
	dst := c.dst
	// Use CAS to not reopen if channel was closed in the meantime.
	atomic.CompareAndSwapInt32(&c.state, recvWait, ready)
	barrier.Compiler()

	// Signal other senders and receivers that now c can be used by next pair.
	atomic.StorePointer(&c.dst, nil)
	atomic.StorePointer(&c.src, nil)
	c.event.Send()

	return *(*unsafe.Pointer)(dst), uintptr(dst)
}

func (c *chanS) getSrcAndUnlock() (unsafe.Pointer, uintptr) {
	// Save c.src. It is only thing that is need for sender and receiver to
	// complete this communication.
	src := c.src
	// Use CAS to not reopen if channel was closed in the meantime.
	atomic.CompareAndSwapInt32(&c.state, sendWait, ready)
	barrier.Compiler()

	// Signal other senders and receivers that now c can be used by next pair.
	atomic.StorePointer(&c.dst, nil)
	atomic.StorePointer(&c.src, nil)
	c.event.Send()

	return *(*unsafe.Pointer)(src), uintptr(src)
}

// Send tries to send the value of variable pointed by e.
// If p != nil the data transfer need to be performed by sender and after that
// sender need to call c.Done(d). Otherwise the data transfer was performed by
// receiver.
func (c *chanS) Send(e unsafe.Pointer) (p unsafe.Pointer, d uintptr) {
	panicIfClosed(atomic.LoadInt32(&c.state))
	for !atomic.CompareAndSwapPointer(&c.src, nil, unsafe.Pointer(&e)) {
		// Other sender locked this channel.
		c.event.Wait()
		panicIfClosed(atomic.LoadInt32(&c.state))
	}
	if atomic.CompareAndSwapInt32(&c.state, ready, sendWait) {
		// Passive sender.
		for atomic.LoadPointer(&e) != nil {
			c.event.Wait()
			panicIfClosed(atomic.LoadInt32(&c.state))
		}
		return nil, 0
	}
	// Active sender
	return c.getDstAndUnlock()
}

// Recv tries to receive a value from channel and store it into variable
// pointed by e.
// If p != nil the data transfer need to be performed by receiver and
// after that receiver need to call c.Done(d). Otherwise d == 0 means that the
// data transfer was performed by sender, d != 0 means that channel is closed
// and there was no data to receive.
func (c *chanS) Recv(e unsafe.Pointer) (p unsafe.Pointer, d uintptr) {
	if atomic.LoadInt32(&c.state) == closed {
		return nil, 1
	}
	for !atomic.CompareAndSwapPointer(&c.dst, nil, unsafe.Pointer(&e)) {
		// Other receiver locked this channel.
		c.event.Wait()
		if atomic.LoadInt32(&c.state) == closed {
			return nil, 1
		}
	}
	if atomic.CompareAndSwapInt32(&c.state, ready, recvWait) {
		// Passive receiver.
		for atomic.LoadPointer(&e) != nil {
			c.event.Wait()
			if atomic.LoadInt32(&c.state) == closed {
				return nil, 1
			}
		}
		return nil, 0
	}
	// Active receiver.
	return c.getSrcAndUnlock()
}

// TrySend tries to send the value of variable pointed by e.
// If p != nil the data transfer need to be performed by sender and after that
// sender need to call c.Done(d). Otherwise d == 0 and any data wasn't send
// because there wasn't any receiver waiting for communication.
func (c *chanS) TrySend(e unsafe.Pointer) (p unsafe.Pointer, d uintptr) {
	if !atomic.CompareAndSwapPointer(&c.src, nil, unsafe.Pointer(&e)) {
		panicIfClosed(atomic.LoadInt32(&c.state))
		// Other sender locked this channel.
		return nil, 0
	}
	if state := atomic.LoadInt32(&c.state); state != recvWait {
		panicIfClosed(state)
		// No any receiver waiting for communication.
		atomic.StorePointer(&c.src, nil)
		c.event.Send()
		return nil, 0
	}
	// Active sender
	return c.getDstAndUnlock()
}

// TryRecv tries to receive a value from channel and store it into variable
// pointed by e.
// If p != nil the data transfer need to be performed by receiver and after that
// receiver need to call c.Done(d). Otherwise, any data wasn't received and
// d == 0 means that there wasn't any sender waiting for communication, d != 0
// means that channel is closed.
func (c *chanS) TryRecv(e unsafe.Pointer) (p unsafe.Pointer, d uintptr) {
	if !atomic.CompareAndSwapPointer(&c.dst, nil, unsafe.Pointer(&e)) {
		// Other receiver locked this channel.
		if atomic.LoadInt32(&c.state) == closed {
			return nil, 1
		}
		return nil, 0
	}
	if state := atomic.LoadInt32(&c.state); state != sendWait {
		if state == closed {
			return nil, 1
		}
		// No any sender waiting for communication.
		atomic.StorePointer(&c.dst, nil)
		c.event.Send()
		return nil, 0
	}
	// Active receiver
	return c.getSrcAndUnlock()
}

// Done must be called when the data transfer that sender or reveiver need to
// perform was completed.
func (c *chanS) Done(d uintptr) {
	barrier.Memory()
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(d)), nil)
	c.event.Send()
}

type chanSMethods struct {
	Send    func(c *chanS, e unsafe.Pointer) (p unsafe.Pointer, d uintptr)
	Recv    func(c *chanS, e unsafe.Pointer) (p unsafe.Pointer, d uintptr)
	TrySend func(c *chanS, e unsafe.Pointer) (p unsafe.Pointer, d uintptr)
	TryRecv func(c *chanS, e unsafe.Pointer) (p unsafe.Pointer, d uintptr)
	Done    func(c *chanS, d uintptr)
	Close   func(c *chanS)
}

var csm = chanSMethods{
	(*chanS).Send,
	(*chanS).Recv,
	(*chanS).TrySend,
	(*chanS).TryRecv,
	(*chanS).Done,
	(*chanS).Close,
}

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
	event  Event
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
	atomic.StoreInt32(&c.closed, 1)
	c.event.Send()
}

func (c *chanA) panicIfClosed() {
	if atomic.LoadInt32(&c.closed) != 0 {
		panic("send: closed channel")
	}
}

// TrySend tries to send the value of variable pointed by e.
// If p != nil the data transfer need to be performed by sender and after that
// sender need to call c.Done(d). Otherwise d == 0 and any data wasn't send
// because the internal buffer was full.
func (c *chanA) TrySend(e unsafe.Pointer) (p unsafe.Pointer, d uintptr) {
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
func (c *chanA) TryRecv(e unsafe.Pointer) (p unsafe.Pointer, d uintptr) {
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
func (c *chanA) Send(e unsafe.Pointer) (p unsafe.Pointer, d uintptr) {
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
func (c *chanA) Recv(e unsafe.Pointer) (p unsafe.Pointer, d uintptr) {
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
func (c *chanA) Done(n uintptr) {
	barrier.Memory()
	atomicToggleBit(&c.rd[n>>5], n&31)
	c.event.Send()
}

type chanAMethods struct {
	Send    func(c *chanA, e unsafe.Pointer) (p unsafe.Pointer, d uintptr)
	Recv    func(c *chanA, e unsafe.Pointer) (p unsafe.Pointer, d uintptr)
	TrySend func(c *chanA, e unsafe.Pointer) (p unsafe.Pointer, d uintptr)
	TryRecv func(c *chanA, e unsafe.Pointer) (p unsafe.Pointer, d uintptr)
	Done    func(c *chanA, d uintptr)
	Close   func(c *chanA)
}

var cam = chanAMethods{
	(*chanA).Send,
	(*chanA).Recv,
	(*chanA).TrySend,
	(*chanA).TryRecv,
	(*chanA).Done,
	(*chanA).Close,
}
