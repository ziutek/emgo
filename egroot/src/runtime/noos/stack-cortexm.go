// +build cortexm0 cortexm3 cortexm4 cortexm4f

package noos

import (
	"unsafe"

	"arch/cortexm"
	"arch/cortexm/mpu"
)

func stacksBegin() uintptr
func isrStackSize() uintptr
func mainStackSize() uintptr
func taskStackSize() uintptr
func stacksEnd() uintptr

func stackTop(n int) uintptr {
	if n == 0 {
		return stacksEnd()
	}
	return stacksEnd() - mainStackSize() - uintptr(n-1)*taskStackSize()
}

func allocStackFrame(sp uintptr) (*cortexm.StackFrame, uintptr) {
	sp -= unsafe.Sizeof(cortexm.StackFrame{})
	return (*cortexm.StackFrame)(unsafe.Pointer(sp)), sp
}

// setStackGuard uses MPU to setup 32 byte stack guard area at the bottom of the
// stack of n-th task.
func setStackGuard(n int) {
	mpu.SetRegion(mpu.BaseAttr{
		stacksEnd() - mainStackSize() - uintptr(n)*taskStackSize() +
			mpu.VALID + regionStackGuard,
		mpu.ENA | mpu.SIZE(5) | mpu.C | mpu.S | mpu.A____ | mpu.XN,
	})
}
