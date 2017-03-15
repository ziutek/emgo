package semihosting

type Error int

func (err Error) Error() string {
	return "semihosting error"
}

type File struct {
	fd int
}

func (f File) Fd() int {
	return f.fd
}

type FopenMode int

const (
	R   FopenMode = 0
	RB  FopenMode = 1
	Rp  FopenMode = 2
	RpB FopenMode = 3
	W   FopenMode = 4
	WB  FopenMode = 5
	Wp  FopenMode = 6
	WpB FopenMode = 7
	A   FopenMode = 8
	AB  FopenMode = 9
	Ap  FopenMode = 10
	ApB FopenMode = 11
)

func OpenFile(name string, mode FopenMode) (File, error) {
	return openFile(name, mode)
}

func (f File) Close() error {
	return f.close()
}

func (f File) WriteString(s string) (int, error) {
	return f.writeString(s)
}

func (f File) Write(b []byte) (int, error) {
	return f.write(b)
}
