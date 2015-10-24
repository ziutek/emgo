package rtos

import (
	"syscall"
)

// Debug allows write debug information.
type Debug int

func (d Debug) Write(b []byte) (int, error) {
	n, e := syscall.DebugOut(int(d), b)
	return n, mkerror(e)
}

func (d Debug) WriteByte(b byte) error {
	_, e := syscall.DebugOut(int(d), []byte{b})
	return mkerror(e)
}

func (d Debug) WriteString(s string) (int, error) {
	n, e := syscall.DebugOutString(int(d), s)
	return n, mkerror(e)
}
