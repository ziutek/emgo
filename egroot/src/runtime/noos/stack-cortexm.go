// +build cortexm0 cortexm3 cortexm4 cortexm4f

package noos

import (
	"arch/cortexm"
	"unsafe"
)

func stackLog2() uint

func stackFrac() uint

func stackEnd() uintptr

func stackTop(i int) uintptr {
	stackSize := uintptr((1 << stackLog2()) * stackFrac() / 8)
	return stackEnd() - uintptr(i)*stackSize
}

func allocStackFrame(sp uintptr) (*cortexm.StackFrame, uintptr) {
	sp -= unsafe.Sizeof(cortexm.StackFrame{})
	return (*cortexm.StackFrame)(unsafe.Pointer(sp)), sp
}
