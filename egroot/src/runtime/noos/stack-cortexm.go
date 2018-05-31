// +build cortexm0 cortexm3 cortexm4 cortexm4f cortexm7f cortexm7d

package noos

import (
	"unsafe"

	"arch/cortexm"
	"arch/cortexm/mpu"
)

//c:inline
func stacksBegin() uintptr

//c:inline
func isrStackSize() uintptr

//c:inline
func mainStackSize() uintptr

//c:inline
func taskStackSize() uintptr

//c:inline
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
// Since taskInfos was placed in stack guard slots, the less restricted
// mpu.Arw__ was used. TODO: Check is there a cheap way to return to mpu.A____.
func setMPUStackGuard(n int) {
	mpu.SetRegion(
		stackGuardBegin(n)+mpu.VALID+mpu5_StackGuard,
		mpu.ENA|mpu.SIZE(5)|mpu.C|mpu.S|mpu.Arw__|mpu.XN,
	)
	cortexm.DSB()
	// Return from exception works like ISB so can ommit it here.
}

func stackGuardArray(n int) *[8]uint32 {
	return (*[8]uint32)(unsafe.Pointer(stackGuardBegin(n)))
}

const stackGuardMagic = 0xffffffee

func resetStackGuard(n int) {
	sg := stackGuardArray(n)
	const m = int(unsafe.Sizeof(taskInfo{})+3) / 4
	for i := 0; i < m; i++ {
		sg[i] = 0
	}
	for i := m; i < len(sg); i++ {
		sg[i] = stackGuardMagic
	}
}

func checkStackGuard(n int) {
	sg := stackGuardArray(n)
	for i := int(unsafe.Sizeof(taskInfo{})+3) / 4; i < len(sg); i++ {
		if sg[i] != stackGuardMagic {
			panic("nnos: stack guard")
		}
	}
}
