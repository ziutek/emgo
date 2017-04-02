// +build linux

package internal

var (
	Argv uintptr
	Env  uintptr
	Auxv uintptr
)

func Syscall1(trap, a1 uintptr) uintptr
func Syscall2(trap, a1, a2 uintptr) uintptr
func Syscall3(trap, a1, a2, a3 uintptr) uintptr
func Syscall4(trap, a1, a2, a3, a4 uintptr) uintptr
func Syscall5(trap, a1, a2, a3, a4, a5 uintptr) uintptr
func Syscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) uintptr
