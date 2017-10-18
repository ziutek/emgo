// Package semihosting provieds access to files located on debuging host
// (debuger must support it).
package semihosting

type Error int

func (err Error) Error() string {
	return "semihosting error"
}

// File represents an open file on debuging host.
type File struct {
	fd int
}

// Fd returns the integer Unix file descriptor referencing the open file.
func (f File) Fd() int {
	return f.fd
}

// FopenMode describes file read/write mode. It is mimics the meaningo of mode
// parameter of C standard library fopen function.
type FopenMode int

const (
	R   FopenMode = 0  // Open text file for reading from beggining of file.
	RB  FopenMode = 1  // Like R but for binary file.
	Rp  FopenMode = 2  // Open text file for read/writing at beggining of file.
	RpB FopenMode = 3  // Like Rp but for binary file.
	W   FopenMode = 4  // Truncate or create text file for writing.
	WB  FopenMode = 5  // Like W but for binary file.
	Wp  FopenMode = 6  // Truncate or create text file for writing and reading.
	WpB FopenMode = 7  // Like Wp but for binary file.
	A   FopenMode = 8  // Open or create text file for appending.
	AB  FopenMode = 9  // Like A but for binary file.
	Ap  FopenMode = 10 // Open or create text file for appending. and reading.
	ApB FopenMode = 11 // Like Ap but for binary file.
)

// OpenFile opens the named file for operations specified by mode. Use name
// ":tt" to read/write from/to standard input/output.
func OpenFile(name string, mode FopenMode) (File, error) {
	return openFile(name, mode)
}

// Close closes file.
func (f File) Close() error {
	return f.close()
}

// WriteString works like Write but allows to write strings.
func (f File) WriteString(s string) (int, error) {
	return f.writeString(s)
}

// Write writes len(s) bytes to the File. It returns the number of bytes written
// and an error, if any. Write returns a non-nil error when n != len(b).
func (f File) Write(b []byte) (int, error) {
	return f.write(b)
}
