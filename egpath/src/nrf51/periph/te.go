package periph

import (
	"arch/cortexm/exce"
	"unsafe"
)

// TasksEvents implements tasks and events.
// It should be the first field on any Periph struct.
// It takes 0x400 bytes of memory.
type TasksEvents struct {
	tasks    [32]uint32
	_        [32]uint32
	events   [32]uint32
	_        [32]uint32
	shorts   uint32
	_        [64]uint32
	intenset uint32
	intenclr uint32
	_        [14]uint32
	evtenset uint32 // Enable event routing to PPI.
	evtenclr uint32 // Disable event routing to PPI.
	_        [45]uint32
} //c:volatile

type Task byte

// TrigTask triggers task t.
func (te *TasksEvents) TrigTask(t Task) {
	te.tasks[t] = 1
}

type Event byte

// Event returns true if event e occured.
func (te *TasksEvents) Event(e Event) bool {
	return te.events[e] != 0
}

// ClearEvent clears event flag for event e.
func (te *TasksEvents) ClearEvent(e Event) {
	te.events[e] = 0
}

type Shorts uint32

// Shorts returns value of shortcuts register.
func (te *TasksEvents) Shorts() Shorts {
	return Shorts(te.shorts)
}

// SetShorts sets value of shortcuts register.
func (te *TasksEvents) SetShorts(s Shorts) {
	te.shorts = uint32(s)
}

// IRQ returns exception number associated to events.
func (te *TasksEvents) IRQ() exce.Exce {
	addr := uintptr(unsafe.Pointer(te))
	return exce.IRQ0 + exce.Exce((addr-BaseAPB)>>12)
}

// IntEnabled tells whether the event e generates interrupts.
func (te *TasksEvents) IntEnabled(e Event) bool {
	return te.intenset&(uint32(1)<<e) != 0
}

// EnableInt enables generating interrupts by event e.
func (te *TasksEvents) EnableInt(e Event) {
	te.intenset = uint32(1) << e
}

// DisableInt disables generating interrupts by event e.
func (te *TasksEvents) DisableInt(e Event) {
	te.intenclr |= uint32(1) << e
}

// EventEnabled tells whether the event e is enabled. Only
// some peripherals (eg. RTC) implements this method.
func (te *TasksEvents) EventEnabled(e Event) bool {
	return te.evtenset&(uint32(1)<<e) != 0
}

// EnableEvent enables generating event e. Only some
// peripherals (eg. RTC) implements this method.
func (te *TasksEvents) EnableEvent(e Event) {
	te.evtenset = uint32(1) << e
}

// DisableEvent disables generating event e. Only
// some peripherals (eg. RTC) implements this method.
func (te *TasksEvents) DisableEvent(e Event) {
	te.evtenclr |= uint32(1) << e
}
