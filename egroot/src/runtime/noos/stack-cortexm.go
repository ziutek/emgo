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

func stackGuardBegin(n int) uintptr {
	return stacksEnd() - mainStackSize() - uintptr(n)*taskStackSize()
}

// setMPUStackGuard uses MPU to setup 32 byte (8 words) stack guard area at the
// bottom of the stack of n-th task.
//
// When taskInfo's are placed in stack guard slots, the less restricted
// mpu.Arw__ is used. TODO: Check is there a cheap way to return to mpu.A____.
func setMPUStackGuard(n int) {
	mpu.SetRegion(mpu.BaseAttr{
		stackGuardBegin(n) + mpu.VALID + regionStackGuard,
		mpu.ENA | mpu.SIZE(5) | mpu.C | mpu.S | mpu.Arw__ | mpu.XN,
	})
}

func stackGuardArray(n int) *[8]uint32 {
	return (*[8]uint32)(unsafe.Pointer(stackGuardBegin(n)))
}

const stackGuardMagic = 0xffffffee

func resetStackGuard(n int) {
	sg := stackGuardArray(n)
	const m = len(sg) - 2
	for i := 0; i < m; i++ {
		sg[i] = 0
	}
	sg[m] = stackGuardMagic
	sg[m+1] = stackGuardMagic
}

func checkStackGuard(n int) {
	sg := stackGuardArray(n)
	const m = len(sg) - 2
	for i := int(unsafe.Sizeof(taskInfo{})+3) / 4; i < m; i++ {
		if sg[i] != 0 {
			goto violation
		}
	}
	if sg[m] != stackGuardMagic || sg[m+1] != stackGuardMagic {
		goto violation
	}
	return
violation:
	panic("nnos: stack guard")
}
