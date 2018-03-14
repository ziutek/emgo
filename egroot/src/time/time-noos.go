// +build noos

package time

import (
	"rtos"
)

var start struct {
	sec  int64
	nsec int32
}

// Set sets the current time. Ns should be the value of rtos.Nanosec()
// corresponding to t. Local is set to t.Location(). Set causes discontinuites
// in time flow mesured by Now and is not thread-safe.
func Set(t Time, ns int64) {
	if ns != 0 {
		t.sec -= ns / 1e9
		t.nsec -= int32(ns % 1e9)
		if t.nsec < 0 {
			t.sec--
			t.nsec += 1e9
		}
	}
	start.sec = t.sec
	start.nsec = t.nsec
	Local = t.Location()
}

func now() (sec int64, nsec int32) {
	ns := rtos.Nanosec()
	sec = start.sec + ns/1e9
	nsec = start.nsec + int32(ns%1e9)
	if nsec >= 1e9 {
		sec++
		nsec -= 1e9
	}
	return
}
