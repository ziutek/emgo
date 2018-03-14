// +build noos

package time

import (
	"rtos"
	"sync"
	"sync/atomic"
	"sync/fence"
)

type timeStamp struct {
	sec  int64
	nsec int32
}

type sysStart struct {
	n  uintptr
	ts [2]timeStamp
	sync.Mutex
}

var start sysStart

// Set sets the current time. Ns should be the value of rtos.Nanosec()
// corresponding to t. Local is set to t.Location().
func Set(t Time, ns int64) {
	if ns != 0 {
		t.sec -= ns / 1e9
		t.nsec -= int32(ns % 1e9)
		if t.nsec < 0 {
			t.sec--
			t.nsec += 1e9
		}
	}
	start.Lock()
	n := start.n + 1
	start.ts[n&1] = timeStamp{t.sec, t.nsec}
	fence.W_SMP()
	start.n = n
	Local = t.Location()
	start.Unlock()
}

func now() (int64, int32) {
	var ts timeStamp
	n := atomic.LoadUintptr(&start.n)
	for {
		fence.R_SMP()
		ts = start.ts[n&1]
		fence.R_SMP()
		m := atomic.LoadUintptr(&start.n)
		if m == n {
			break
		}
		n = m
	}
	ns := rtos.Nanosec()
	ts.sec += ns / 1e9
	ts.nsec += int32(ns % 1e9)
	if ts.nsec >= 1e9 {
		ts.sec++
		ts.nsec -= 1e9
	}
	return ts.sec, ts.nsec
}
