// +build linux

package sync

import (
	"sync/atomic"
	"syscall"
)

type mutex struct {
	f int32
}

const (
	futexWait = syscall.FUTEX_WAIT | syscall.FUTEX_PRIVATE_FLAG
	futexWake = syscall.FUTEX_WAKE | syscall.FUTEX_PRIVATE_FLAG
)

func lock(m *Mutex) {
	for {
		if atomic.CompareAndSwapInt32(&m.f, 0, 1) {
			return
		}
		syscall.Futex(&m.f, futexWait, 1, nil, nil, 0)
	}
}

func unlock(m *Mutex) {
	if atomic.AddInt32(&m.f, -1) != 0 {
		panic("sync: unlock of unlocked mutex")
	}
	syscall.Futex(&m.f, futexWake, 1, nil, nil, 0)
}
