package nvic

import "unsafe"

type Vector func()

type Table struct {
	_          Vector `C:"__attribute__((aligned(16*4)))"`
	Reset      Vector // prio -3 (fixed)
	NMI        Vector // prio -2 (fixed)
	HardFault  Vector // prio -1 (fixed)
	MemManage  Vector // prio 0
	BusFault   Vector // prio 1
	UsageFault Vector // prio 2
	_          Vector
	_          Vector
	_          Vector
	_          Vector
	SVCall     Vector // prio 3
	DebugMon   Vector // prio 4
	_          Vector
	PendSV     Vector // prio 5
	SysTick    Vector // prio 6
}

func (t *Table) Set(irq IRQ, v Vector) {
	(*[16]Vector)(unsafe.Pointer(t))[irq] = v
}
