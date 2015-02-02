package gpio

import "stm32/f4/periph"

func (p *Port) Periph() periph.AHB1Dev {
	return periph.GPIOA << uint(p.Number())
}