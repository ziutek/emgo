// +build noos

package sync

import (
	"syscall"
	"sync/atomic"
	"sync/fence"
)

type mutex struct {
	state uintptr
}

func (m *Mutex) lock() {
	state := atomic.LoadUintptr(&m.state)
	if state == 0 {
		state = uintptr(syscall.AssignEventFlag())
		if !atomic.CompareAndSwapUintptr(&m.state, 0, state) {
			state = m.state
		}
	}
	unlocked, locked := state&^1, state|1
	for {
		if atomic.CompareAndSwapUintptr(&m.state, unlocked, locked) {
			return
		}
		syscall.Event(unlocked).Wait()
	}
}

func (m *Mutex) unlock() {
	unlocked := m.state &^ 1
	fence.Memory()
	if atomic.AddUintptr(&m.state, ^uintptr(0)) != unlocked {
		panic("sync: unlock of unlocked mutex")
	}
	syscall.Event(unlocked).Send()
}
