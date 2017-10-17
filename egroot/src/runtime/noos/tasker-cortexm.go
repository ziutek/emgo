// +build cortexm0 cortexm3 cortexm4 cortexm4f cortexm7f cortexm7d

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

type task struct {
	info *taskInfo
	rng  rand.XorShift64
	at   int64
}

// Comment for future separate tasker package:
// 1. Exported methods can be called only from thread mode (using SVC handler)
//    or from PendSV handler.
// 2. Describe which method can be called from which handler.
type taskSched struct {
	alarm     int64
	period    uint32
	nanosec   func() int64
	setWakeup func(int64, bool)
	tasks     []task
	curTask   int
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

func (ts *taskSched) SetSysTimer(nanosec func() int64, setWakeup func(int64, bool)) {
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
}

func (ts *taskSched) deliverEvent(e syscall.Event) {
	for n := range ts.tasks {
		ti := ts.tasks[n].info
		switch ti.state() {
		case taskEmpty:
			// skip
		case taskWaitEvent:
			if ti.event&e != 0 {
				ti.event = 0
				ti.setState(taskReady)
			}
		default:
			ti.event |= e
		}
	}
}

const (
	mpuFlash = iota
	mpuSRAM
	mpuPeriph
	mpuExtRAM
	mpuExtDev
	mpuStackGuard
)

func (ts *taskSched) init() {
	if unsafe.Sizeof(taskInfo{}) > 24 {
		panic("noos: taskInfo to big")
	}

	ts.tasks = make([]task, maxTasks())
	for n := range ts.tasks {
		// Place taskinfo for n-th task in its stack guard slot.
		resetStackGuard(n)
		ts.tasks[n].info = (*taskInfo)(unsafe.Pointer(stackGuardBegin(n)))
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

	if useMPU {
		// MPU is used mainly to detect stack overflows: mpuStackGuard region is
		// used to setup the stack guard area at the bottom of the stack of
		// active task (setStackGuard function is used).
		//
		// Below there is the MPU configuration that more or less corresponds to
		// the default behavior, without MPU enabled, but all RAM and peripheral
		// regions are declared as shareable (usually shared with DMA). In case
		// of Cortex-M7, shareable regions are not cacheable (set acc.SIWT bit
		// to allow use cache but only in WT mode or reconfigure MPU).
		var (
			maPeriph = mpu.B | mpu.Arwrw | mpu.XN
			maFlash  = mpu.C | mpu.Ar_r_
			maRAM    = mpu.TEX(1) | mpu.C | mpu.B | mpu.S | mpu.Arwrw
		)
		mpu.SetRegion(mpu.BaseAttr{
			0x00000000 | mpu.VALID | mpuFlash,
			mpu.ENA | mpu.SIZE(29) | maFlash,
		})
		mpu.SetRegion(mpu.BaseAttr{
			0x20000000 | mpu.VALID | mpuSRAM,
			mpu.ENA | mpu.SIZE(29) | maRAM,
		})
		mpu.SetRegion(mpu.BaseAttr{
			0x40000000 | mpu.VALID | mpuPeriph,
			mpu.ENA | mpu.SIZE(29) | maPeriph,
		})
		mpu.SetRegion(mpu.BaseAttr{
			0x60000000 | mpu.VALID | mpuExtRAM,
			mpu.ENA | mpu.SIZE(30) | maRAM,
		})
		mpu.SetRegion(mpu.BaseAttr{
			0xA0000000 | mpu.VALID | mpuExtDev,
			mpu.ENA | mpu.SIZE(29) | maPeriph,
		})
		mpu.Set(mpu.ENABLE | mpu.PRIVDEFENA)
		cortexm.DSB()
		cortexm.ISB()
	}

	// Set taskInfo for initial (current) task.
	ts.initTask(0, 0)
	if useMPU {
		setMPUStackGuard(0)
	}

	// Leave privileged level.
	cortexm.SetCONTROL(cortexm.CONTROL() | cortexm.Unpriv)
}

func (ts *taskSched) initTask(n int, sp uintptr) {
	t := &ts.tasks[n]
	t.info.init(ts.curTask, sp)
	ns := ts.nanosec()
	if ns == 0 {
		t.rng.Seed(int64(uintptr(unsafe.Pointer(t))))
	} else {
		t.rng.Seed(ns + 1)
	}
}

func (ts *taskSched) newTask(pc uintptr, psr uint32, lock bool) (tid int, err syscall.Errno) {
	n := ts.curTask
	for {
		if n++; n >= len(ts.tasks) {
			n = 0
		}
		if ts.tasks[n].info.state() == taskEmpty {
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
		ts.tasks[ts.curTask].info.setState(taskLocked)
	}
	raisePendSV() // Ensure that scheduler knowns that has more than one task.
	return n + 1, syscall.OK
}

func (ts *taskSched) killTask(tid int) syscall.Errno {
	n := ts.curTask
	if tid != 0 {
		n = tid - 1.
	}
	ti := ts.tasks[n].info
	if n >= len(ts.tasks) || ti.state() == taskEmpty {
		return syscall.ENFOUND
	}
	ti.setState(taskEmpty)
	for i := range ts.tasks {
		if ti := ts.tasks[i].info; int(ti.parent) == n {
			ti.parent = -1
		}
	}
	if n == ts.curTask {
		raisePendSV()
	}
	return syscall.OK
}

func (ts *taskSched) unlockParent() {
	parent := ts.tasks[ts.curTask].info.parent
	if parent == -1 {
		return
	}
	if ti := ts.tasks[parent].info; ti.state() == taskLocked {
		ti.setState(taskReady)
	}
}

func (ts *taskSched) waitEvent(e syscall.Event) {
	ti := ts.tasks[ts.curTask].info
	if e == 0 || ti.event&e != 0 {
		ti.event = 0
		return
	}
	ti.setState(taskWaitEvent)
	ti.event = e
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

// SetAt can be called only from thread mode (through SVCall).
func (ts *taskSched) SetAt(t int64) {
	tasker.tasks[tasker.curTask].at = t
	ts.SetAlarm(t)
}
