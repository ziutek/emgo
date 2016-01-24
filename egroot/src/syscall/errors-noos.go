// build +noos

package syscall

const (
	OK Errno = iota
	ENORES
	ENFOUND
	EPERM
	ERANGE
	ENOSYS
)

//emgo:const
var errnos = []string{
	OK:      "success",
	ENORES:  "no resources",
	ENFOUND: "not found",
	EPERM:   "no permissions",
	ERANGE:  "out of range",
	ENOSYS:  "function not implemented",
}
