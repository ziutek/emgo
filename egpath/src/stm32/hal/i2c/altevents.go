package i2c

import (
	"rtos"

	"stm32/hal/dma"

	"stm32/hal/raw/i2c"
)

func i2cPollEvent(p *i2c.I2C_Periph, ev i2c.SR1_Bits, deadline int64) Error {
	for {
		sr1 := p.SR1.Load()
		if e := getError(sr1); e != 0 {
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

func i2cWaitIRQ(p *i2c.I2C_Periph, evflag *rtos.EventFlag, ev i2c.SR1_Bits, deadline int64) Error {
	irqen := i2c.ITEVTEN | i2c.ITERREN
	if ev&(i2c.RXNE|i2c.TXE) != 0 {
		irqen |= i2c.ITBUFEN
	}
	for {
		p.CR2.SetBits(irqen)
		if !evflag.Wait(deadline) {
			return SoftTimeout
		}
		evflag.Clear()
		sr1 := p.SR1.Load()
		if e := getError(sr1); e != 0 {
			return e
		}
		if sr1&ev != 0 {
			return 0
		}
	}
}

func dmaPoolTRCE(ch *dma.Channel, deadline int64) Error {
	for {
		ev, err := ch.Status()
		if err&dmaErrMask != 0 {
			return DMAErr
		}
		if ev&dma.Complete != 0 {
			return 0
		}
		if rtos.Nanosec() >= deadline {
			ch.Disable()
			return SoftTimeout
		}
	}
}

func dmaWait(ch *dma.Channel, evflag *rtos.EventFlag, deadline int64) Error {
	done := evflag.Wait(deadline)
	ch.Disable() // To be compatible with STM32F1.
	if !done {
		ch.DisableIRQ(dma.EvAll, dma.ErrAll)
		return SoftTimeout
	}
	_, err := ch.Status()
	if err&dmaErrMask != 0 {
		return DMAErr
	}
	return 0
}
