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
	state := atomic.LoadUintptr(&m.state)
	if state == 0 {
		state = uintptr(noos.AssignEventFlag())
		if !atomic.CompareAndSwapUintptr(&m.state, 0, state) {
			state = m.state
		}
	}
	unlocked, locked := state&^1, state|1
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
