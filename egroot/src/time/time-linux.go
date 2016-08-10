// +build linux

package time

import (
	"syscall"
)

func now() Time {
	var tp syscall.Timespec
	syscall.ClockGettime(syscall.CLOCK_REALTIME, &tp)
	return Time{sec: tp.Sec + unixToInternal, nsec: int32(tp.Nsec)}
}
