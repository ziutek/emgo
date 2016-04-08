// +build linux

package syscall

const (
	EPERM           Errno = 1
	ENOENT          Errno = 2
	ESRCH           Errno = 3
	EINTR           Errno = 4
	EIO             Errno = 5
	ENXIO           Errno = 6
	E2BIG           Errno = 7
	ENOEXEC         Errno = 8
	EBADF           Errno = 9
	ECHILD          Errno = 10
	EAGAIN          Errno = 11
	ENOMEM          Errno = 12
	EACCES          Errno = 13
	EFAULT          Errno = 14
	ENOTBLK         Errno = 15
	EBUSY           Errno = 16
	EEXIST          Errno = 17
	EXDEV           Errno = 18
	ENODEV          Errno = 19
	ENOTDIR         Errno = 20
	EISDIR          Errno = 21
	EINVAL          Errno = 22
	ENFILE          Errno = 23
	EMFILE          Errno = 24
	ENOTTY          Errno = 25
	ETXTBSY         Errno = 26
	EFBIG           Errno = 27
	ENOSPC          Errno = 28
	ESPIPE          Errno = 29
	EROFS           Errno = 30
	EMLINK          Errno = 31
	EPIPE           Errno = 32
	EDOM            Errno = 33
	ERANGE          Errno = 34
	EDEADLK         Errno = 35
	ENAMETOOLONG    Errno = 36
	ENOLCK          Errno = 37
	ENOSYS          Errno = 38
	ENOTEMPTY       Errno = 39
	ELOOP           Errno = 40
	EWOULDBLOCK     Errno = 41
	ENOMSG          Errno = 42
	EIDRM           Errno = 43
	ECHRNG          Errno = 44
	EL2NSYNC        Errno = 45
	EL3HLT          Errno = 46
	EL3RST          Errno = 47
	ELNRNG          Errno = 48
	EUNATCH         Errno = 49
	ENOCSI          Errno = 50
	EL2HLT          Errno = 51
	EBADE           Errno = 52
	EBADR           Errno = 53
	EXFULL          Errno = 54
	ENOANO          Errno = 55
	EBADRQC         Errno = 56
	EBADSLT         Errno = 57
	EDEADLOCK       Errno = 58
	EBFONT          Errno = 59
	ENOSTR          Errno = 60
	ENODATA         Errno = 61
	ETIME           Errno = 62
	ENOSR           Errno = 63
	ENONET          Errno = 64
	ENOPKG          Errno = 65
	EREMOTE         Errno = 66
	ENOLINK         Errno = 67
	EADV            Errno = 68
	ESRMNT          Errno = 69
	ECOMM           Errno = 70
	EPROTO          Errno = 71
	EMULTIHOP       Errno = 72
	EDOTDOT         Errno = 73
	EBADMSG         Errno = 74
	EOVERFLOW       Errno = 75
	ENOTUNIQ        Errno = 76
	EBADFD          Errno = 77
	EREMCHG         Errno = 78
	ELIBACC         Errno = 79
	ELIBBAD         Errno = 80
	ELIBSCN         Errno = 81
	ELIBMAX         Errno = 82
	ELIBEXEC        Errno = 83
	EILSEQ          Errno = 84
	ERESTART        Errno = 85
	ESTRPIPE        Errno = 86
	EUSERS          Errno = 87
	ENOTSOCK        Errno = 88
	EDESTADDRREQ    Errno = 89
	EMSGSIZE        Errno = 90
	EPROTOTYPE      Errno = 91
	ENOPROTOOPT     Errno = 92
	EPROTONOSUPPORT Errno = 93
	ESOCKTNOSUPPORT Errno = 94
	EOPNOTSUPP      Errno = 95
	EPFNOSUPPORT    Errno = 96
	EAFNOSUPPORT    Errno = 97
	EADDRINUSE      Errno = 98
	EADDRNOTAVAIL   Errno = 99
	ENETDOWN        Errno = 100
	ENETUNREACH     Errno = 101
	ENETRESET       Errno = 102
	ECONNABORTED    Errno = 103
	ECONNRESET      Errno = 104
	ENOBUFS         Errno = 105
	EISCONN         Errno = 106
	ENOTCONN        Errno = 107
	ESHUTDOWN       Errno = 108
	ETOOMANYREFS    Errno = 109
	ETIMEDOUT       Errno = 110
	ECONNREFUSED    Errno = 111
	EHOSTDOWN       Errno = 112
	EHOSTUNREACH    Errno = 113
	EALREADY        Errno = 114
	EINPROGRESS     Errno = 115
	ESTALE          Errno = 116
	EUCLEAN         Errno = 117
	ENOTNAM         Errno = 118
	ENAVAIL         Errno = 119
	EISNAM          Errno = 120
	EREMOTEIO       Errno = 121
	EDQUOT          Errno = 122
	ENOMEDIUM       Errno = 123
	EMEDIUMTYPE     Errno = 124
	ECANCELED       Errno = 125
	ENOKEY          Errno = 126
	EKEYEXPIRED     Errno = 127
	EKEYREVOKED     Errno = 128
	EKEYREJECTED    Errno = 129
	EOWNERDEAD      Errno = 130
	ENOTRECOVERABLE Errno = 131
	ERFKILL         Errno = 132
	EHWPOISON       Errno = 133
)

//emgo:const
var errnos = []string{
	0:               "success",
	EPERM:           "operation not permitted",
	ENOENT:          "no such file or directory",
	ESRCH:           "no such process",
	EINTR:           "interrupted system call",
	EIO:             "I/O error",
	ENXIO:           "no such device or address",
	E2BIG:           "argument list too long",
	ENOEXEC:         "exec format error",
	EBADF:           "bad file number",
	ECHILD:          "no child processes",
	EAGAIN:          "try again",
	ENOMEM:          "out of memory",
	EACCES:          "permission denied",
	EFAULT:          "bad address",
	ENOTBLK:         "block device required",
	EBUSY:           "device or resource busy",
	EEXIST:          "file exists",
	EXDEV:           "cross-device link",
	ENODEV:          "no such device",
	ENOTDIR:         "not a directory",
	EISDIR:          "is a directory",
	EINVAL:          "invalid argument",
	ENFILE:          "file table overflow",
	EMFILE:          "too many open files",
	ENOTTY:          "not a typewriter",
	ETXTBSY:         "text file busy",
	EFBIG:           "file too large",
	ENOSPC:          "no space left on device",
	ESPIPE:          "illegal seek",
	EROFS:           "read-only file system",
	EMLINK:          "too many links",
	EPIPE:           "broken pipe",
	EDOM:            "math argument out of domain of func",
	ERANGE:          "math result not representable",
	EDEADLK:         "resource deadlock would occur",
	ENAMETOOLONG:    "file name too long",
	ENOLCK:          "no record locks available",
	ENOSYS:          "invalid system call number",
	ENOTEMPTY:       "directory not empty",
	ELOOP:           "too many symbolic links encountered",
	EWOULDBLOCK:     "deprecated. Use EAGAIN.",
	ENOMSG:          "no message of desired type",
	EIDRM:           "identifier removed",
	ECHRNG:          "channel number out of range",
	EL2NSYNC:        "level 2 not synchronized",
	EL3HLT:          "level 3 halted",
	EL3RST:          "level 3 reset",
	ELNRNG:          "link number out of range",
	EUNATCH:         "protocol driver not attached",
	ENOCSI:          "no CSI structure available",
	EL2HLT:          "level 2 halted",
	EBADE:           "invalid exchange",
	EBADR:           "invalid request descriptor",
	EXFULL:          "exchange full",
	ENOANO:          "no anode",
	EBADRQC:         "invalid request code",
	EBADSLT:         "invalid slot",
	EDEADLOCK:       "deprecated. Use EDEADLK",
	EBFONT:          "bad font file format",
	ENOSTR:          "device not a stream",
	ENODATA:         "no data available",
	ETIME:           "timer expired",
	ENOSR:           "out of streams resources",
	ENONET:          "machine is not on the network",
	ENOPKG:          "package not installed",
	EREMOTE:         "object is remote",
	ENOLINK:         "link has been severed",
	EADV:            "advertise error",
	ESRMNT:          "srmount error",
	ECOMM:           "communication error on send",
	EPROTO:          "protocol error",
	EMULTIHOP:       "multihop attempted",
	EDOTDOT:         "rFS specific error",
	EBADMSG:         "not a data message",
	EOVERFLOW:       "value too large for defined data type",
	ENOTUNIQ:        "name not unique on network",
	EBADFD:          "file descriptor in bad state",
	EREMCHG:         "remote address changed",
	ELIBACC:         "can not access a needed shared library",
	ELIBBAD:         "accessing a corrupted shared library",
	ELIBSCN:         ".lib section in a.out corrupted",
	ELIBMAX:         "attempting to link in too many shared libraries",
	ELIBEXEC:        "cannot exec a shared library directly",
	EILSEQ:          "illegal byte sequence",
	ERESTART:        "interrupted system call should be restarted",
	ESTRPIPE:        "streams pipe error",
	EUSERS:          "too many users",
	ENOTSOCK:        "socket operation on non-socket",
	EDESTADDRREQ:    "destination address required",
	EMSGSIZE:        "message too long",
	EPROTOTYPE:      "protocol wrong type for socket",
	ENOPROTOOPT:     "protocol not available",
	EPROTONOSUPPORT: "protocol not supported",
	ESOCKTNOSUPPORT: "socket type not supported",
	EOPNOTSUPP:      "operation not supported on transport endpoint",
	EPFNOSUPPORT:    "protocol family not supported",
	EAFNOSUPPORT:    "address family not supported by protocol",
	EADDRINUSE:      "address already in use",
	EADDRNOTAVAIL:   "cannot assign requested address",
	ENETDOWN:        "network is down",
	ENETUNREACH:     "network is unreachable",
	ENETRESET:       "network dropped connection because of reset",
	ECONNABORTED:    "software caused connection abort",
	ECONNRESET:      "connection reset by peer",
	ENOBUFS:         "no buffer space available",
	EISCONN:         "transport endpoint is already connected",
	ENOTCONN:        "transport endpoint is not connected",
	ESHUTDOWN:       "cannot send after transport endpoint shutdown",
	ETOOMANYREFS:    "too many references: cannot splice",
	ETIMEDOUT:       "connection timed out",
	ECONNREFUSED:    "connection refused",
	EHOSTDOWN:       "host is down",
	EHOSTUNREACH:    "no route to host",
	EALREADY:        "operation already in progress",
	EINPROGRESS:     "operation now in progress",
	ESTALE:          "stale file handle",
	EUCLEAN:         "structure needs cleaning",
	ENOTNAM:         "not a XENIX named type file",
	ENAVAIL:         "no XENIX semaphores available",
	EISNAM:          "is a named type file",
	EREMOTEIO:       "remote I/O error",
	EDQUOT:          "quota exceeded",
	ENOMEDIUM:       "no medium found",
	EMEDIUMTYPE:     "wrong medium type",
	ECANCELED:       "operation Canceled",
	ENOKEY:          "required key not available",
	EKEYEXPIRED:     "key has expired",
	EKEYREVOKED:     "key has been revoked",
	EKEYREJECTED:    "key was rejected by service",
	EOWNERDEAD:      "owner died",
	ENOTRECOVERABLE: "state not recoverable",
	ERFKILL:         "operation not possible due to RF-kill",
	EHWPOISON:       "memory page has hardware error",
}
