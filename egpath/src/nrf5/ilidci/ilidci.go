package ilidci

import (
	"nrf5/hal/gpio"
	"nrf5/hal/spi"
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

func (dci *DCI) SetWordSize(_ int) {
}

func (dci *DCI) Cmd(b byte) {
	dci.dc.Clear()
	dci.spi.WriteReadByte(b)
	dci.dc.Set()
}

func (dci *DCI) WriteByte(b byte) {
	dci.spi.WriteReadByte(b)
}

func (dci *DCI) Cmd2(w uint16) {
	dci.dc.Clear()
	dci.spi.WriteReadWord16(w)
	dci.dc.Set()
}

func (dci *DCI) WriteWord(w uint16) {
	dci.spi.WriteReadWord16(w)
}

func (dci *DCI) Write(data []uint16) {
	dci.spi.WriteRead16(data, nil)
}

func (dci *DCI) Fill(w uint16, n int) {
	dci.spi.RepeatWord16(w, n)
}

func (dci *DCI) Err() error {
	return nil
}
