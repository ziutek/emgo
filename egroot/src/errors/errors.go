package errors

type strErr struct {
	s string
}

func (e strErr) Error() string {
	return e.s
}

func New(s string) error {
	return strErr{s}
}
