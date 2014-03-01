package scb

type wordReg struct {
	r uint32 `C:"volatile"`
}

func (b *wordReg) setBit(n uint) {
	b.r |= uint32(1) << n
}

func (b *wordReg) clearBit(n uint){
	b.r &^= uint32(1) << n
}
