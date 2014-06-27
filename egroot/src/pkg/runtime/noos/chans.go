package noos

import (
	"sync/atomic"
	"sync/barrier"
	"unsafe"
)

// Synchronous channels

func loadWaiter(wptr **waiter) *waiter {
	return (*waiter)(atomic.LoadPointer(
		(*unsafe.Pointer)(unsafe.Pointer(wptr)),
	))
}

func addWaiter(head **waiter, w *waiter) {
	w.next = nil
	if *head == nil {
		*head = w
		return
	}
	last := (*head)
	for last.next != nil {
		last = last.next
	}
	last.next = w
}

func delWaiter(head **waiter, w *waiter) bool {
	if *head == nil {
		return false
	}
	if *head == w {
		*head = w.next
		return true
	}
	prev := *head
	for prev.next != nil {
		if prev.next == w {
			prev.next = w.next
			return true
		}
		prev = prev.next
	}
	return false
}

type chanS struct {
	event Event // Event must be the first field (see chanSelect).
	src   *waiter
	dst   *waiter
	state int32
}

func makeChanS() *chanS {
	c := new(chanS)
	c.event = AssignEvent()
	return c
}

func (c *chanS) Close() {
	if c == nil {
		panicCloseNil()
	}
	atomic.StoreInt32(&c.state, 2)
	c.event.Send()
}

func (c *chanS) isClosed() bool {
	return atomic.LoadInt32(&c.state) == 2
}

func (c *chanS) lock() bool {
	for !atomic.CompareAndSwapInt32(&c.state, 0, 1) {
		if c.isClosed() {
			return true
		}
		c.event.Wait()
	}
	return false
}

func (c *chanS) unlock() {
	barrier.Memory()
	atomic.CompareAndSwapInt32(&c.state, 1, 0)
	c.event.Send()
}

// TrySend tries to send the value of variable pointed by e.
// w is used to guarantee exclusive communication in select and to signal that
// data transfer was completed. After return, if p != nil the data transfer need
// to be performed by sender and after that sender need to call c.Done(d).
// Otherwise, d can be equal to cagain or cok. cagain means that c isn't ready
// for communication and TrySend can be called again. cok means that data
// transfer was performed by receiver.
func (c *chanS) TrySend(e unsafe.Pointer, w *waiter) (p unsafe.Pointer, d uintptr) {
	if w == nil {
		// Fast path.
		if loadWaiter(&c.dst) == nil {
			if c.isClosed() {
				panicClosed()
			}
			return nil, cagain
		}
	} else {
		if w.next != w {
			// Altready waiting for communication.
			if atomic.LoadPointer(&w.addr) == nil {
				// Receiver is ready.
				atomic.StorePointer(&w.addr, e)
				c.event.Send()
				// Wait for receiver to complete communication.
				for atomic.LoadPointer(&w.addr) != nil {
					c.event.Wait()
				}
				return nil, cok
			}
			return nil, cagain
		}
	}
	if c.lock() {
		panicClosed()
	}
	// Try comunicate with oldest waiting receiver.
	for c.dst != nil {
		dst := c.dst
		c.dst = dst.next
		// Try lock select.
		if !atomic.CompareAndSwapInt32((*int32)(dst.addr), 0, 1) {
			continue
		}
		// Inform receiver that sender is ready.
		dst.addr = nil
		c.unlock() // This wakeups receivers.
		// Wait for receiver's decision.
		for {
			if p = atomic.LoadPointer(&dst.addr); p != nil {
				break
			}
			c.event.Wait()
		}
		if p != unsafe.Pointer(&dst.addr) {
			// Receiver is ready for communication.
			return p, uintptr(unsafe.Pointer(&dst.addr))
		}
		// Receiver cancel communication.
		atomic.StorePointer(&dst.addr, nil)
		c.event.Send()
		if c.lock() {
			panicClosed()
		}
	}
	if w != nil {
		// Add itself to the list of waiting senders.
		addWaiter(&c.src, w)
	}
	c.unlock()
	return nil, cagain
}

// TryRecv tries to receive a value channel and store it into variable
// pointed by e.
// w is used to guarantee exclusive communication in select and to signal that
// data transfer was completed. After return, if p != nil the data transfer need
// to be performed by receiver and after that receiver need to call c.Done(d).
// Otherwise, d can be equal to cagain, cok or cclosed. cagain means that c
// isn't ready for communication and TryRecv can be called again. cok means
// that data transfer was performed by receiver. cclosed means that c is closed.
func (c *chanS) TryRecv(e unsafe.Pointer, w *waiter) (p unsafe.Pointer, d uintptr) {
	if w == nil {
		// Fast path.
		if loadWaiter(&c.src) == nil {
			if c.isClosed() {
				return nil, cclosed
			}
			return nil, cagain
		}
	} else {
		if w.next != w {
			// Altready waiting for communication.
			if atomic.LoadPointer(&w.addr) == nil {
				// Sender is ready.
				atomic.StorePointer(&w.addr, e)
				c.event.Send()
				// Wait for sender to complete communication.
				for atomic.LoadPointer(&w.addr) != nil {
					c.event.Wait()
				}
				return nil, cok
			}
			return nil, cagain
		}
	}
	if c.lock() {
		return nil, cclosed
	}
	// Try comunicate with oldest waiting sender.
	for c.src != nil {
		src := c.src
		c.src = src.next
		// Try lock select.
		if !atomic.CompareAndSwapInt32((*int32)(src.addr), 0, 1) {
			continue
		}
		// Inform sender that receiver is ready.
		src.addr = nil
		c.unlock() // This wakeups senders.
		// Wait for sender's decision.
		for {
			if p = atomic.LoadPointer(&src.addr); p != nil {
				break
			}
			c.event.Wait()
		}
		if p != unsafe.Pointer(&src.addr) {
			// Sender is ready for communication.
			return p, uintptr(unsafe.Pointer(&src.addr))
		}
		// Sender cancel communication.
		atomic.StorePointer(&src.addr, nil)
		c.event.Send()
		if c.lock() {
			return nil, cclosed
		}
	}
	if w != nil {
		// Add itself to the list of waiting receivers.
		addWaiter(&c.dst, w)
	}
	c.unlock()
	return nil, cagain
}

// Done must be called after the data transfer was completed.
func (c *chanS) Done(d uintptr) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(d)), nil)
	c.event.Send()
}

func (c *chanS) cancel(wptr **waiter, w *waiter) {
	c.lock()
	if !delWaiter(wptr, w) && atomic.LoadInt32((*int32)(w.addr)) != 2 {
		atomic.StorePointer(&w.addr, unsafe.Pointer(w.addr))
		for atomic.LoadPointer(&w.addr) != nil {
			c.event.Wait()
		}
	}
	c.unlock()
}

func (c *chanS) CancelSend(w *waiter) {
	c.cancel(&c.src, w)
}

func (c *chanS) CancelRecv(w *waiter) {
	c.cancel(&c.dst, w)
}

// Send tries to send the value of variable pointed by e.
// If p != nil the data transfer need to be performed by sender and after that
// sender need to call c.Done(d). Otherwise the data transfer was performed by
// receiver.
func (c *chanS) Send(e unsafe.Pointer) (p unsafe.Pointer, d uintptr) {
	var (
		sel int32
		w   waiter
	)
	w.addr = unsafe.Pointer(&sel)
	w.next = &w
	for {
		p, d = c.TrySend(e, &w)
		if p != nil || d != cagain {
			return
		}
		c.event.Wait()
	}
}

// Recv tries to receive a value from channel and store it into variable pointed
// by e.
// If p != nil the data transfer need to be performed by receiver and after that
// receiver need to call c.Done(d). Otherwise d == cok means that the data
// transfer was performed by sender, d == cclosed means that channel is closed
// and there was no data to receive.
func (c *chanS) Recv(e unsafe.Pointer) (p unsafe.Pointer, d uintptr) {
	var (
		sel int32
		w   waiter
	)
	w.addr = unsafe.Pointer(&sel)
	w.next = &w
	for {
		p, d = c.TryRecv(e, &w)
		if p != nil || d != cagain {
			return
		}
		c.event.Wait()
	}
}

func (_ *chanS) Len() int {
	return 0
}

func (_ *chanS) Cap() int {
	return 0
}

type chanSMethods struct {
	Send       func(c *chanS, e unsafe.Pointer) (p unsafe.Pointer, d uintptr)
	Recv       func(c *chanS, e unsafe.Pointer) (p unsafe.Pointer, d uintptr)
	TrySend    func(c *chanS, e unsafe.Pointer, w *waiter) (p unsafe.Pointer, d uintptr)
	TryRecv    func(c *chanS, e unsafe.Pointer, w *waiter) (p unsafe.Pointer, d uintptr)
	CancelSend func(c *chanS, w *waiter)
	CancelRecv func(c *chanS, w *waiter)
	Done       func(c *chanS, d uintptr)
	Close      func(c *chanS)
	Len        func(c *chanS) int
	Cap        func(c *chanS) int
}

var csm = chanSMethods{
	(*chanS).Send,
	(*chanS).Recv,
	(*chanS).TrySend,
	(*chanS).TryRecv,
	(*chanS).CancelSend,
	(*chanS).CancelRecv,
	(*chanS).Done,
	(*chanS).Close,
	(*chanS).Len,
	(*chanS).Cap,
}
