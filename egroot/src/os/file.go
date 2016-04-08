// +build linux

package os

import (
	"syscall"
)

type FileMode uint32

const (
	O_RDONLY int = syscall.O_RDONLY // open the file read-only.
	O_WRONLY int = syscall.O_WRONLY // open the file write-only.
	O_RDWR   int = syscall.O_RDWR   // open the file read-write.
	O_APPEND int = syscall.O_APPEND // append data to the file when writing.
	O_CREATE int = syscall.O_CREAT  // create a new file if none exists.
	O_EXCL   int = syscall.O_EXCL   // used with O_CREATE, file must not exist
	O_SYNC   int = syscall.O_SYNC   // open for synchronous I/O.
	O_TRUNC  int = syscall.O_TRUNC  // if possible, truncate file when opened.
)

type File struct {
	fd int
}

var (
	Stdin  = File{0}
	Stdout = File{1}
	Stderr = File{2}
)

func OpenFile(name string, flag int, perm FileMode) (File, error) {
	fd, err := syscall.Open(name, flag, uint32(perm))
	if err != nil {
		return File{-1}, err
	}
	return File{fd}, nil
}

func Open(name string) (File, error) {
	return OpenFile(name, O_RDONLY, 0)
}

func Create(name string) (File, error) {
	return OpenFile(name, O_RDWR|O_CREATE|O_TRUNC, 0666)
}

func (f File) WriteString(s string) (int, error) {
	return syscall.WriteString(f.fd, s)
}

func (f File) Write(b []byte) (int, error) {
	return syscall.Write(f.fd, b)
}

func (f File) Read(b []byte) (int, error) {
	return syscall.Read(f.fd, b)
}

func (f File) Close() error {
	return syscall.Close(f.fd)
}
