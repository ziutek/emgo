// +build cortexm0 cortexm3 cortexm4 cortexm4f

package noos

import (
	"arch/cortexm"
	"unsafe"
)

func stackExp() uint

func stackFrac() uint

func stackEnd() uintptr

var stackCap = uintptr((1 << stackExp()) * stackFrac() / 8)

func stackTop(i int) uintptr {
	return stackEnd() - uintptr(i)*stackCap
}

func allocStackFrame(sp uintptr) (*cortexm.StackFrame, uintptr) {
	sp -= unsafe.Sizeof(cortexm.StackFrame{})
	return (*cortexm.StackFrame)(unsafe.Pointer(sp)), sp
}
