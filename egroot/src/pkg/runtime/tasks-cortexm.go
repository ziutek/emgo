// +build cortexm0 cortexm3 cortexm4 cortexm4f

package runtime

import (
	"unsafe"

	"cortexm"
	"sync/barrier"
)

func MaxTasks() int

type task struct {
	sp    uintptr
	flags uint8
	prio  uint8
}

var (
	tasks       []task
	stacksStart uintptr
	stackCap    uintptr
	heapStart   uintptr
)

const pspLen = 1024

func init() {
	freeStart := freeStart()
	freeEnd := freeEnd()

	setSlice(
		unsafe.Pointer(&tasks), 
		unsafe.Pointer(freeStart),
		uint(MaxTasks()), uint(MaxTasks()),
	)
	
	heapStart = freeStart + uintptr(len(tasks)) * unsafe.Sizeof(task{})
	stacksStart = heapStart + HeapSize()
	if stacksStart&3 != 0 {
		// align stacks to next 4 bytes
		stacksStart = (stacksStart + 4) &^ 3
	}

	if stacksStart > freeEnd {
		panicMemory()
	}

	stackCap = (freeEnd - stacksStart) / uintptr(len(tasks)+1)
	stackCap &^= 3

	for i := 0; i < len(tasks); i++ {
		tasks[i] = task{sp: stacksStart + uintptr(i+1)*stackCap}
	}

	msp := cortexm.MSP()
	mspLen := freeEnd - msp

	psp := tasks[0].sp - mspLen
	tasks[0].sp = psp
	
	Copy(unsafe.Pointer(psp), unsafe.Pointer(msp), uint(mspLen))
	cortexm.SetPSP(unsafe.Pointer(psp))
	barrier.Sync()
	cortexm.SetCtrl(cortexm.Ctrl() | cortexm.UsePSP)
}
