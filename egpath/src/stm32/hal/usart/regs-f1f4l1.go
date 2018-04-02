// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl f10x_cl f40_41xxx f411xe l1xx_md l1xx_mdp l1xx_hd l1xx_xl

package usart

import (
	"stm32/hal/raw/usart"
)

func (p *Periph) status() (Event, Error) {
	sr := p.raw.SR.Load()
	return Event(sr >> 4), Error(sr & 0xf)
}

func (p *Periph) clear(ev Event, _ Error) {
	if ev != 0 {
		p.raw.SR.Store(^(usart.SR(ev) << 4))
	}
}

func (p *Periph) store(d int) {
	p.raw.DR.Store(usart.DR(d))
}

func (p *Periph) load() int {
	return int(p.raw.DR.Load())
}

func (p *Periph) rdAddr() uintptr {
	return p.raw.DR.Addr()
}

func (p *Periph) tdAddr() uintptr {
	return p.raw.DR.Addr()
}
