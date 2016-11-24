// +build noos

package syscall

import (
	"internal"
	"sync/atomic"
	"sync/fence"
	"unsafe"
)

// Event bitmask that represents events that multiple tasks or ISRs can send
// and multiple tasks (but not ISRs) can wait for. Events are intended to be
// used by low-level library rutines to implement higher level communication and
// synchronization primitives like channels and mutexes. They are specific to
// noos runtime so can be unavailable if RTOS is used.
type Event uintptr

const eventBits = uintptr(unsafe.Sizeof(Event(0)) * 8)

var (
	eventReg Event
	gen      uintptr
)

// AssignEvent returns event from some internal event pool.
// There is no any guarantee that subsequent calls of AssignEvent assigns
// different events, which means that AssignEvent can return Event already
// assigned by current or another task.
func AssignEvent() Event {
	u := atomic.AddUintptr(&gen, 1)
	return Event(1) << (u & (eventBits - 1))
}

// AssignEventFlag works like AssignEvent but guarantees that the least
// significant bit of returned value is zero.
func AssignEventFlag() Event {
	u := atomic.AddUintptr(&gen, 1)
	return Event(2) << (u % (eventBits - 1))
}

// Send sends event that means it waking up all tasks that wait for e. If some 
// task isn't waiting for any event, e is saved for this task for possible
// future call of Wait. Send do not execute until all memory write operations
// before it, in program order, will be completed. In very simple implementation
// Send can do nothing.
func (e Event) Send() {
	fence.W_SMP() // Send signals that some shared variable was changed.
	atomic.OrUintptr((*uintptr)(&eventReg), uintptr(e))
	schedNext()
}

// TakeEventReg is intended to be used by runtime to obtain accumulated events.
// It returns value of shared event register and clears it in one atomic
// operation.
func TakeEventReg() Event {
	return Event(atomic.SwapUintptr((*uintptr)(&eventReg), 0))
}

// Wait waits for event. If e == 0 it returns immediately. Wait clears all saved
// events for current task so the information about recieved events, that Wait
// hasn't waited for, is lost. Wait ensures that any memory operation after it,
// in program order, do not be executed until Wait return. In very simple
// implementation Wait can do nothing and always return immediatelly.
func (e Event) Wait() {
	internal.Syscall1(EVENTWAIT, uintptr(e))
}

// Alarm is an event that is sent by runtime when asked by using SetAlarm.
var Alarm = AssignEventFlag()

func AtomicLoadEvent(p *Event) Event {
	return Event(atomic.LoadUintptr((*uintptr)(p)))
}

func AtomicStoreEvent(p *Event, e Event) {
	atomic.StoreUintptr((*uintptr)(p), uintptr(e))
}

func AtomicAddEvent(p *Event, delta int) Event {
	return Event(atomic.AddUintptr((*uintptr)(p), uintptr(delta)))
}

func AtomicCompareAndSwapEvent(p *Event, old, new Event) bool {
	return atomic.CompareAndSwapUintptr(
		(*uintptr)(p), uintptr(old), uintptr(new),
	)
}
