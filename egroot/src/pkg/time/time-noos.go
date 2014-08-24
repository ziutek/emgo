// +build noos

package time

import "runtime/noos"

var start Time

func Set(t Time) {
	up := noos.Uptime()
	sec := int64(up / 1e9)
	nsec := int32(up - uint64(sec)*1e9)
	if nsec < 0 {
		sec--
		nsec += 1e9
	}
	start = Time{sec: sec, nsec: nsec}
}

func now() (sec int64, nsec int32) {
	up := noos.Uptime()
	sec = int64(up / 1e9)
	nsec = int32(up - uint64(sec)*1e9)
	sec = start.sec + sec
	nsec = start.nsec + nsec
	if nsec >= 1e9 {
		sec++
		nsec -= 1e9
	}
	return
}
