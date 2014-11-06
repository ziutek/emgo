// build +noos

package syscall

const (
	ENORES Errno = iota + 1
)

var errnos = []string{
	ENORES: "no resources",
}

const (
	NEWTASK = iota
	DELTASK
	NEXTTASK
	WAITEVENT
)

func Syscall0(trap uintptr) (uintptr, Errno)

func Syscall1(trap, a1 uintptr) (uintptr, Errno)

func NewTask(wait bool)
