// +build f746xx l476xx

package usart

import (
	"stm32/hal/raw/usart"
)

const lbd = usart.LBDF

func (p *Periph) status() (Event, Error) {
	isr := p.raw.ISR.Load()
	return Event(isr >> 4), Error(isr & 0xf)
}

func (p *Periph) clear(ev Event, err Error) {
	p.raw.ICR.Store(usart.ICR(ev)<<4 | usart.ICR(err))
}

func (p *Periph) store(d int) {
	p.raw.TDR.Store(usart.TDR(d))
}

func (p *Periph) load() int {
	return int(p.raw.RDR.Load())
}

func (p *Periph) rdAddr() uintptr {
	return p.raw.RDR.Addr()
}

func (p *Periph) tdAddr() uintptr {
	return p.raw.TDR.Addr()
}
