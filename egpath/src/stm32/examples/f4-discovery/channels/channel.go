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

func atomicSetBit(addr *uint32, n uint32) {
	mask := uint32(1) << n
	for {
		old := atomic.LoadUint32(addr)
		if atomic.CompareAndSwapUint32(addr, old, old|mask) {
			return
		}
	}
}

func atomicClearBit(addr *uint32, n uint32) {
	mask := uint32(1) << n
	for {
		old := atomic.LoadUint32(addr)
		if atomic.CompareAndSwapUint32(addr, old, old&^mask) {
			return
		}
	}
}

type Elem int

const ccap = 4

type Chan struct {
	event  noos.Event
	tosend uint32
	torecv uint32
	rd     [1 + ccap/32]uint32
	buf    [ccap]Elem
}

const b31 = 1 << 31

func (c *Chan) TrySend(e Elem) bool {
	var n, wordn, bitn uint32

	for {
		tosend := atomic.LoadUint32(&c.tosend)
		if atomic.LoadUint32(&c.torecv) == tosend^b31 {
			// Channel is full.
			return false
		}
		n = tosend &^ b31
		wordn = n >> 5
		bitn = n & 31

		if bit(&c.rd[wordn], bitn) {
			// This element is still being received.
			return false
		}

		var next uint32

		if n+1 == uint32(len(c.buf)) {
			next = ^tosend & b31
		} else {
			next = tosend + 1
		}
		if atomic.CompareAndSwapUint32(&c.tosend, tosend, next) {
			break
		}
	}

	c.buf[n] = e
	barrier.Memory()

	atomicSetBit(&c.rd[wordn], bitn)
	c.event.Send()

	return true
}

func (c *Chan) Send(e Elem) {
	for !c.TrySend(e) {
		c.event.Wait()
	}
}

func (c *Chan) TryRecv() (Elem, bool) {
	var n, wordn, bitn uint32

	for {
		torecv := atomic.LoadUint32(&c.torecv)
		if atomic.LoadUint32(&c.tosend) == torecv {
			// Channel is empty.
			return 0, false
		}
		n = torecv &^ b31
		wordn = n >> 5
		bitn = n & 31
		if !bit(&c.rd[wordn], bitn) {
			// This element is still being sent.
			return 0, false
		}

		var next uint32

		if n+1 == uint32(len(c.buf)) {
			next = ^torecv & b31
		} else {
			next = torecv + 1
		}
		if atomic.CompareAndSwapUint32(&c.torecv, torecv, next) {
			break
		}
	}

	e := c.buf[n]
	barrier.Memory()

	atomicClearBit(&c.rd[wordn], bitn)
	c.event.Send()

	return e, true
}

func (c *Chan) Recv() Elem {
	for {
		if e, ok := c.TryRecv(); ok {
			return e
		}
		c.event.Wait()
	}
}
