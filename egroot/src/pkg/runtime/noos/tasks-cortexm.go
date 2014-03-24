// +build cortexm0 cortexm3 cortexm4 cortexm4f

package noos

import (
	"unsafe"

	"sync/barrier"

	"cortexm"
	"cortexm/irq"
	"cortexm/sleep"
	"cortexm/systick"
)

var stackCap int

func initSP(i int) uintptr {
	return heapStackEnd() - uintptr(i*stackCap)
}

type taskFlags byte

const (
	taskEmpty taskFlags = iota
	taskReady
	taskSleep
)

type taskCtx struct {
	sp    uintptr
	flags taskFlags
	prio  uint8
}

var (
	tasks   []taskCtx
	actTask int
)

func init() {
	stackExp := (stackSize()>>8)&0xff + 1
	stackNum := int(stackSize()&0xff)>>5 + 1
	stackCap = (1 << stackExp) * stackNum / 8

	heapStack := heapStack()

	if MaxTasks() == 0 {
		if stackCap > len(heapStack) {
			panicMemory()
		}
		stackEnd := len(heapStack) - stackCap
		heap = heapStack[:stackEnd:stackEnd]
		return
	}

	stacksEnd := len(heapStack) - (MaxTasks()+1)*stackCap

	var vt []irq.Vector
	vtlen := 1 << maxIRQ()
	vtsize := vtlen * int(unsafe.Sizeof(irq.Vector{}))

	heap = allocTop(
		unsafe.Pointer(&vt), heapStack[:stacksEnd],
		vtlen, unsafe.Sizeof(irq.Vector{}), uintptr(vtsize),
	)
	if heap == nil {
		panicMemory()
	}

	// There can be some free memory after vector table.
	irqHeap := heapStack[len(heap)+vtsize : stacksEnd]

	if h := allocBottom(
		unsafe.Pointer(&tasks), irqHeap,
		MaxTasks(), unsafe.Sizeof(taskCtx{}), unsafe.Alignof(taskCtx{}),
	); h != nil {
		irqHeap = h
	} else {
		heap = allocTop(
			unsafe.Pointer(&tasks), heap,
			MaxTasks(), unsafe.Sizeof(taskCtx{}), unsafe.Alignof(taskCtx{}),
		)
		if heap == nil {
			panicMemory()
		}
	}

	tasks[0].flags = taskReady
	for i := 1; i < len(tasks); i++ {
		tasks[i] = taskCtx{flags: 0}
	}

	// Use PSP as stack pointer for thread mode.
	cortexm.SetPSP(unsafe.Pointer(cortexm.MSP()))
	barrier.Sync()
	cortexm.SetCtrl(cortexm.Ctrl() | cortexm.UsePSP)
	cortexm.ISB()

	// Now MSP is used only by exceptions handlers.
	cortexm.SetMSP(unsafe.Pointer(initSP(len(tasks))))

	// Setup interrupt table
	vt[irq.Reset] = irq.VectorFor(resetHandler)
	vt[irq.NMI] = irq.VectorFor(nmiHandler)
	vt[irq.HardFault] = irq.VectorFor(hardFaultHandler)
	vt[irq.MemFault] = irq.VectorFor(memFaultHandler)
	vt[irq.BusFault] = irq.VectorFor(busFaultHandler)
	vt[irq.UsageFault] = irq.VectorFor(usageFaultHandler)
	vt[irq.PendSV] = irq.VectorFor(pendSVHandler)
	vt[irq.SysTick] = irq.VectorFor(sysTickHandler)
	irq.UseTable(vt)

	irq.MemFault.Enable()
	irq.BusFault.Enable()
	irq.UsageFault.Enable()

	irq.PendSV.SetPrio(irq.Lowest)

	systick.SetReload(1e6 - 1)
	systick.SetFlags(systick.Enable | systick.TickInt | systick.ClkCPU)
}

func resetHandler() {
	for {
	}
}

func nmiHandler() {
	for {
	}
}

type cfs struct {
	mmfs uint8  `C:"volatile"`
	bfs  uint8  `C:"volatile"`
	ufs  uint16 `C:"volatile"`
}

var cfsr = (*cfs)(unsafe.Pointer(uintptr(0xE000ED28)))

func hardFaultHandler() {
	for {
	}
}

func memFaultHandler() {
	// Check cfsr.mmfs.
	for {
	}
}

func busFaultHandler() {
	// Check cfsr.bfs.
	for {
	}
}

func usageFaultHandler() {
	// Check cfsr.ufs.
	for {
	}
}

func nextTask(sp uintptr) uintptr {
	n := actTask
	for {
		if n++; n >= len(tasks) {
			n = 0
		}
		if tasks[n].flags&3 == taskReady {
			break
		}
		if n == actTask {
			sleep.WFI()
		}
	}
	if n == actTask {
		return 0
	}
	tasks[actTask].sp = sp
	actTask = n
	return tasks[n].sp
}

// pendSVHandler calls nextTask with PSP for current task. It does context
// swich if nextTask returns new non-zero value for PSP.
func pendSVHandler()

var Tick uint32

func sysTickHandler() {
	Tick++
	irq.PendSV.SetPending()
}
