package syscall

type Errno uintptr

func (e Errno) Error() string {
	if int(e) >= len(errnos) {
		return "unknown error"
	}
	return errnos[e]
}