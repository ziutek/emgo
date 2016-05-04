// +build linux

package syscall

import (
	"internal"
	"unsafe"
)

const minerr = ^uintptr(4095) + 1

func syscallPath(trap uintptr, path string, a2, a3 uintptr) uintptr

func Read(fd int, b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}
	ret := internal.Syscall3(
		sys_READ,
		uintptr(fd), uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)),
	)
	if ret >= minerr {
		return 0, -Errno(ret)
	}
	return int(ret), nil
}

func WriteString(fd int, s string) (int, error) {
	if len(s) == 0 {
		return 0, nil
	}
	p := (*internal.String)(unsafe.Pointer(&s))
	ret := internal.Syscall3(sys_WRITE, uintptr(fd), p.Addr, p.Len)
	if ret >= minerr {
		return 0, -Errno(ret)
	}
	return int(ret), nil
}

func Write(fd int, b []byte) (int, error) {
	return WriteString(fd, *(*string)(unsafe.Pointer(&b)))
}

func Open(path string, mode int, perm uint32) (int, error) {
	ret := syscallPath(sys_OPEN, path, uintptr(mode), uintptr(perm))
	if ret >= minerr {
		return 0, -Errno(ret)
	}
	return int(ret), nil
}

func Close(fd int) error {
	ret := internal.Syscall1(sys_CLOSE, uintptr(fd))
	if ret >= minerr {
		return -Errno(ret)
	}
	return nil
}

func Mmap(addr, length uintptr, prot, flags, fd, offset4k int) (uintptr, error) {
	ret := internal.Syscall6(
		sys_MMAP,
		addr, length, uintptr(prot), uintptr(flags), uintptr(fd), uintptr(offset4k),
	)
	if ret >= minerr {
		return 0, -Errno(ret)
	}
	return ret, nil
}

func Brk(brk unsafe.Pointer) uintptr {
	return internal.Syscall1(sys_BRK, uintptr(brk))
}

func Socket(domain, typ, proto int) (int, error) {
	ret := internal.Syscall3(
		sys_SOCKET, uintptr(domain), uintptr(typ), uintptr(proto),
	)
	if ret >= minerr {
		return 0, -Errno(ret)
	}
	return int(ret), nil
}

type Sockaddr interface {
	sockaddr() (ptr, size uintptr, err error)
}

type RawSockaddrInet4 struct {
	Family uint16
	Port   uint16
	Addr   [4]byte
	Zero   [8]uint8
}

//emgo:minfo
func (sa *RawSockaddrInet4) sockaddr() (ptr, size uintptr, err error) {
	return uintptr(unsafe.Pointer(sa)), unsafe.Sizeof(*sa), nil
}

func Bind(fd int, sa Sockaddr) error {
	ptr, size, err := sa.sockaddr()
	if err != nil {
		return err
	}
	ret := internal.Syscall3(sys_BIND, uintptr(fd), ptr, size)
	if ret >= minerr {
		return -Errno(ret)
	}
	return nil
}

func SetsockoptInt(fd, level, opt, value int) error {
	ret := internal.Syscall5(
		sys_SETSOCKOPT, uintptr(fd), uintptr(level), uintptr(opt),
		uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(value),
	)
	if ret >= minerr {
		return -Errno(ret)
	}
	return nil
}

func Exit(code int) {
	internal.Syscall1(sys_EXIT, uintptr(code))
}

type Timespec struct {
	Sec  int64
	Nsec int64
}

func ClockGettime(clkid int, tp *Timespec) error {
	ret := clock_gettime(clkid, tp)
	if ret >= minerr {
		return -Errno(ret)
	}
	return nil
}
