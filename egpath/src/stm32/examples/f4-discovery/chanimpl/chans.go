package main

import (
	"runtime/noos"
	"sync/atomic"
	"sync/barrier"
	"unsafe"
)

const (
	ready int32 = iota
	recvWait
	sendWait
	closed
)

type ChanS struct {
	event noos.Event
	src   unsafe.Pointer
	dst   unsafe.Pointer
	state int32
}

func NewChanS() *ChanS {
	c := new(ChanS)
	c.event = noos.AssignEvent()
	return c
}

func (c *ChanS) Close() {
	atomic.StoreInt32(&c.state, closed)
	c.event.Send()
}

func panicIfClosed(state int32) {
	if state == closed {
		panic("send: closed channel")
	}
}

func (c *ChanS) getDstAndUnlock() (unsafe.Pointer, uintptr) {
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

func (c *ChanS) getSrcAndUnlock() (unsafe.Pointer, uintptr) {
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
func (c *ChanS) Send(e unsafe.Pointer) (p unsafe.Pointer, d uintptr) {
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
func (c *ChanS) Recv(e unsafe.Pointer) (p unsafe.Pointer, d uintptr) {
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
func (c *ChanS) TrySend(e unsafe.Pointer) (p unsafe.Pointer, d uintptr) {
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
func (c *ChanS) TryRecv(e unsafe.Pointer) (p unsafe.Pointer, d uintptr) {
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
func (c *ChanS) Done(d uintptr) {
	barrier.Memory()
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(d)), nil)
	c.event.Send()
}