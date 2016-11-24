// +build noos

package sync

import (
	"sync/fence"
	"syscall"
)

type mutex struct {
	state syscall.Event
}

func (m *Mutex) lock() {
	state := syscall.AtomicLoadEvent(&m.state)
	if state == 0 {
		state = syscall.AssignEventFlag() | 1
		if !syscall.AtomicCompareAndSwapEvent(&m.state, 0, state) {
			state = syscall.AtomicLoadEvent(&m.state)
		}
	}
	unlocked, locked := state|1, state&^1
	for {
		if syscall.AtomicCompareAndSwapEvent(&m.state, unlocked, locked) {
			break
		}
		locked.Wait()
	}
	fence.RW_SMP() // Lock has ACQUIRE semantic.
}

func (m *Mutex) unlock() {
	fence.RW_SMP() // Unlock has RELEASE semantic.
	state := syscall.AtomicLoadEvent(&m.state)
	if syscall.AtomicAddEvent(&m.state, 1) != state|1 {
		panic("sync: unlock of unlocked mutex")
	}
	state.Send()
}
