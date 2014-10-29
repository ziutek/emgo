// +build cortexm3 cortexm4 cortexm4f

package noos

import (
	"math/rand"
	"sync/barrier"
	"unsafe"

	"cortexm"
	"cortexm/exce"
	"cortexm/mpu"
)

type taskState byte

const (
	taskEmpty taskState = iota
	taskReady
	taskWaitEvent
)

// taskInfo
// sp contains value of SP after automatic stacking during exception entry. So
// sp points to the last register in set automatically stacked by CPU and just
// after the register set stacked by tasker. pendSVHandler can use two least
// significant bits of sp for its flags.
type taskInfo struct {
	sp    uintptr
	event Event
	rng   rand.XorShift64
	flags taskState
	prio  uint8
}

func (ti *taskInfo) init() {
	*ti = taskInfo{prio: 255}
	ti.rng.Seed(Uptime())
}

func (ti *taskInfo) state() taskState {
	return taskState(ti.flags & 3)
}

func (ti *taskInfo) setState(s taskState) {
	ti.flags = ti.flags&^3 | s
}

type taskSched struct {
	tasks     []taskInfo
	curTask   int
	forceNext int
	onSysTick bool
}

var tasker taskSched

func (ts *taskSched) run() {
	barrier.Memory()
	ts.onSysTick = true
}

func (ts *taskSched) stop() {
	ts.onSysTick = false
	barrier.Memory()
}

func (ts *taskSched) deliverEvent(e Event) {
	for i := range ts.tasks {
		t := &ts.tasks[i]
		switch t.state() {
		case taskEmpty:
			// skip

		case taskWaitEvent:
			if t.event&e != 0 {
				t.event = 0
				t.setState(taskReady)
			}

		default:
			t.event |= e
		}
	}
}

func irtExp() uint

func (ts *taskSched) init() {
	var vt []exce.Vector
	vtlen := 1 << irtExp()
	vtsize := vtlen * int(unsafe.Sizeof(exce.Vector{}))

	Heap = allocTop(
		unsafe.Pointer(&vt), Heap,
		vtlen, unsafe.Sizeof(exce.Vector{}), unsafe.Alignof(exce.Vector{}),
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

	// Use PSP as stack pointer for thread mode. Current (zero) task has stack
	// at top of the stacks area.
	cortexm.SetPSP(unsafe.Pointer(cortexm.MSP()))
	barrier.Sync()
	cortexm.SetCtrl(cortexm.Ctrl() | cortexm.UsePSP)
	cortexm.ISB()

	// Use MSP only for exceptions handlers. MSP will point to stack at boottom
	// of stacks area, which is at the same time, the beginning of the RAM, so
	// stack overflow in exception handler is always caught (even if MPU isn't
	// used).
	cortexm.SetMSP(unsafe.Pointer(stackTop(len(ts.tasks))))

	// Setup interrupt table.
	// Consider setup at link time using GCC weak functions to support Cortex-M0
	// and (in case of Cortex-M3,4) to allow vector load on ICode bus
	// simultaneously with registers stacking on DCode bus.
	vt[exce.NMI] = exce.VectorFor(NMIHandler)
	vt[exce.HardFault] = exce.VectorFor(FaultHandler)
	vt[exce.MemManage] = exce.VectorFor(FaultHandler)
	vt[exce.BusFault] = exce.VectorFor(FaultHandler)
	vt[exce.UsageFault] = exce.VectorFor(FaultHandler)
	vt[exce.SVCall] = exce.VectorFor(svcHandler)
	vt[exce.PendSV] = exce.VectorFor(pendSVHandler)
	vt[exce.SysTick] = exce.VectorFor(sysTickHandler)
	exce.UseTable(vt)

	exce.SysTick.SetPrio(exce.Highest)
	exce.SVCall.SetPrio(exce.Lowest)
	exce.PendSV.SetPrio(exce.Lowest)
	for irq := exce.IRQ0; int(irq) < len(vt); irq++ {
		irq.SetPrio((exce.Lowest + exce.Highest) / 2)
	}
	exce.MemManage.Enable()
	exce.BusFault.Enable()
	exce.UsageFault.Enable()

	// Setup MPU.
	mpu.SetMode(mpu.PrivDef)
	//mpu.Enable()

	// Start SysTick before setup taskInfo for initial task to
	// allow correct rng seed.
	sysTickStart()

	// Set taskInfo for initial (current) task.
	ts.tasks[0].init()
	ts.tasks[0].setState(taskReady)

	tasker.forceNext = -1
	tasker.run()

	// Leave privilege level.
	//cortexm.SetCtrl(cortexm.Ctrl() | cortexm.Unpriv)
}
