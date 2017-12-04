package evedci

import (
	"rtos"

	"stm32/hal/exti"
	"stm32/hal/gpio"
	"stm32/hal/spi"
)

// SPI implements eve.DCI using Serial Peripheral Interface.
type SPI struct {
	spi      *spi.Driver
	irqline  exti.Lines
	irqflag  rtos.EventFlag
	pdn, csn gpio.Pin
	started  bool
}

func NewSPI(spidrv *spi.Driver, csn, pdn gpio.Pin, irqline exti.Lines) *SPI {
	spidrv.P.SetConf(
		spi.Master | spi.MSBF | spi.CPOL0 | spi.CPHA0 |
			spidrv.P.BR(11e6) | // 11 MHz max. before configure PCLK.
			spi.SoftSS | spi.ISSHigh,
	)
	spidrv.P.SetWordSize(8)
	spidrv.P.Enable()
	csn.Set()
	irqline.EnableFallTrig()
	irqline.EnableIRQ()
	dci := new(SPI)
	dci.spi = spidrv
	dci.csn = csn
	dci.pdn = pdn
	dci.irqline = irqline
	return dci
}

func (dci *SPI) SPI() *spi.Driver {
	return dci.spi
}

func (dci *SPI) PDN() gpio.Pin {
	return dci.pdn
}

func (dci *SPI) IRQL() exti.Lines {
	return dci.irqline
}

func (dci *SPI) IRQF() *rtos.EventFlag {
	return &dci.irqflag
}

func (dci *SPI) ISR() {
	dci.irqflag.Signal(1)
}

func (dci *SPI) Err(clear bool) error {
	return dci.spi.Err(clear)
}

func (dci *SPI) End() {
	dci.started = false
	dci.csn.Set()
}

func (dci *SPI) Read(s []byte) {
	dci.spi.WriteRead(nil, s)
}

func (dci *SPI) Write(s []byte) {
	if !dci.started {
		if len(s) == 0 {
			return
		}
		dci.started = true
		dci.csn.Clear()
	}
	dci.spi.WriteRead(s, nil)
}
