package mmio

type Reg8 struct {
	r uint8 `C:"volatile"`
}

func (r *Reg8) SetBit(n int) {
	r.r |= uint8(1) << uint(n)
}

func (r *Reg8) ClearBit(n int) {
	r.r &^= uint8(1) << uint(n)
}

func (r *Reg8) Bit(n int) bool {
	return r.r&(uint8(1)<<uint(n)) != 0
}

func (r *Reg8) Read() uint8 {
	return r.r
}

func (r *Reg8) Write(v uint8) {
	r.r = v
}

type Reg16 struct {
	r uint16 `C:"volatile"`
}

func (r *Reg16) SetBit(n int) {
	r.r |= uint16(1) << uint(n)
}

func (r *Reg16) ClearBit(n int) {
	r.r &^= uint16(1) << uint(n)
}

func (r *Reg16) Bit(n int) bool {
	return r.r&(uint16(1)<<uint(n)) != 0
}

func (r *Reg16) Read() uint16 {
	return r.r
}

func (r *Reg16) Write(v uint16) {
	r.r = v
}

type Reg32 struct {
	r uint32 `C:"volatile"`
}

func (r *Reg32) SetBit(n int) {
	r.r |= uint32(1) << uint(n)
}

func (r *Reg32) ClearBit(n int) {
	r.r &^= uint32(1) << uint(n)
}

func (r *Reg32) Bit(n int) bool {
	return r.r&(uint32(1)<<uint(n)) != 0
}

func (r *Reg32) Read() uint32 {
	return r.r
}

func (r *Reg32) Write(v uint32) {
	r.r = v
}
