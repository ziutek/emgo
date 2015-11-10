// +build cortexm0 cortexm3 cortexm4 cortexm4f

package noos

import (
	"arch/cortexm/exce"
	"unsafe"
)

func stackExp() uint

func stackFrac() uint

func stackEnd() uintptr

var stackCap = uintptr((1 << stackExp()) * stackFrac() / 8)

func stackTop(i int) uintptr {
	return stackEnd() - uintptr(i)*stackCap
}

func allocStackFrame(sp uintptr) (*exce.StackFrame, uintptr) {
	sp -= unsafe.Sizeof(exce.StackFrame{})
	return (*exce.StackFrame)(unsafe.Pointer(sp)), sp
}