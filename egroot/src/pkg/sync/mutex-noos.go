// +build noos

package sync

import (
	"runtime/noos"
	"sync/atomic"
	"sync/barrier"
)

type mutex struct {
	state uintptr
}

func (m *Mutex) lock() {
	if m.state == 0 {
		m.state = uintptr(noos.AssignEventFlag())
	}
	unlocked, locked := m.state&^1, m.state|1
	for {
		if atomic.CompareAndSwapUintptr(&m.state, unlocked, locked) {
			return
		}
		noos.Event(unlocked).Wait()
	}
}

func (m *Mutex) unlock() {
	unlocked := m.state &^ 1
	barrier.Memory()
	if atomic.AddUintptr(&m.state, ^uintptr(0)) != unlocked {
		panic("sync: unlock of unlocked mutex")
	}
	noos.Event(unlocked).Send()
}
