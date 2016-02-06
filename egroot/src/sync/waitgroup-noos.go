// +build noos

package sync

import (
	"sync/atomic"
	"syscall"
)

type waitgroup struct {
	e   syscall.Event
	cnt int32
}

func add(wg *WaitGroup, delta int) {
	cnt := atomic.AddInt32(&wg.cnt, int32(delta))
	if cnt < 0 {
		panic("sync: negative WaitGroup counter")
	}
	if delta < 0 && cnt == 0 {
		e := syscall.AtomicLoadEvent(&wg.e)
		if e != 0 {
			// Waiter should check cnt.
			e.Send()
		}
	}
}

func wait(wg *WaitGroup) {
	e := syscall.AssignEvent()
	syscall.AtomicStoreEvent(&wg.e, e)
	for atomic.LoadInt32(&wg.cnt) != 0 {
		e.Wait()
	}
}
