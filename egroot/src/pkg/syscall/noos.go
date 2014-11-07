// build +noos

package syscall

const (
	ENORES Errno = iota + 1 // No resources.
)

var errnos = []string{
	ENORES: "no resources",
}

const (
	NEWTASK = iota
	DELTASK
	TASKREADY
	WAITEVENT
)

func syscall0(trap uintptr) (uintptr, Errno)

func syscall1(trap, a1 uintptr) (uintptr, Errno)

func syscall2(trap, a1, a2 uintptr) (uintptr, Errno)

// NewTask creates new task that starts execute f. If wait is true
// tasker stops and waits until new task will call TaskReady.
func NewTask(f func(), wait bool) Errno {
	_, err := syscall2(NEWTASK, f2p(f), b2p(wait))
	return err
}
