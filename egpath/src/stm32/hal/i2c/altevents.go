package i2c

import (
	"rtos"

	"stm32/hal/dma"

	"stm32/hal/raw/i2c"
)

const byteTimeout = 100e6 // 100 ms

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
		ev := ch.Events()
		errmask := dma.ERR &^ dma.FFERR // Ignore FIFO error.
		if ev&errmask != 0 {
			return DMAErr
		}
		if ev&dma.TRCE != 0 {
			return 0
		}
		if rtos.Nanosec() >= deadline {
			return SoftTimeout
		}
	}
}

func dmaWaitTRCE(ch *dma.Channel, evflag *rtos.EventFlag, deadline int64) Error {
	errmask := dma.ERR &^ dma.FFERR // Ignore FIFO error.
	for {
		ch.EnableInt(dma.TRCE | errmask)
		if !evflag.Wait(deadline) {
			return SoftTimeout
		}
		evflag.Clear()
		ev := ch.Events()
		if ev&errmask != 0 {
			return DMAErr
		}
		if ev&dma.TRCE != 0 {
			return 0
		}
	}
}
