package errors

type strerr string

func New(s string) error {
	return strerr(s)
}

func (e strerr) Error() string {
	return string(e)
}
