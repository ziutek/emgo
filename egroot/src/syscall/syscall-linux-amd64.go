// +build linux
// +build amd64

package syscall

const (
	sys_READ  = 0
	sys_WRITE = 1
	sys_OPEN  = 2
	sys_CLOSE = 3

	sys_MMAP = 9

	sys_BRK = 12

	sys_SOCKET  = 41
	sys_CONNECT = 42

	sys_BIND   = 49
	sys_LISTEN = 50

	sys_SETSOCKOPT = 54

	sys_CLONE = 56

	sys_EXIT = 60

	sys_UNLINK = 81
)

const (
	O_ACCMODE   = 0x3
	O_APPEND    = 0x400
	O_ASYNC     = 0x2000
	O_CLOEXEC   = 0x80000
	O_CREAT     = 0x40
	O_DIRECT    = 0x4000
	O_DIRECTORY = 0x10000
	O_DSYNC     = 0x1000
	O_EXCL      = 0x80
	O_FSYNC     = 0x101000
	O_LARGEFILE = 0x0
	O_NDELAY    = 0x800
	O_NOATIME   = 0x40000
	O_NOCTTY    = 0x100
	O_NOFOLLOW  = 0x20000
	O_NONBLOCK  = 0x800
	O_RDONLY    = 0x0
	O_RDWR      = 0x2
	O_RSYNC     = 0x101000
	O_SYNC      = 0x101000
	O_TRUNC     = 0x200
	O_WRONLY    = 0x1
)

const (
	PROT_EXEC      = 0x4
	PROT_GROWSDOWN = 0x1000000
	PROT_GROWSUP   = 0x2000000
	PROT_NONE      = 0x0
	PROT_READ      = 0x1
	PROT_WRITE     = 0x2

	MAP_32BIT      = 0x40
	MAP_ANON       = 0x20
	MAP_ANONYMOUS  = 0x20
	MAP_DENYWRITE  = 0x800
	MAP_EXECUTABLE = 0x1000
	MAP_FILE       = 0x0
	MAP_FIXED      = 0x10
	MAP_GROWSDOWN  = 0x100
	MAP_HUGETLB    = 0x40000
	MAP_LOCKED     = 0x2000
	MAP_NONBLOCK   = 0x10000
	MAP_NORESERVE  = 0x4000
	MAP_POPULATE   = 0x8000
	MAP_PRIVATE    = 0x2
	MAP_SHARED     = 0x1
	MAP_STACK      = 0x20000
	MAP_TYPE       = 0xf
)

const (
	AF_ALG        = 0x26
	AF_APPLETALK  = 0x5
	AF_ASH        = 0x12
	AF_ATMPVC     = 0x8
	AF_ATMSVC     = 0x14
	AF_AX25       = 0x3
	AF_BLUETOOTH  = 0x1f
	AF_BRIDGE     = 0x7
	AF_CAIF       = 0x25
	AF_CAN        = 0x1d
	AF_DECnet     = 0xc
	AF_ECONET     = 0x13
	AF_FILE       = 0x1
	AF_IEEE802154 = 0x24
	AF_INET       = 0x2
	AF_INET6      = 0xa
	AF_IPX        = 0x4
	AF_IRDA       = 0x17
	AF_ISDN       = 0x22
	AF_IUCV       = 0x20
	AF_KEY        = 0xf
	AF_LLC        = 0x1a
	AF_LOCAL      = 0x1
	AF_MAX        = 0x27
	AF_NETBEUI    = 0xd
	AF_NETLINK    = 0x10
	AF_NETROM     = 0x6
	AF_PACKET     = 0x11
	AF_PHONET     = 0x23
	AF_PPPOX      = 0x18
	AF_RDS        = 0x15
	AF_ROSE       = 0xb
	AF_ROUTE      = 0x10
	AF_RXRPC      = 0x21
	AF_SECURITY   = 0xe
	AF_SNA        = 0x16
	AF_TIPC       = 0x1e
	AF_UNIX       = 0x1
	AF_UNSPEC     = 0x0
	AF_WANPIPE    = 0x19
	AF_X25        = 0x9
)

const (
	SOCK_CLOEXEC   = 0x80000
	SOCK_DCCP      = 0x6
	SOCK_DGRAM     = 0x2
	SOCK_NONBLOCK  = 0x800
	SOCK_PACKET    = 0xa
	SOCK_RAW       = 0x3
	SOCK_RDM       = 0x4
	SOCK_SEQPACKET = 0x5
	SOCK_STREAM    = 0x1
)

const (
	SOL_AAL    = 0x109
	SOL_ATM    = 0x108
	SOL_DECNET = 0x105
	SOL_ICMPV6 = 0x3a
	SOL_IP     = 0x0
	SOL_IPV6   = 0x29
	SOL_IRDA   = 0x10a
	SOL_PACKET = 0x107
	SOL_RAW    = 0xff
	SOL_SOCKET = 0x1
	SOL_TCP    = 0x6
	SOL_X25    = 0x106
)

const (
	SO_ACCEPTCONN                    = 0x1e
	SO_ATTACH_FILTER                 = 0x1a
	SO_BINDTODEVICE                  = 0x19
	SO_BROADCAST                     = 0x6
	SO_BSDCOMPAT                     = 0xe
	SO_DEBUG                         = 0x1
	SO_DETACH_FILTER                 = 0x1b
	SO_DOMAIN                        = 0x27
	SO_DONTROUTE                     = 0x5
	SO_ERROR                         = 0x4
	SO_KEEPALIVE                     = 0x9
	SO_LINGER                        = 0xd
	SO_MARK                          = 0x24
	SO_NO_CHECK                      = 0xb
	SO_OOBINLINE                     = 0xa
	SO_PASSCRED                      = 0x10
	SO_PASSSEC                       = 0x22
	SO_PEERCRED                      = 0x11
	SO_PEERNAME                      = 0x1c
	SO_PEERSEC                       = 0x1f
	SO_PRIORITY                      = 0xc
	SO_PROTOCOL                      = 0x26
	SO_RCVBUF                        = 0x8
	SO_RCVBUFFORCE                   = 0x21
	SO_RCVLOWAT                      = 0x12
	SO_RCVTIMEO                      = 0x14
	SO_REUSEADDR                     = 0x2
	SO_RXQ_OVFL                      = 0x28
	SO_SECURITY_AUTHENTICATION       = 0x16
	SO_SECURITY_ENCRYPTION_NETWORK   = 0x18
	SO_SECURITY_ENCRYPTION_TRANSPORT = 0x17
	SO_SNDBUF                        = 0x7
	SO_SNDBUFFORCE                   = 0x20
	SO_SNDLOWAT                      = 0x13
	SO_SNDTIMEO                      = 0x15
	SO_TIMESTAMP                     = 0x1d
	SO_TIMESTAMPING                  = 0x25
	SO_TIMESTAMPNS                   = 0x23
	SO_TYPE                          = 0x3
)
