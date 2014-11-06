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
	TASKREADY
	WAITEVENT
)

func syscall0(trap uintptr) (uintptr, Errno)

func syscall1(trap, a1 uintptr) (uintptr, Errno)

func syscall2(trap, a1, a2 uintptr) (uintptr, Errno)

// NewTask creates new task that starts execution at pc. If wait is true
// tasker stops and waits until new task will call TaskReady.
func NewTask(pc uintptr, wait bool) Errno {
	var w uintptr
	if wait {
		w = 1
	}
	_, err := syscall2(NEWTASK, pc, w)
	return err
}
