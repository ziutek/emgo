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
	// Don't use "r [60*4]byte" because Cortex-M0 doesn't support byte access
	// to this registers.
	r [60]uint32 `C:"volatile"`
}

func (b *byteReg) setByte(n IRQ, v byte) {
	shift := (n & 3) * 8
	val := uint32(v) << shift
	mask := uint32(0xff) << shift
	n >>= 2
	b.r[n] = b.r[n]&^mask | val
}

func (b *byteReg) byte(n IRQ) byte {
	shift := (n & 3) * 8
	return byte(b.r[n>>2] >> shift)
}
