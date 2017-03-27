package ppi

import (
	"mmio"
	"unsafe"

	"nrf5/hal/te"

	"nrf5/hal/internal/mmap"
)

type channel struct {
	eep mmio.U32
	tep mmio.U32
}

type Periph struct {
	te.Regs

	_       [64]mmio.U32
	chen    mmio.U32
	chenset mmio.U32
	chenclr mmio.U32
	_       mmio.U32
	ch      [20]channel
	_       [148]mmio.U32
	chg     [6]mmio.U32
	_       [62]mmio.U32
	forktep [32]mmio.U32
}

//emgo:const
var PPI = (*Periph)(unsafe.Pointer(mmap.BaseAPB + 0x1F000))

type Task byte

const (
	CHG0EN  Task = 0  // Enable channel group 0.
	CHG0DIS Task = 1  // Disable channel group 0.
	CHG1EN  Task = 2  // Enable channel group 1.
	CHG1DIS Task = 3  // Disable channel group 1.
	CHG2EN  Task = 4  // Enable channel group 2.
	CHG2DIS Task = 5  // Disable channel group 2.
	CHG3EN  Task = 6  // Enable channel group 3.
	CHG3DIS Task = 7  // Disable channel group 3.
	CHG4EN  Task = 8  // Enable channel group 4.
	CHG4DIS Task = 9  // Disable channel group 4.
	CHG5EN  Task = 10 // Enable channel group 5.
	CHG5DIS Task = 11 // Disable channel group 5.
)

func (p *Periph) Task(t Task) *te.Task { return p.Regs.Task(int(t)) }

