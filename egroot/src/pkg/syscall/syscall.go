package syscall

type Errno uintptr

func (e Errno) Error() string {
	if e == 0 {
		return "success"
	}
	if int(e) >= len(errnos) {
		return "unknown error"
	}
	return errnos[e]
}
