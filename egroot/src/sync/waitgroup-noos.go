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
		if e := syscall.AtomicLoadEvent(&wg.e); e != 0 {
			e.Send()
		}
	}
}

func wait(wg *WaitGroup) {
	if wg.e == 0 {
		syscall.AtomicStoreEvent(&wg.e, syscall.AssignEvent())
	}
	for atomic.LoadInt32(&wg.cnt) != 0 {
		wg.e.Wait()
	}
}
