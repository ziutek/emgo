package exti

import "log"

type Lines uint32

type extiRegs struct {
	im   Lines `C:"volatile"`
	em   Lines `C:"volatile"`
	rts  Lines `C:"volatile"`
	fts  Lines `C:"volatile"`
	swie Lines `C:"volatile"`
	p    Lines `C:"volatile"`
}

func IntEnabled() Lines {
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
	regs1.swie |= l
}

func Pending() Lines {
	return regs1.p
}

func (l Lines) ClearPending() {
	regs1.p = l
}

func (l Lines) connect(port uint32) {
	if l >= L15<<1 {
		log.Panic("exti: only lines 0...15 can be connected to the external source")
	}
	for i, r := range regs2 {
		if l&0x0f != 0 {
			if l&1 != 0 {
				r = r&^0x000f | port
			}
			if l&2 != 0 {
				r = r&^0x00f0 | port<<4
			}
			if l&4 != 0 {
				r = r&^0x0f00 | port<<8
			}
			if l&8 != 0 {
				r = r&^0xf000 | port<<16
			}
			regs2[i] = r
		}
		l = l >> 4
	}
}
