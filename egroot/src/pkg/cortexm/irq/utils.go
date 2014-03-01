package irq

type bitReg struct {
	r [8]uint32 `C:"volatile"`
}

func (b *bitReg) setBit(n IRQ) {
	val := uint32(1) << (n & 31)
	n >>= 5
	b.r[n] = val
}

func (b *bitReg) bit(n IRQ) bool {
	mask := uint32(1) << (n & 31)
	n >>= 5
	return b.r[n]&mask != 0
}

type byteReg struct {
	r [60 * 4]byte `C:"volatile"`
}

type wordReg struct {
	r uint32 `C:"volatile"`
}

func (w *wordReg) setBit(n IRQ) {
	w.r |= uint32(1) << n
}

func (w *wordReg) clearBit(n IRQ) {
	w.r &^= uint32(1) << n
}

func (w *wordReg) bit(n IRQ) bool {
	return w.r&(uint32(1)<<n) != 0
}
