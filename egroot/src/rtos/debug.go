package rtos

import (
	"syscall"
)

// Debug allows write debug information.
type Debug int

func (d Debug) Write(b []byte) (int, error) {
	return syscall.DebugOut(int(d), b)
}

func (d Debug) WriteByte(b byte) error {
	_, err := syscall.DebugOut(int(d), []byte{b})
	return err
}

func (d Debug) WriteString(s string) (int, error) {
	return syscall.DebugOutString(int(d), s)
}
