// +build noos

package builtin

const (
	NEWTASK = iota
	DELTASK
	TASKUNLOCK
)

func Syscall0(trap uintptr) (r, e uintptr)
func Syscall1(trap, a1 uintptr) (r, e uintptr)
func Syscall2(trap, a1, a2 uintptr) (r, e uintptr)