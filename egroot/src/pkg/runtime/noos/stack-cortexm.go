// +build cortexm0 cortexm3 cortexm4 cortexm4f

package noos

import "unsafe"

func stackExp() uint

func stackFrac() uint

func stackEnd() uintptr

var stackCap = uintptr((1 << stackExp()) * stackFrac() / 8)

func initSP(i int) uintptr {
	return stackEnd() - uintptr(i)*stackCap
}

type stackFrame struct {
	r    [4]uintptr
	ip   uintptr
	lr   uintptr
	pc   uintptr
	xpsr uint32
}

func allocStackFrame(sp uintptr) (*stackFrame, uintptr) {
	sp -= unsafe.Sizeof(stackFrame{})
	return (*stackFrame)(unsafe.Pointer(sp)), sp
}

