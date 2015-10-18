package rtos

import "syscall"

func mkerror(e syscall.Errno) error {
	if e == syscall.OK {
		return nil
	}
	return e
}