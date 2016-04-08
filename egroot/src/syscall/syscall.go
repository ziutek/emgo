package syscall

type Errno int

func (e Errno) Error() string {
	if int(e) >= len(errnos) {
		return "unknown error"
	}
	return errnos[e]
}