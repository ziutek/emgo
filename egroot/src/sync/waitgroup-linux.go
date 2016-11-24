// +build linux

package sync

import (
	"sync/atomic"
	"sync/fence"
	"syscall"
)

type waitgroup struct {
	cnt int
	f   int32
}

func (wg *WaitGroup) add(delta int) {
	fence.RW_SMP() // Add used as RELEASE (delta < 0).
	cnt := atomic.AddInt(&wg.cnt, delta)
	if cnt < 0 || cnt < delta {
		panic("sync: negative WaitGroup counter")
	}
	switch cnt {
	case 0: // cnt == 0 implies delta < 0.
		atomic.StoreInt32(&wg.f, 1)
		syscall.Futex(&wg.f, futexWake, 1<<31-1, nil, nil, 0)
	case delta:
		if delta > 0 {
			// cnt == delta && delta > 0 implies cnt == 0 before add.
			atomic.StoreInt32(&wg.f, 0) // To allow reuse wg.
		}
	}
	fence.RW_SMP() // Add used as ACQUIRE (delta > 0).
}

func (wg *WaitGroup) wait() {
	for atomic.LoadInt(&wg.cnt) == 0 {
		syscall.Futex(&wg.f, futexWait, 1, nil, nil, 0)
	}
}
