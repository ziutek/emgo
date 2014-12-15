// build +noos

package syscall

const (
	OK Errno = iota
	ENORES
	ENFOUND
	EPERM
)

var errnos = []string{
	OK:      "success",
	ENORES:  "no resources",
	ENFOUND: "not found",
	EPERM:   "no permissions",
}
