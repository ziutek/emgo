package evedci

import (
	"stm32/hal/exti"
	"stm32/hal/gpio"
	"stm32/hal/spi"
)

// SPI implements eve.DCI using Serial Peripheral Interface.
type SPI struct {
	spi      *spi.Driver
	irqline  exti.Lines
	irqchan  chan struct{}
	pdn, csn gpio.Pin
	started  bool
}

func NewSPI(spidrv *spi.Driver, csn, pdn gpio.Pin, irqline exti.Lines) *SPI {
	p := spidrv.P
	p.SetConf(
		spi.Master | spi.MSBF | spi.CPOL0 | spi.CPHA0 |
			p.BR(11e6) | // 11 MHz max. before configure PCLK.
			spi.SoftSS | spi.ISSHigh,
	)
	p.SetWordSize(8)
	p.Enable()
	csn.Set()
	irqline.EnableFallTrig()
	irqline.EnableIRQ()
	dci := new(SPI)
	dci.spi = spidrv
	dci.csn = csn
	dci.pdn = pdn
	dci.irqline = irqline
	dci.irqchan = make(chan struct{}, 1)
	return dci
}

func (dci *SPI) SetBaudrate(baud int) {
	p := dci.spi.P
	p.SetConf(p.Conf()&^spi.BR256 | p.BR(30e6))
}

func (dci *SPI) SPI() *spi.Driver {
	return dci.spi
}

func (dci *SPI) SetPDN(pdn int) {
	dci.pdn.Store(pdn)
}

func (dci *SPI) IRQ() <-chan struct{} {
	return dci.irqchan
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

func (dci *SPI) ISR() {
	select {
	case dci.irqchan <- struct{}{}:
	default:
	}
}

func (dci *SPI) EXTI() exti.Lines {
	return dci.irqline
}
