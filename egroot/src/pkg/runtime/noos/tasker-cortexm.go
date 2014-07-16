// +build cortexm3 cortexm4 cortexm4f

package noos

import (
	"math/rand"
	"sync/barrier"
	"unsafe"

	"cortexm"
	"cortexm/exce"
	"cortexm/systick"
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
	ti.rng.Seed(Ticks() ^ (uint64(systick.Val()) << 32))
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
	vt[exce.NMI] = exce.VectorFor(nmiHandler)
	vt[exce.HardFault] = exce.VectorFor(hardFaultHandler)
	vt[exce.MemFault] = exce.VectorFor(memFaultHandler)
	vt[exce.BusFault] = exce.VectorFor(busFaultHandler)
	vt[exce.UsageFault] = exce.VectorFor(usageFaultHandler)
	vt[exce.SVCall] = exce.VectorFor(svcHandler)
	vt[exce.PendSV] = exce.VectorFor(pendSVHandler)
	vt[exce.SysTick] = exce.VectorFor(sysTickHandler)
	exce.UseTable(vt)

	exce.MemFault.Enable()
	exce.BusFault.Enable()
	exce.UsageFault.Enable()

	exce.SVCall.SetPrio(exce.Lowest)
	exce.PendSV.SetPrio(exce.Lowest)

	// Start SysTick before setup taskInfo for initial task to
	// allow correct rng seed.
	sysTickStart()

	// Set taskInfo for initial (current) task.
	ts.tasks[0].init()
	ts.tasks[0].setState(taskReady)

	tasker.forceNext = -1
	tasker.run()
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
