// +build cortexm3 cortexm4 cortexm4f

package noos

import (
	"unsafe"

	"sync/barrier"

	"cortexm"
	"cortexm/irq"
	"cortexm/sleep"
	"cortexm/systick"
)

func evtExp() uint

var stackCap = uintptr((1 << stackExp()) * stackFrac() / 8)

func initSP(i int) uintptr {
	return stackEnd() - uintptr(i)*stackCap
}

type taskState byte

const (
	taskEmpty taskState = iota
	taskReady
)

func (s taskState) Ready() bool {
	return s&3 == taskReady
}

func (s *taskState) SetReady() {
	*s = *s&^3 | taskReady
}

func (s taskState) Empty() bool {
	return s&3 == taskEmpty
}

func (s *taskState) SetEmpty() {
	*s = *s&^3 | taskEmpty
}

// taskInfo
// sp contains value of SP after automatic stacking during exception entry. So
// sp points to the last register in set automatically stacked by CPU and just
// after the register set stacked by tasker. pendSVHandler can use two least
// significant bits of sp as flags.
type taskInfo struct {
	sp    uintptr
	state taskState
	prio  uint8
}

var (
	tasks   []taskInfo
	curTask int
)

func initTasker() {
	var vt []irq.Vector
	vtlen := 1 << evtExp()
	vtsize := vtlen * int(unsafe.Sizeof(irq.Vector{}))

	Heap = allocTop(
		unsafe.Pointer(&vt), Heap,
		vtlen, unsafe.Sizeof(irq.Vector{}), unsafe.Alignof(irq.Vector{}),
		uintptr(vtsize),
	)
	if Heap == nil {
		panicMemory()
	}

	Heap = allocTop(
		unsafe.Pointer(&tasks), Heap,
		MaxTasks(), unsafe.Sizeof(taskInfo{}), unsafe.Alignof(taskInfo{}),
		unsafe.Alignof(taskInfo{}),
	)
	if Heap == nil {
		panicMemory()
	}

	tasks[0] = taskInfo{prio: 255}
	tasks[0].state.SetReady()
	for i := 1; i < len(tasks); i++ {
		tasks[i].state.SetEmpty()
	}

	// Use PSP as stack pointer for thread mode.
	cortexm.SetPSP(unsafe.Pointer(cortexm.MSP()))
	barrier.Sync()
	cortexm.SetCtrl(cortexm.Ctrl() | cortexm.UsePSP)
	cortexm.ISB()

	// Now MSP is used only by exceptions handlers.
	cortexm.SetMSP(unsafe.Pointer(initSP(len(tasks))))

	// Setup interrupt table.
	// Consider setup at link time using GCC weak functions to support Cortex-M0
	// and (in case of Cortex-M3,4) to allow vector load on the ICode bus
	// simultaneously with registers stacking on DCode bus.
	vt[irq.Reset] = irq.VectorFor(resetHandler)
	vt[irq.NMI] = irq.VectorFor(nmiHandler)
	vt[irq.HardFault] = irq.VectorFor(hardFaultHandler)
	vt[irq.MemFault] = irq.VectorFor(memFaultHandler)
	vt[irq.BusFault] = irq.VectorFor(busFaultHandler)
	vt[irq.UsageFault] = irq.VectorFor(usageFaultHandler)
	vt[irq.SVCall] = irq.VectorFor(svcHandler)
	vt[irq.PendSV] = irq.VectorFor(pendSVHandler)
	vt[irq.SysTick] = irq.VectorFor(sysTickHandler)
	irq.UseTable(vt)

	irq.MemFault.Enable()
	irq.BusFault.Enable()
	irq.UsageFault.Enable()

	irq.SVCall.SetPrio(irq.Lowest)
	irq.PendSV.SetPrio(irq.Lowest)

	// One context switch per 5e5 SysTicks (140/s for 70 Mhz, 336/s for 168 MHz)
	systick.SetReload(5e5 - 1)
	systick.WriteFlags(systick.Enable | systick.TickInt | systick.ClkCPU)
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

var Tick uint32

func sysTickHandler() {
	Tick++
	irq.PendSV.SetPending()
}

// pendSVHandler calls nextTask with PSP for current task. It performs
// context swich if nextTask returns new non-zero value for PSP.
func pendSVHandler()

// nextTask returns taskInfo.sp for nextTask or 0. pendSVHandler can use two
// least significant bits of sp as flags.
func nextTask(sp uintptr) uintptr {
	n := curTask
	for {
		if n++; n >= len(tasks) {
			n = 0
		}
		if tasks[n].state.Ready() {
			break
		}
		if n == curTask {
			sleep.WFE()
		}
	}
	if n == curTask {
		return 0
	}
	tasks[curTask].sp = sp
	curTask = n
	barrier.Memory()

	return tasks[n].sp
}

func newTask(pc uintptr, xpsr uint32) {
	n := curTask
	for {
		if n++; n >= len(tasks) {
			n = 0
		}
		if tasks[n].state.Empty() {
			break
		}
		if n == curTask {
			panic("too many tasks")
		}
	}

	sf, sp := allocStackFrame(initSP(n))
	tasks[n] = taskInfo{sp: sp, prio: 255} // (re)initialization

	// Use parent's xPSR as initial xPSR for new task.
	sf.xpsr = xpsr
	sf.pc = pc

	tasks[n].state.SetReady()
	barrier.Memory()
}

func delTask() {
	tasks[curTask].state.SetEmpty()
	barrier.Memory()
	irq.PendSV.SetPending()
}
