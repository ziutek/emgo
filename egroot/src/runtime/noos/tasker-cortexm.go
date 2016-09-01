// +build cortexm0 cortexm3 cortexm4 cortexm4f

package noos

import (
	"math/rand"
	"syscall"
	"unsafe"

	"arch/cortexm"
	"arch/cortexm/mpu"
	"arch/cortexm/nvic"
	"arch/cortexm/scb"
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

func (ti *taskInfo) Init(parent int, t int64) {
	*ti = taskInfo{parent: int16(parent), prio: 255}
	if t == 0 {
		ti.rng.Seed(uint64(uintptr(unsafe.Pointer(ti))))
	} else {
		ti.rng.Seed(uint64(t + 1))
	}
}

func (ti *taskInfo) State() taskState {
	return taskState(ti.flags & 3)
}

func (ti *taskInfo) SetState(s taskState) {
	ti.flags = ti.flags&^3 | s
}

// Comment for future separate tasker package:
// 1. Exported methods can be called only from thread mode (using SVC handler)
//    or from PendSV handler.
// 2. Describe which method can be called from which handler.
type taskSched struct {
	alarm   int64
	period  uint32
	nanosec func() int64
	wakeup  func(int64)
	tasks   []taskInfo
	curTask int
}

func dummyNanosec() int64 { return 0 }
func dummyWakeup(int64)   {}

var tasker = taskSched{
	alarm:   1<<63 - 1,
	period:  2e6, // 2 ms
	nanosec: dummyNanosec,
	wakeup:  dummyWakeup,
}

func (ts *taskSched) Nanosec() int64 {
	return ts.nanosec()
}

func (ts *taskSched) SetNanosec(nanosec func() int64) {
	if nanosec == nil {
		ts.nanosec = dummyNanosec
	} else {
		ts.nanosec = nanosec
	}
}

func (ts *taskSched) SetWakeup(wakeup func(int64)) {
	if wakeup == nil {
		ts.wakeup = dummyWakeup
	} else {
		ts.wakeup = wakeup
	}
}

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

const (
	regionFlash = iota
	regionSRAM
	regionPeriph
	regionStackGuard
)

func (ts *taskSched) init() {
	ts.tasks = make([]taskInfo, maxTasks())

	// Use PSP as stack pointer for thread mode. Current (zero) task has stack
	// at top of the stacks area.
	cortexm.SetPSP(unsafe.Pointer(cortexm.MSP()))
	cortexm.DSB()
	cortexm.SetCONTROL(cortexm.CONTROL() | cortexm.UsePSP)
	cortexm.ISB()

	// Use MSP only for exceptions handlers. MSP will point to stack at boottom
	// of stacks area, which is at the same time, the beginning of the RAM, so
	// stack overflow in exception handler is always caught (even if MPU isn't
	// used).
	cortexm.SetMSP(unsafe.Pointer(stacksBegin() + isrStackSize()))

	// After reset all exceptions have highest priority. Change this to
	// achieve folowing assumptions (even in case of Cortex-M0 which has only
	// four effective priority levels:
	// 1. PendSV has absolutely lowest priority.
	// 2. SVC has lowest possible priority that is effectively higher than
	//    priority of PendSV exception.
	spnum := cortexm.PrioStep * cortexm.PrioNum
	SCB := scb.SCB
	SCB.PRI_PendSV().Store(scb.PRI_PendSV.J(cortexm.PrioLowest + spnum*0/4))
	SCB.PRI_SVCall().Store(scb.PRI_SVCall.J(cortexm.PrioLowest + spnum*1/4))
	for irq := nvic.IRQ(0); irq < 240; irq++ {
		irq.SetPrio(cortexm.PrioLowest + spnum*2/4)
	}

	// Setup MPU. MPU is used mainly to detect stack overflows: regionStackGuard
	// is used to setup stack guard area at the bottom of the stack of active
	// task. Below there is configuration that more or less corresponds to
	// default behavior without MPU enabled. Stack guard region for active task
	// is configured by setStackGuard function.
	mpu.SetRegion(mpu.BaseAttr{
		0x00000000 + mpu.VALID + regionFlash,
		mpu.ENA | mpu.SIZE(29) | mpu.C | mpu.Ar_r_,
	})
	mpu.SetRegion(mpu.BaseAttr{
		0x20000000 + mpu.VALID + regionSRAM,
		mpu.ENA | mpu.SIZE(29) | mpu.C | mpu.S | mpu.Arwrw,
	})
	mpu.SetRegion(mpu.BaseAttr{
		0x40000000 + mpu.VALID + regionPeriph,
		mpu.ENA | mpu.SIZE(29) | mpu.B | mpu.S | mpu.Arwrw | mpu.XN,
	})
	mpu.Set(mpu.ENABLE | mpu.PRIVDEFENA)

	// Set taskInfo for initial (current) task.
	ts.tasks[0].Init(0, ts.nanosec())
	ts.tasks[0].SetState(taskReady)
	setStackGuard(0)

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
	newt.Init(ts.curTask, ts.nanosec())
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

// SetAlarm can be called only from thread mode (through SVCall).
func (ts *taskSched) SetAlarm(t int64) {
	if t > ts.alarm {
		return
	}
	// Can be read only in SVC or PendSV so non-atomic assignment is good.
	ts.alarm = t
}
