package linux

import (
	"internal"
	"syscall"
)

func init() {
	internal.Panic = panic_
	internal.Alloc = alloc
}

func exit(code int) {
	syscall.Exit(code)
}
