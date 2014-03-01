package nvic

import "unsafe"

type IRQ byte

const (
	Reset      IRQ = iota + 1 // prio -3 (fixed)
	NMI                       // prio -2 (fixed)
	HardFault                 // prio -1 (fixed)
	MemManage                 // prio 0
	BusFault                  // prio 1
	UsageFault                // prio 2
	_
	_
	_
	_
	SVCall   // prio 3
	DebugMon // prio 4
	_
	PendSV  // prio 5
	SysTick // prio 6
)

type bitReg struct {
	r [8]uint32 `C:"volatile"`
}

func (b *bitReg) setBit(n IRQ) {
	val := uint32(1) << (n & 31)
	n >>= 5
	b.r[n] = val
}

func (b *bitReg) bit(n IRQ) bool {
	val := uint32(1) << (n & 31)
	n >>= 5
	return b.r[n]&val != 0
}

type byteReg struct {
	r [60 * 4]byte `C:"volatile"`
}

type wordReg struct {
	r uint32 `C:"volatile"`
}

var (
	ise = (*bitReg)(unsafe.Pointer(uintptr(0xe000e100)))
	ice = (*bitReg)(unsafe.Pointer(uintptr(0xe000e180)))
	isp = (*bitReg)(unsafe.Pointer(uintptr(0xe000e200)))
	icp = (*bitReg)(unsafe.Pointer(uintptr(0xe000e280)))
	iab = (*bitReg)(unsafe.Pointer(uintptr(0xe000e300)))
	ip  = (*byteReg)(unsafe.Pointer(uintptr(0xe000e400)))
	sti = (*wordReg)(unsafe.Pointer(uintptr(0xe000eF00)))
)

func (irq IRQ) Enable() {
	ise.setBit(irq)
}

func (irq IRQ) Enabled() bool {
	return ise.bit(irq)
}

func (irq IRQ) Disable() {
	ice.setBit(irq)
}

func (irq IRQ) SetPending() {
	isp.setBit(irq)
}

func (irq IRQ) Pending() bool {
	return isp.bit(irq)
}

func (irq IRQ) ClearPending() {
	icp.setBit(irq)
}

func (irq IRQ) Active() bool {
	return iab.bit(irq)
}

func (irq IRQ) SetPriority(prio byte) {
	ip.r[irq] = prio
}

func (irq IRQ) Priority() byte {
	return ip.r[irq]
}

func (irq IRQ) Trig() {
	sti.r = uint32(irq)
}
