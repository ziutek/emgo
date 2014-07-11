package noos

import (
	"builtin"
	"sync/atomic"
	"unsafe"
)

const (
	cok = builtin.ChanOK + iota
	cclosed
	cagain
)

func panicClosed() {
	panic("send: closed channel")
}

func panicCloseNil() {
	panic("close: nil channel")
}

func makeChan(cap int, size, align uintptr) (c builtin.Chan) {
	if cap == 0 {
		c.C = unsafe.Pointer(makeChanS())
		c.M = (*builtin.ChanMethods)(unsafe.Pointer(&csm))
	} else {
		c.C = unsafe.Pointer(makeChanA(cap, size, align))
		c.M = (*builtin.ChanMethods)(unsafe.Pointer(&cam))
	}
	return
}

type waiter struct {
	addr unsafe.Pointer
	next *waiter
}

func selectComm(comms []*builtin.Comm, dflt unsafe.Pointer) (jmp, p unsafe.Pointer, d uintptr) {
	// BUG: comms need to be shufled there. Waiting for rand package...

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
		e   Event
		sel int32
		w   waiter
	)
	for _, comm := range comms {
		if comm.C != nil {
			e = e.Sum(*(*Event)(comm.C))
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

func init() {
	builtin.MakeChan = makeChan
	builtin.Select = selectComm
}
