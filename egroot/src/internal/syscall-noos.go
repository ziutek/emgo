// +build noos

package internal

const (
	NEWTASK = iota
	KILLTASK
	TASKUNLOCK
)

func Syscall0(trap uintptr) (r, e uintptr)
func Syscall1(trap, a1 uintptr) (r, e uintptr)
func Syscall2(trap, a1, a2 uintptr) (r, e uintptr)
func Syscall3(trap, a1, a2, a3 uintptr) (r, e uintptr)

func Syscall1i64(trap, uintptr int64) (r, e uintptr)
func Syscall0r64(trap uintptr) int64

//c:inline
func NewTask(f func(), lock bool)
