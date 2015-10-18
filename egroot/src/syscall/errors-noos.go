// build +noos

package syscall

const (
	OK Errno = iota
	ENORES
	ENFOUND
	EPERM
	ERANGE
)

var errnos = []string{
	OK:      "success",
	ENORES:  "no resources",
	ENFOUND: "not found",
	EPERM:   "no permissions",
	ERANGE:  "out of range",
}
