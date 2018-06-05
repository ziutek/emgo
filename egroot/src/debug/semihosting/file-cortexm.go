// +builg cortexm0 cortexm0p cortexm3 cortexm4 cortexm4f cortexm7f cortexm7d

package semihosting

import (
	"reflect"
	"unsafe"
)

//emgo:export
//c:inline
func hostIO(cmd int, p unsafe.Pointer) int

func hostErrno() Error {
	return Error(hostIO(0x13, nil))
}

func openFile(name string, mode FopenMode) (File, error) {
	type args struct {
		path    uintptr
		mode    int
		pathLen int
	}
	p := &args{
		(*reflect.StringHeader)(unsafe.Pointer(&name)).Data,
		int(mode),
		len(name),
	}
	ret := hostIO(0x01, unsafe.Pointer(p))
	if ret == -1 {
		return File{}, hostErrno()
	}
	return File{ret}, nil
}

func (f File) close() error {
	if hostIO(0x02, unsafe.Pointer(&f.fd)) == -1 {
		return hostErrno()
	}
	return nil
}

func (f File) writeString(s string) (int, error) {
	type args struct {
		fd      int
		data    uintptr
		dataLen int
	}
	p := &args{
		f.fd,
		(*reflect.StringHeader)(unsafe.Pointer(&s)).Data,
		len(s),
	}
	ret := hostIO(0x05, unsafe.Pointer(p))
	if ret != 0 {
		return len(s) - ret, hostErrno()
	}
	return len(s), nil
}

func (f File) write(b []byte) (int, error) {
	return f.writeString(*(*string)(unsafe.Pointer(&b)))
}
