// +build noos

package time

import "rtos"

var start Time

// SetStart sets start time of rtos system clock.
func SetStart(t Time) {
	start = t
}

func now() (sec int64, nsec int32) {
	up := rtos.Nanosec()
	sec = up / 1e9
	nsec = int32(up - sec*1e9)
	sec = start.sec + sec
	nsec = start.nsec + nsec
	if nsec >= 1e9 {
		sec++
		nsec -= 1e9
	}
	return
}
