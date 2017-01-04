package ilidci

import (
	"rtos"

	"stm32/hal/gpio"
	"stm32/hal/spi"
)

// DCI implements ili9341.DCI interface.
type DCI struct {
	spi *spi.Driver
	dc  gpio.Pin
}

func NewDCI(spidrv *spi.Driver, dc gpio.Pin) *DCI {
	dci := new(DCI)
	dci.spi = spidrv
	dci.dc = dc
	return dci
}

func (dci *DCI) SPI() *spi.Driver {
	return dci.spi
}

func (dci *DCI) SetWordSize(size int) {
	dci.spi.P.SetWordSize(size)
}

func (dci *DCI) waitBusy() {
	for {
		if ev, _ := dci.spi.P.Status(); ev&spi.Busy == 0 {
			break
		}
		rtos.SchedYield()
	}
}

func (dci *DCI) Cmd(b byte) {
	dci.waitBusy()
	dci.dc.Clear()
	dci.spi.WriteReadByte(b)
	dci.dc.Set()
}

func (dci *DCI) Byte(b byte) {
	dci.spi.WriteReadByte(b)
}

func (dci *DCI) Cmd16(w uint16) {
	dci.waitBusy()
	dci.dc.Clear()
	dci.spi.WriteReadWord16(w)
	dci.dc.Set()
}

func (dci *DCI) Word(w uint16) {
	dci.spi.WriteReadWord16(w)
}

func (dci *DCI) Data(data []uint16) {
	dci.spi.WriteRead16(data, nil)
}

func (dci *DCI) Fill(w uint16, n int) {
	dci.spi.RepeatWord16(w, n)
}

func (dci *DCI) Err() error {
	return dci.spi.Err()
}
