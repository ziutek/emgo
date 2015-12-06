// +build noos

package builtin

const (
	NEWTASK = iota
	KILLTASK
	TASKUNLOCK
)

// TODO: All following functions tell compiler that they can (as side effect)
// modify memory. Consider provide second kind of syscalls that guarantee
// no memory modification to allow compiler better optimize generated code.

func Syscall0(trap uintptr) (r, e uintptr)
func Syscall1(trap, a1 uintptr) (r, e uintptr)
func Syscall2(trap, a1, a2 uintptr) (r, e uintptr)
func Syscall3(trap, a1, a2, a3 uintptr) (r, e uintptr)

func Syscall1i64(trap, uintptr uint64)
func Syscall0r64(trap uintptr) uint64
