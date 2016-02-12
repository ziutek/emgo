// +build cortexm0 cortexm3 cortexm4 cortexm4f

package noos

import (
	"arch/cortexm"
	"unsafe"
)

func stackTaskLog2() uint

func stackTaskFrac() uint

func stackEnd() uintptr

func stackTop(i int) uintptr {
	stackSize := uintptr((1 << stackTaskLog2()) * stackTaskFrac() / 8)
	return stackEnd() - uintptr(i)*stackSize
}

func allocStackFrame(sp uintptr) (*cortexm.StackFrame, uintptr) {
	sp -= unsafe.Sizeof(cortexm.StackFrame{})
	return (*cortexm.StackFrame)(unsafe.Pointer(sp)), sp
}
