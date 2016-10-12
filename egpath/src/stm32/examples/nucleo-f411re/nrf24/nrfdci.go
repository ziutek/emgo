package main

import (
	"mmio"
	"rtos"

	"arch/cortexm/bitband"

	"stm32/hal/exti"
	"stm32/hal/spi"
	"stm32/hal/system"

	"stm32/hal/raw/tim"
)

type NRFDCI struct {
	spi  *spi.Driver
	csn  bitband.Bit
	cet  *tim.TIM_Periph
	ocmn uint
	irq  exti.Lines
	flag rtos.EventFlag
}

func NewNRFDCI(spidrv *spi.Driver, csn bitband.Bit, pclk uint, cet *tim.TIM_Periph, ch int, irqline exti.Lines) *NRFDCI {
	dci := new(NRFDCI)
	dci.spi = spidrv
	spidrv.P.SetConf(
		spi.Master | spi.MSBF | spi.CPOL0 | spi.CPHA0 |
			spidrv.P.BR(10e6) | // 10 MHz max.
			spi.SoftSS | spi.ISSHigh,
	)
	spidrv.P.Enable()

	dci.csn = csn
	csn.Set()

	dci.cet = cet
	if pclk != system.AHB.Clock() {
		pclk *= 2
	}

	// PSC=1 gives shortest possible delay (1/pclk) before CE pulse.
	cet.PSC.U16.Store(1)

	// ARR=(pclk+1e5-1)/1e5 corresponds to the shortest posible pulse but
	// not less than 10 us. CE will be asserted after 1/pclk for ARR/pclk.
	cet.ARR.U32.Store(uint32(pclk+1e5-1) / 1e5)

	// Reset CNT and transfer PSC, ARR to the corresponding shadow registers.
	cet.EGR.Store(tim.UG)

	var (
		ccr *mmio.U32
		cce tim.CCER_Bits
	)
	switch ch {
	case 1:
		ccr = &cet.CCR1.U32
		cce = tim.CC1E
		dci.ocmn = tim.OC1Mn
	case 2:
		ccr = &cet.CCR2.U32
		cce = tim.CC2E
		dci.ocmn = tim.OC2Mn
	case 3:
		ccr = &cet.CCR3.U32
		cce = tim.CC3E
		dci.ocmn = tim.OC3Mn
	case 4:
		ccr = &cet.CCR4.U32
		cce = tim.CC4E
		dci.ocmn = tim.OC4Mn
	}
	ccr.Store(200)
	cet.CCER.Store(cce)

	dci.irq = irqline
	irqline.EnableFallTrig()
	irqline.EnableInt()
	return dci
}

func (dci *NRFDCI) WriteRead(oi ...[]byte) (n int, err error) {
	dci.csn.Clear()
	dci.spi.WriteReadMany(oi...)
	dci.csn.Set()
	return n, dci.spi.Err()
}

func (dci *NRFDCI) SetCE(v int) error {
	switch v {
	case 0:
		dci.cet.CCMR2.Store(4 << dci.ocmn)
	case 1:
		dci.cet.CCMR2.Store(5 << dci.ocmn)
	case 2:
		dci.cet.CCMR2.Store(7 << dci.ocmn)
		dci.cet.CR1.Store(tim.OPM | tim.CEN)
	}
	return nil
}

func (dci *NRFDCI) IRQ() exti.Lines {
	return dci.irq
}

func (dci *NRFDCI) ISR() {
	dci.flag.Set()
}

func (dci *NRFDCI) Baudrate() uint {
	return dci.spi.P.Baudrate(dci.spi.P.Conf())
}

func (dci *NRFDCI) Wait(deadline int64) bool {
	if dci.flag.Wait(deadline) {
		dci.flag.Clear()
		return true
	}
	return false
}
