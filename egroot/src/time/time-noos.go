// +build noos

package time

import (
	"rtos"
	"sync"
	"sync/atomic"
	"sync/fence"
)

type sysStart struct {
	n uintptr
	t [2]Time
	sync.Mutex
}

var start sysStart

// Set sets the current time. Local is set to t.Location().
func Set(t Time) {
	ns := rtos.Nanosec()
	t.sec -= ns / 1e9
	t.nsec -= int32(ns % 1e9)
	if t.nsec < 0 {
		t.sec--
		t.nsec += 1e9
	}
	start.Lock()
	n := start.n + 1
	start.t[n&1] = t
	fence.W_SMP()
	start.n = n
	Local = t.loc
	start.Unlock()
}

func now() Time {
	var t Time
	n := atomic.LoadUintptr(&start.n)
	for {
		fence.R_SMP()
		t = start.t[n&1]
		fence.R_SMP()
		m := atomic.LoadUintptr(&start.n)
		if m == n {
			break
		}
		n = m
	}
	ns := rtos.Nanosec()
	t.sec += ns / 1e9
	t.nsec += int32(ns % 1e9)
	if t.nsec >= 1e9 {
		t.sec++
		t.nsec -= 1e9
	}
	return t
}
