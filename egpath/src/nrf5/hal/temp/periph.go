package temp

import (
	"mmio"
	"unsafe"

	"nrf5/hal/internal/mmap"
	"nrf5/hal/te"
)

type Periph struct {
	te.Regs

	_    [66]mmio.U32
	temp mmio.U32
	_    [5]mmio.U32
	a    [6]mmio.U32
	_    [2]mmio.U32
	b    [6]mmio.U32
	_    [2]mmio.U32
	t    [5]mmio.U32
}

//emgo:const
var TEMP = (*Periph)(unsafe.Pointer(mmap.APB_BASE + 0x0C000))

type Task byte

const (
	START Task = 0 // Start temperature measurement.
	STOP  Task = 1 // Stop temperature measurement.
)

type Event byte

const (
	DATARDY Event = 0 // Temperature measurement complete, data ready.
)

func (p *Periph) Task(t Task) *te.Task    { return p.Regs.Task(int(t)) }
func (p *Periph) Event(e Event) *te.Event { return p.Regs.Event(int(e)) }

// LoadTEMP returns temperature (in 0.25 °C units) mesured by last START task.
func (p *Periph) LoadTEMP() int {
	return int(p.temp.Load())
}

func (p *Periph) A() []mmio.U32 {
	return p.a[:]
}

func (p *Periph) B() []mmio.U32 {
	return p.b[:]
}

func (p *Periph) T() []mmio.U32 {
	return p.t[:]
}
