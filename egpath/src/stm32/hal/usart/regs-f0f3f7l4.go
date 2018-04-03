// +build f030x6 f030x8 f303xe f746xx l476xx

package usart

import (
	"stm32/hal/raw/usart"
)

func (p *Periph) status() (Event, Error) {
	isr := p.raw.ISR.Load()
	return Event(isr >> 4), Error(isr & 0xf)
}

func (p *Periph) clear(ev Event, err Error) {
	raw := &p.raw
	raw.ICR.Store(usart.ICR(ev)<<4 | usart.ICR(err))
	if ev&RxNotEmpty != 0 {
		raw.RQR.Store(usart.RXFRQ)
	}
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
