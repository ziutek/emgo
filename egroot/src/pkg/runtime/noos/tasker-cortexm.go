// +build cortexm3 cortexm4 cortexm4f

package noos

import (
	"unsafe"

	"sync/barrier"

	"cortexm"
	"cortexm/irq"
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
	next  uint16
	state taskState
	prio  uint8
}

type taskSched struct {
	tasks     []taskInfo
	curTask   int
	forceNext int
	onSysTick bool
}

var tasker taskSched

func (ts *taskSched) run() {
	ts.onSysTick = true
}

func (ts *taskSched) stop() {
	ts.onSysTick = false
}

func (ts *taskSched) init() {
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
		unsafe.Pointer(&ts.tasks), Heap,
		MaxTasks(), unsafe.Sizeof(taskInfo{}), unsafe.Alignof(taskInfo{}),
		unsafe.Alignof(taskInfo{}),
	)
	if Heap == nil {
		panicMemory()
	}

	ts.tasks[0] = taskInfo{prio: 255}
	ts.tasks[0].state.SetReady()
	for i := 1; i < len(ts.tasks); i++ {
		ts.tasks[i].state.SetEmpty()
	}

	// Use PSP as stack pointer for thread mode.
	cortexm.SetPSP(unsafe.Pointer(cortexm.MSP()))
	barrier.Sync()
	cortexm.SetCtrl(cortexm.Ctrl() | cortexm.UsePSP)
	cortexm.ISB()

	// Now MSP is used only by exceptions handlers.
	cortexm.SetMSP(unsafe.Pointer(initSP(len(ts.tasks))))

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

	tasker.forceNext = -1
	tasker.run()
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

var cfsr = (*cfs)(unsafe.Pointer(uintptr(0xe000ed28)))

func hardFaultHandler() {
	for {
	}
}

func memFaultHandler() {
	mmfs := cfsr.mmfs
	_ = mmfs
	for {
	}
}

func busFaultHandler() {
	bfs := cfsr.bfs
	_ = bfs
	for {
	}
}

func usageFaultHandler() {
	ufs := cfsr.ufs
	pfp := (*stackFrame)(unsafe.Pointer(cortexm.PSP()))
	_, _ = ufs, pfp
	for {
	}
}

var Tick uint32

func sysTickHandler() {
	Tick++
	if tasker.onSysTick {
		irq.PendSV.SetPending()
	}
}

// pendSVHandler calls nextTask with PSP for current task. It performs
// context swich if nextTask returns new non-zero value for PSP.
func pendSVHandler()

// nextTask returns taskInfo.sp for next task or 0. pendSVHandler can use two
// least significant bits of sp as flags so sp isn't real stack pointer.
// TODO: better scheduler
func nextTask(sp uintptr) uintptr {
	n := tasker.forceNext
	if n >= 0 {
		tasker.forceNext = -1
	} else {
		n = tasker.curTask
		for {
			if n++; n >= len(tasker.tasks) {
				n = 0
			}
			if tasker.tasks[n].state.Ready() {
				break
			}
			if n == tasker.curTask {
				panic("no task to run")
			}
		}
		if n == tasker.curTask {
			return 0
		}
	}
	tasker.tasks[tasker.curTask].sp = sp
	tasker.curTask = n
	return tasker.tasks[n].sp
}

func (ts *taskSched) newTask(pc uintptr, xpsr uint32, wait bool) {
	n := ts.curTask
	for {
		if n++; n >= len(ts.tasks) {
			n = 0
		}
		if ts.tasks[n].state.Empty() {
			break
		}
		if n == ts.curTask {
			panic("too many tasks")
		}
	}

	sf, sp := allocStackFrame(initSP(n))
	ts.tasks[n] = taskInfo{sp: sp, prio: 255} // (re)initialization

	// Use parent's xPSR as initial xPSR for new task.
	sf.xpsr = xpsr
	sf.pc = pc

	ts.tasks[n].state.SetReady()

	if wait {
		ts.stop()
		// This badly affects scheduling but don't care for now.
		ts.forceNext = n
		irq.PendSV.SetPending()
	}
}

func (ts *taskSched) delTask(n int) {
	ts.tasks[n].state.SetEmpty()
	irq.PendSV.SetPending()
}
