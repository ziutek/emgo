// +build linux

package os

import "syscall"

func Exit(code int) {
	syscall.Exit(code)
}
