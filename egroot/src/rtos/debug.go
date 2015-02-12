package rtos

import (
	"syscall"
)

// Debug allows write debug information.
type Debug int

func (d Debug) Write(b []byte) (int, error) {
	return syscall.DebugOut(int(d), b)
}

func (d Debug) WriteString(s string) (int, error) {
	return syscall.DebugOutString(int(d), s)
}

const (
	DbgOut Debug = 1
	DbgErr Debug = 2
)
