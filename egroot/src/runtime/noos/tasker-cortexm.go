// +build cortexm0 cortexm3 cortexm4 cortexm4f

package noos

import (
	"arch/cortexm"
	"arch/cortexm/nvic"
	"arch/cortexm/scb"
	"math/rand"
	"sync/fence"
	"syscall"
	"unsafe"
)

type taskState byte

const (
	taskEmpty taskState = iota
	taskReady
	taskLocked
	taskWaitEvent
)

// taskInfo
// sp contains value of SP after automatic stacking during exception entry. So
// sp points to the last register in set automatically stacked by CPU and just
// after the register set stacked by tasker. pendSVHandler can use two least
// significant bits of sp for its flags.
type taskInfo struct {
	rng    rand.XorShift64
	sp     uintptr
	event  syscall.Event
	parent int16
	flags  taskState
	prio   uint8
}

func (ti *taskInfo) Init(parent int) {
	*ti = taskInfo{parent: int16(parent), prio: 255}
	if uptime == nil {
		ti.rng.Seed(uint64(uintptr(unsafe.Pointer(ti))))
	} else {
		ti.rng.Seed(uptime() + 1)
	}
}

func (ti *taskInfo) State() taskState {
	return taskState(ti.flags & 3)
}

func (ti *taskInfo) SetState(s taskState) {
	ti.flags = ti.flags&^3 | s
}

type taskSched struct {
	tasks   []taskInfo
	curTask int
}

var tasker taskSched

func (ts *taskSched) deliverEvent(e syscall.Event) {
	for i := range ts.tasks {
		t := &ts.tasks[i]
		switch t.State() {
		case taskEmpty:
			// skip

		case taskWaitEvent:
			if t.event&e != 0 {
				t.event = 0
				t.SetState(taskReady)
			}

		default:
			t.event |= e
		}
	}
}

/*
func setupVectorTable(vtLenExp int) {
	// vt schould be allocated before anything other (first allocation in
	// program run) to satisfy NVIC allignment restrictions.
	vt := make([]exce.Vector, 1<<vtLenExp)

	// Setup interrupt table.
	// Consider setup at link time using GCC weak functions to support
	// Cortex-M0 and (in case of Cortex-M3,4) to allow vector load on ICode bus
	// simultaneously with registers stacking on DCode bus.
	vt[exce.NMI] = exce.VectorFor(nmiHandler)
	vt[exce.HardFault] = exce.VectorFor(FaultHandler)
	vt[exce.MemManage] = exce.VectorFor(FaultHandler)
	vt[exce.BusFault] = exce.VectorFor(FaultHandler)
	vt[exce.UsageFault] = exce.VectorFor(FaultHandler)
	vt[exce.SVC] = exce.VectorFor(svcHandler)
	vt[exce.PendSV] = exce.VectorFor(pendSVHandler)
	vt[exce.SysTick] = exce.VectorFor(sysTickHandler)
	exce.UseTable(vt)
}
*/

func (ts *taskSched) init() {
	//setupVectorTable(irtExp) - disabled (we use static VT, set at link time)

	ts.tasks = make([]taskInfo, MaxTasks())

	// Use PSP as stack pointer for thread mode. Current (zero) task has stack
	// at top of the stacks area.
	cortexm.SetPSP(unsafe.Pointer(cortexm.MSP()))
	fence.Sync()
	cortexm.SetCONTROL(cortexm.CONTROL() | cortexm.UsePSP)
	cortexm.ISB()

	// Use MSP only for exceptions handlers. MSP will point to stack at boottom
	// of stacks area, which is at the same time, the beginning of the RAM, so
	// stack overflow in exception handler is always caught (even if MPU isn't
	// used).
	cortexm.SetMSP(unsafe.Pointer(stackTop(len(ts.tasks))))

	// After reset all exceptions have highest priority. Change this to achieve:
	// 1. SysTick has higher priority than SVC to ensure that any user of Uptime 	//    syscall observes both systick counter reset and sysCounter update as
	//    one atomic operation.
	// 2. PendSV has lowest priority (can be preempt by any exception).
	// 3. SVC can be used in external interrupt handlers.
	// 4. It should be possible to increase the priority of any external
	//    interrupt to allow it to preempt SVC.
	prange := cortexm.PrioHighest - cortexm.PrioLowest
	scb.PRI_SVCall.Store(cortexm.PrioLowest + prange/2)
	scb.PRI_PendSV.Store(cortexm.PrioLowest)
	for irq := nvic.IRQ(0); irq < 240; irq++ {
		irq.SetPrio(cortexm.PrioLowest + prange/4)
	}

	// Setup MPU.
	//mpu.SetMode(mpu.PrivDef)
	//mpu.Enable()

	// Set taskInfo for initial (current) task.
	ts.tasks[0].Init(0)
	ts.tasks[0].SetState(taskReady)

	// Run tasker.
	//sysTickStart()

	// Leave privileged level.
	cortexm.SetCONTROL(cortexm.CONTROL() | cortexm.Unpriv)
}

func (ts *taskSched) newTask(pc uintptr, psr uint32, lock bool) (tid int, err syscall.Errno) {
	n := ts.curTask
	for {
		if n++; n >= len(ts.tasks) {
			n = 0
		}
		if ts.tasks[n].State() == taskEmpty {
			break
		}
		if n == ts.curTask {
			return 0, syscall.ENORES
		}
	}

	sf, sp := allocStackFrame(stackTop(n))
	sf.PSR = psr // Use parent's PSR as initial PSR for new task.
	sf.PC = pc

	newt := &ts.tasks[n]
	newt.Init(ts.curTask)
	newt.sp = sp
	newt.SetState(taskReady)

	if lock {
		ts.tasks[ts.curTask].SetState(taskLocked)
		raisePendSV()
	}

	return n + 1, syscall.OK
}

func (ts *taskSched) killTask(tid int) syscall.Errno {
	n := ts.curTask
	if tid != 0 {
		n = tid - 1.
	}
	if n >= len(ts.tasks) || ts.tasks[n].State() == taskEmpty {
		return syscall.ENFOUND
	}
	ts.tasks[n].SetState(taskEmpty)
	for i := range ts.tasks {
		if t := &ts.tasks[i]; int(t.parent) == n {
			t.parent = -1
		}
	}
	if n == ts.curTask {
		raisePendSV()
	}
	return syscall.OK
}

func (ts *taskSched) unlockParent() {
	parent := ts.tasks[ts.curTask].parent
	if parent == -1 {
		return
	}
	if pt := &ts.tasks[parent]; pt.State() == taskLocked {
		pt.SetState(taskReady)
	}
}

func (ts *taskSched) waitEvent(e syscall.Event) {
	t := &ts.tasks[ts.curTask]
	if e == 0 || t.event&e != 0 {
		t.event = 0
		return
	}
	t.SetState(taskWaitEvent)
	t.event = e
	raisePendSV()
}
