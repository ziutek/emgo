package exti

type Lines uint32

type extiRegs struct {
	im   Lines
	em   Lines
	rts  Lines
	fts  Lines
	swie Lines
	p    Lines
} //c:volatile

func IntrEnabled() Lines {
	return regs1.im
}

func (l Lines) IntEnable() {
	regs1.im |= l
}

func (l Lines) IntDisable() {
	regs1.im &^= l
}

func EventEnabled() Lines {
	return regs1.em
}

func (l Lines) EventEnable() {
	regs1.em |= l
}

func (l Lines) EventDisable() {
	regs1.em &^= l
}

func RiseTrigEnabled() Lines {
	return regs1.rts
}

func (l Lines) RiseTrigEnable() {
	regs1.rts |= l
}

func (l Lines) RiseTrigDisable() {
	regs1.rts &^= l
}

func FallTrigEnabled() Lines {
	return regs1.fts
}

func (l Lines) FallTrigEnable() {
	regs1.fts |= l
}

func (l Lines) FallTrigDisable() {
	regs1.fts &^= l
}

func SoftReq() Lines {
	return regs1.swie
}

func (l Lines) SoftReqGen() {
	regs1.swie = l
}

func Pending() Lines {
	return regs1.p
}

func (l Lines) ClearPending() {
	regs1.p = l
}

func (l Lines) connect(port uint32) {
	if l >= L15<<1 {
		panic("exti: only lines 0...15 can be connected to the external source")
	}
	for i := range regs2 {
		if l&0x0f != 0 {
			r := &regs2[i]
			if l&1 != 0 {
				r.StoreBits(port, 0x000f)
			}
			if l&2 != 0 {
				r.StoreBits(port<<4, 0x00f0)
			}
			if l&4 != 0 {
				r.StoreBits(port<<8, 0x0f00)
			}
			if l&8 != 0 {
				r.StoreBits(port<<12, 0xf000)
			}
		}
		l = l >> 4
	}
}
