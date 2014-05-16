// +build cortexm3 cortexm4 cortexm4f

package noos

import (
	"unsafe"

	"sync/barrier"

	"cortexm"
	"cortexm/irq"
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
	flags taskState
	prio  uint8
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
	var vt []irq.Vector
	vtlen := 1 << irtExp()
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

	// Set taskInfo for initial (current) task.
	ts.tasks[0].prio = 255
	ts.tasks[0].setState(taskReady)

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

var Tick uint32

func sysTickHandler() {
	Tick++
	if tasker.onSysTick {
		irq.PendSV.SetPending()
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
