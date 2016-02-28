package i2c

import (
	"rtos"

	"arch/cortexm/nvic"

	"stm32/hal/dma"

	"stm32/hal/raw/i2c"
)

type irqs struct {
	evflag rtos.EventFlag
	i2cev  nvic.IRQ
	i2cerr nvic.IRQ
	dmatx  nvic.IRQ
	dmarx  nvic.IRQ
}

func i2cPollEvent(p *i2c.I2C_Periph, ev i2c.SR1_Bits, deadline int64) Error {
	for {
		sr1 := p.SR1.Load()
		if e := Error(sr1 >> 8); e != 0 {
			return e
		}
		if sr1&ev != 0 {
			return 0
		}
		if rtos.Nanosec() >= deadline {
			return SoftTimeout
		}
	}
}

func i2cWaitIRQ(p *i2c.I2C_Periph, irqs *irqs, ev i2c.SR1_Bits, deadline int64) Error {
	for {
		rtos.IRQ(irqs.i2cev).Enable()
		rtos.IRQ(irqs.i2cerr).Enable()
		if !irqs.evflag.Wait(deadline) {
			return SoftTimeout
		}
		irqs.evflag.Clear()
		sr1 := p.SR1.Load()
		if e := Error(sr1 >> 8); e != 0 {
			return e
		}
		if sr1&ev != 0 {
			return 0
		}
	}
}

func i2cWaitEvent(p *i2c.I2C_Periph, irqs *irqs, ev i2c.SR1_Bits) Error {
	deadline := rtos.Nanosec() + 100e6 // 100 ms
	if irqs.i2cev == 0 {
		return i2cPollEvent(p, ev, deadline)
	}
	return i2cWaitIRQ(p, irqs, ev, deadline)
}

func dmaPoolTCE(ch dma.Channel, deadline int64) Error {
	for {
		cur := ch.Events()
		if cur&dma.ERR != 0 {
			return DMAErr
		}
		if cur&dma.TCE != 0 {
			return 0
		}
		if rtos.Nanosec() >= deadline {
			return SoftTimeout
		}
	}
}
