// +build noos

package time

import "runtime/noos"

var start Time

func set(t Time) {
	up := noos.Uptime()
	sec := int64(up / 1e9)
	nsec := int32(up - uint64(sec)*1e9)
	if nsec < 0 {
		sec--
		nsec += 1e9
	}
	start = Time{sec, nsec}
}

func now() (t Time) {
	up := noos.Uptime()
	sec := int64(up / 1e9)
	nsec := int32(up - uint64(sec)*1e9)
	t.sec = start.sec + sec
	t.nsec = start.nsec + nsec
	if t.nsec >= 1e9 {
		t.sec++
		t.nsec -= 1e9
	}
	return
}
