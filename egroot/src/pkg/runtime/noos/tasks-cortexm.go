// +build cortexm0 cortexm3 cortexm4 cortexm4f

package noos

import (
	"unsafe"

	"cortexm"
	"sync/barrier"
)

type taskCtx struct {
	sp    uintptr
	flags uint8
	prio  uint8
}

var tasks []taskCtx


func init() {
	heapStack := heapStack()
	if heapSize() + minStack > uintptr(len(heapStack)) {
		panicMemory()
	}
	heap = heapStack[:heapSize()]
	
	if MaxTasks() <= 0 {
		return
	}

	heap = alloc(
		unsafe.Pointer(&tasks), heap, 
		MaxTasks(), unsafe.Sizeof(taskCtx{}),
	)

	stacksSetup(heapStack[heapSize():])
	
	tasks[0].sp = cortexm.MSP()
	for i := 1; i < len(tasks); i++ {
		tasks[i] = taskCtx{sp: stackInitSP(i)}
	}
	
	// Use PSP as stack pointer for thread mode.
	cortexm.SetPSP(unsafe.Pointer(tasks[0].sp))
	barrier.Sync()
	cortexm.SetCtrl(cortexm.Ctrl() | cortexm.UsePSP)
	cortexm.ISB()

	// Now MSP is used only by exceptions handlers.
	cortexm.SetMSP(unsafe.Pointer(stackInitSP(len(tasks))))
}
