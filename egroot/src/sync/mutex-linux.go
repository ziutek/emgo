// +build linux

package sync

import (
	"sync/atomic"
	"sync/fence"
	"syscall"
)

type mutex struct {
	f int32
}

const (
	futexWait = syscall.FUTEX_WAIT | syscall.FUTEX_PRIVATE_FLAG
	futexWake = syscall.FUTEX_WAKE | syscall.FUTEX_PRIVATE_FLAG
)

func (m *Mutex) lock() {
	for {
		if atomic.CompareAndSwapInt32(&m.f, 0, 1) {
			break
		}
		syscall.Futex(&m.f, futexWait, 1, nil, nil, 0)
	}
	fence.RW_SMP() // Lock has ACQUIRE semantic.
}

func (m *Mutex) unlock() {
	fence.RW_SMP() // Unlock has RELEASE semantic.
	if atomic.AddInt32(&m.f, -1) != 0 {
		panic("sync: unlock of unlocked mutex")
	}
	syscall.Futex(&m.f, futexWake, 1, nil, nil, 0)
}
