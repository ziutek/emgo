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
