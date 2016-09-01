package noos

import (
	"internal"
	"sync/atomic"
	"syscall"
	"unsafe"
)

const (
	cok = internal.ChanOK + iota
	cclosed
	cagain
)

func panicClosed() {
	panic("send: closed channel")
}

func panicCloseNil() {
	panic("close: nil channel")
}

func makeChan(cap int, size, align uintptr) (c internal.Chan) {
	if cap == 0 {
		c.C = unsafe.Pointer(makeChanS())
		c.M = (*internal.ChanMethods)(unsafe.Pointer(&csm))
	} else {
		c.C = unsafe.Pointer(makeChanA(cap, size, align))
		c.M = (*internal.ChanMethods)(unsafe.Pointer(&cam))
	}
	return
}

type waiter struct {
	addr unsafe.Pointer
	next *waiter
}

func shuffle(comms []*internal.Comm) {
	rng := &tasker.tasks[tasker.curTask].rng
	n := uint(len(comms))
	for n > 1 {
		i := uint(rng.Uint64()) % n
		n--
		if i != n {
			comms[i], comms[n] = comms[n], comms[i]
		}
	}
	// TODO: use len(comms) to do this more efficently. Divide value from
	// rng.Uint64() into smaller chunks of nonzero bits. This reduces number of
	// rng.Uint64() calls and will result in fasetr % operation.
}

func selectComm(comms []*internal.Comm, dflt unsafe.Pointer) (jmp, p unsafe.Pointer, d uintptr) {
	shuffle(comms)

	if dflt != nil {
		// "Nonblocking" select.
		for _, comm := range comms {
			if comm.C == nil {
				continue
			}
			p, d = comm.Try(comm.C, comm.E, nil)
			if p != nil || d != cagain {
				jmp = comm.Case
				return
			}
		}
		jmp = dflt
		return
	}
	// Blocking select.
	var (
		e   syscall.Event
		sel int32
		w   waiter
	)
	for _, comm := range comms {
		if comm.C != nil {
			e |= *(*syscall.Event)(comm.C)
		}
	}
	w.addr = unsafe.Pointer(&sel)
	w.next = &w
	n := 0
	for {
		comm := comms[n]
		if comm.C != nil {
			p, d = comm.Try(comm.C, comm.E, unsafe.Pointer(&w))
			if p != nil || d != cagain {
				jmp = comm.Case
				break
			}
		}
		if n++; n == len(comms) {
			n = 0
			e.Wait()
		}
	}
	atomic.CompareAndSwapInt32(&sel, 0, 2)
	for i, comm := range comms {
		if i != n && comm.C != nil && comm.Cancel != nil {
			comm.Cancel(comm.C, unsafe.Pointer(&w))
		}
	}
	return
}
