package main

import (
	"runtime/noos"
	"sync/atomic"
	"sync/barrier"
)

func bit(w *uint32, n uint32) bool {
	barrier.Compiler()
	return *w&(1<<n) != 0
}

func atomicToggleBit(addr *uint32, n uint32) {
	mask := uint32(1) << n
	for {
		old := atomic.LoadUint32(addr)
		if atomic.CompareAndSwapUint32(addr, old, old^mask) {
			return
		}
	}
}

/*func atomicClearBit(addr *uint32, n uint32) {
	mask := uint32(1) << n
	for {
		old := atomic.LoadUint32(addr)
		if atomic.CompareAndSwapUint32(addr, old, old&^mask) {
			return
		}
	}
}*/

const ccap = 4

type ChanA struct {
	event  noos.Event
	cap    uint32
	tosend uint32
	torecv uint32
	rd     []uint32
}

func NewChanA(cap int) *ChanA {
	c := new(ChanA)
	c.event = noos.AssignEvent()
	c.cap = uint32(cap)
	c.rd = make([]uint32, 1+cap/32)
	return c
}

const b31 = 1 << 31

func (c *ChanA) TrySend() int {
	for {
		tosend := atomic.LoadUint32(&c.tosend)
		if atomic.LoadUint32(&c.torecv) == tosend^b31 {
			// Channel is full.
			return -1
		}
		n := tosend &^ b31

		if bit(&c.rd[n>>5], n&31) {
			// This element is still being received.
			return -1
		}

		var next uint32

		if n+1 == c.cap {
			next = ^tosend & b31
		} else {
			next = tosend + 1
		}
		if atomic.CompareAndSwapUint32(&c.tosend, tosend, next) {
			return int(n)
		}
	}
}

func (c *ChanA) Send() int {
	for {
		if n := c.TrySend(); n != -1 {
			return n
		}
		c.event.Wait()
	}
}

func (c *ChanA) TryRecv() int {
	for {
		torecv := atomic.LoadUint32(&c.torecv)
		if atomic.LoadUint32(&c.tosend) == torecv {
			// Channel is empty.
			return -1
		}
		n := torecv &^ b31
		if !bit(&c.rd[n>>5], n&31) {
			// This element is still being sent.
			return -1
		}

		var next uint32

		if n+1 == c.cap {
			next = ^torecv & b31
		} else {
			next = torecv + 1
		}
		if atomic.CompareAndSwapUint32(&c.torecv, torecv, next) {
			return int(n)
		}
	}
}

func (c *ChanA) Recv() int {
	for {
		if n := c.TryRecv(); n != -1 {
			return n
		}
		c.event.Wait()
	}
}

func (c *ChanA) Done(n int) {
	barrier.Memory()
	atomicToggleBit(&c.rd[n>>5], uint32(n)&31)
	c.event.Send()
}
