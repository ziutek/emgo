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

// taskInfo stores information about task state.
//
// sp contains value of SP after automatic stacking during exception entry. So
// sp points to the last register in set automatically stacked by CPU and just
// after the register set stacked by tasker. pendSVHandler can use two least
// significant bits of sp for its flags.
//
// sizeof(taskInfo32) must be less or equal to 24 bytes to fit into 32 B stack
// guard slot (at last 8 bytes must be left free for stack guard algorithm).
type taskInfo struct {
	event  syscall.Event
	parent int16
	flags  taskState
	prio   byte
	sp     uintptr
}

func (ti *taskInfo) init(parent int, sp uintptr) {
	*ti = taskInfo{sp: sp, parent: int16(parent), flags: taskReady}
}

func (ti *taskInfo) state() taskState {
	return taskState(ti.flags & 3)
}

func (ti *taskInfo) setState(s taskState) {
	ti.flags = ti.flags&^3 | s
}

// Comment for future separate tasker package:
// 1. Exported methods can be called only from thread mode (using SVC handler)
//    or from PendSV handler.
// 2. Describe which method can be called from which handler.
type taskSched struct {
	alarm      int64
	lastAlarm  int64
	checkAlarm uint32
	period     uint32
	nanosec    func() int64
	setWakeup  func(int64, bool)
	tasks      []*taskInfo
	rng        []rand.XorShift64
	curTask    int
}

func dummyNanosec() int64                { return 0 }
func dummySetWakeup(t int64, alarm bool) {}

const noalarm = 1<<63 - 1

var tasker = taskSched{
	alarm:     noalarm,
	period:    2e6, // 2 ms
	nanosec:   dummyNanosec,
	setWakeup: dummySetWakeup,
}

func (ts *taskSched) Nanosec() int64 {
	return ts.nanosec()
}

func (ts *taskSched) SetSysTimer(nanosec func() int64, setWakeup func(int64, bool)) *uint32 {
	if nanosec == nil {
		ts.nanosec = dummyNanosec
	} else {
		ts.nanosec = nanosec
	}
	if setWakeup == nil {
		ts.setWakeup = dummySetWakeup
	} else {
		ts.setWakeup = setWakeup
	}
	return &ts.checkAlarm
}

func (ts *taskSched) deliverEvent(e syscall.Event) {
	for _, t := range ts.tasks {
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

const (
	regionFlash = iota
	regionSRAM
	regionPeriph
	regionStackGuard
)

func (ts *taskSched) init() {
	if unsafe.Sizeof(taskInfo{}) > 24 {
		panic("noos: taskInfo to big")
	}

	ts.tasks = make([]*taskInfo, maxTasks())
	ts.rng = make([]rand.XorShift64, maxTasks())
	for n := range ts.tasks {
		// Place taskinfo for n-th task in its stack guard slot.
		resetStackGuard(n)
		ts.tasks[n] = (*taskInfo)(unsafe.Pointer(stackGuardBegin(n)))
	}

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

	// Change the default exceptionpriority levels to achieve folowing
	// assumptions (even in case of Cortex-M0 which has only four effective
	// priority levels:
	// 1. PendSV has absolutely lowest priority.
	// 2. There is an effective priority level between PendSV and SVCall that
	//    can be used by ISRs that want to use system calls.
	// 3. There is effective priority level above SVCall for IRQs that need
	//    preempting system calls (SVCall should not use highest priority).
	SCB := scb.SCB
	SCB.PRI_PendSV().Store(scb.PRI_PendSV.J(cortexm.PrioLowest))
	SCB.PRI_SVCall().Store(scb.PRI_SVCall.J(syscall.SyscallPrio))
	for irq := nvic.IRQ(0); irq < 240; irq++ {
		irq.SetPrio(cortexm.PrioHighest)
	}

	if hasMPU {
		// Setup MPU. MPU is used mainly to detect stack overflows:
		// regionStackGuard is used to setup stack guard area at the bottom of
		// the stack of active task. Below there is configuration that more or
		// less corresponds to default behavior without MPU enabled. Stack guard
		//  region for active task is configured by setStackGuard function.
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
	}

	// Set taskInfo for initial (current) task.
	ts.initTask(0, 0)
	if hasMPU {
		setMPUStackGuard(0)
	}

	// Leave privileged level.
	cortexm.SetCONTROL(cortexm.CONTROL() | cortexm.Unpriv)
}

func (ts *taskSched) initTask(n int, sp uintptr) {
	t := ts.tasks[n]
	t.init(ts.curTask, sp)
	ns := ts.nanosec()
	if ns == 0 {
		ts.rng[n].Seed(uint64(uintptr(unsafe.Pointer(t))))
	} else {
		ts.rng[n].Seed(uint64(ns + 1))
	}
}

func (ts *taskSched) newTask(pc uintptr, psr uint32, lock bool) (tid int, err syscall.Errno) {
	n := ts.curTask
	for {
		if n++; n >= len(ts.tasks) {
			n = 0
		}
		if ts.tasks[n].state() == taskEmpty {
			break
		}
		if n == ts.curTask {
			return 0, syscall.ENORES
		}
	}

	sf, sp := allocStackFrame(stackTop(n))
	sf.PSR = psr // Use parent's PSR as initial PSR for new task.
	sf.PC = pc

	ts.initTask(n, sp)

	if lock {
		ts.tasks[ts.curTask].setState(taskLocked)
		raisePendSV()
	}
	return n + 1, syscall.OK
}

func (ts *taskSched) killTask(tid int) syscall.Errno {
	n := ts.curTask
	if tid != 0 {
		n = tid - 1.
	}
	if n >= len(ts.tasks) || ts.tasks[n].state() == taskEmpty {
		return syscall.ENFOUND
	}
	ts.tasks[n].setState(taskEmpty)
	for i := range ts.tasks {
		if t := ts.tasks[i]; int(t.parent) == n {
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
	if pt := ts.tasks[parent]; pt.state() == taskLocked {
		pt.setState(taskReady)
	}
}

func (ts *taskSched) waitEvent(e syscall.Event) {
	t := ts.tasks[ts.curTask]
	if e == 0 || t.event&e != 0 {
		t.event = 0
		return
	}
	t.setState(taskWaitEvent)
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
