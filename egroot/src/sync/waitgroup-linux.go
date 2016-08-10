// +build linux

package sync

import (
	"sync/atomic"
	"syscall"
)

type waitgroup struct {
	cnt int
	f   int32
}

func add(wg *WaitGroup, delta int) {
	cnt := atomic.AddInt(&wg.cnt, delta)
	if cnt < 0 {
		panic("sync: negative WaitGroup counter")
	}
	switch cnt {
	case 0:
		if delta < 0 {
			atomic.StoreInt32(&wg.f, 1)
			syscall.Futex(&wg.f, futexWake, 1<<31-1, nil, nil, 0)
		}
	case delta:
		if delta > 0 {
			// To allow reuse wg.
			atomic.StoreInt32(&wg.f, 0)
		}
	}
}

func wait(wg *WaitGroup) {
	for atomic.LoadInt(&wg.cnt) == 0 {
		syscall.Futex(&wg.f, futexWait, 1, nil, nil, 0)
	}
}
