// +build linux

package os

import (
	"internal"
	"unsafe"
)

func strlen(s *[1<<31 - 1]byte) int {
	for n, c := range s {
		if c == 0 {
			return n
		}
	}
	panic("strlen overflow")
}

func args(begin, end uintptr) []string {
	argv := (*[1<<31 - 1]*[1<<31 - 1]byte)(unsafe.Pointer(begin))
	args := make([]string, (end-begin)/unsafe.Sizeof(uintptr(0))-1)
	for i := range args {
		a := argv[i]
		s := a[:strlen(a)]
		args[i] = *(*string)(unsafe.Pointer(&s))
	}
	return args
}

var (
	Args []string
	Env  []string
)

func init() {
	Args = args(internal.Argv, internal.Env)
	Env = args(internal.Env, internal.Auxv)
}

func Getenv(key string) string {
	for _, s := range Env {
		if len(s) <= len(key) {
			continue
		}
		if s[len(key)] == '=' && s[:len(key)] == key {
			return s[len(key)+1:]
		}
	}
	return ""
}
