// +build noos

package sync

import (
	"sync/atomic"
	"sync/fence"
	"syscall"
)

type waitgroup struct {
	e   syscall.Event
	cnt int
}

func (wg *WaitGroup) add(delta int) {
	fence.RW_SMP() // Add used as RELEASE (delta < 0).
	cnt := atomic.AddInt(&wg.cnt, delta)
	if cnt < 0 || cnt < delta {
		panic("sync: negative WaitGroup counter")
	}
	if delta < 0 && cnt == 0 {
		fence.RW_SMP() // Store(wg.cnt) must be observed before load(wg.e).
		if e := syscall.AtomicLoadEvent(&wg.e); e != 0 {
			e.Send()
			return
		}
	}
	fence.RW_SMP() // Add used as ACQUIRE (delta > 0).
}

func (wg *WaitGroup) wait() {
	if wg.e == 0 {
		syscall.AtomicStoreEvent(&wg.e, syscall.AssignEvent())
	}
	fence.RW_SMP() // Store(wg.e) must be observed before load(wg.cnt).
	for atomic.LoadInt(&wg.cnt) != 0 {
		wg.e.Wait()
	}
}
