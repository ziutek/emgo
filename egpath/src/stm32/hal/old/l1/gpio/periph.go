package gpio

import "stm32/l1/periph"

func (p *Port) Periph() periph.AHBDev {
	return periph.GPIOA << uint(p.Number())
}