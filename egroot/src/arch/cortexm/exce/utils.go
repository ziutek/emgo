package exce

type bitReg struct {
	r [8]uint32
} //c:volatile

func (b *bitReg) setBit(e Exce) {
	val := uint32(1) << uint(e & 31)
	e >>= 5
	b.r[e] = val
}

func (b *bitReg) bit(e Exce) bool {
	mask := uint32(1) << uint(e & 31)
	e >>= 5
	return b.r[e]&mask != 0
}

type byteReg struct {
	// Don't use "r [60*4]byte" because Cortex-M0 doesn't support byte access
	// to this registers.
	r [60]uint32
} //c:volatile

func (b *byteReg) setByte(e Exce, v byte) {
	shift := uint(e & 3) * 8
	val := uint32(v) << shift
	mask := uint32(0xff) << shift
	e >>= 2
	b.r[e] = b.r[e]&^mask | val
}

func (b *byteReg) byte(e Exce) byte {
	shift := uint(e & 3) * 8
	return byte(b.r[e>>2] >> shift)
}
