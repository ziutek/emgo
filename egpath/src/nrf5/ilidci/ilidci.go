package ilidci

import (
	"nrf5/hal/gpio"
	"nrf5/hal/spi"
)

// DCI implements ili9341.DCI interface.
type DCI struct {
	spi  *spi.Driver
	dc   gpio.Pin
	freq spi.Freq
}

// Make returns initialized DCI struct. It only initializes DCI type and does
// not configure SPI peripheral. Use DCI.Setup to configure SPI.
func Make(spidrv *spi.Driver, dc gpio.Pin, baudrate spi.Freq) DCI {
	return DCI{spidrv, dc, baudrate}
}

// New provides convenient way to create heap allocated DCI. See Make.
func New(spidrv *spi.Driver, dc gpio.Pin, baudrate spi.Freq) *DCI {
	dci := new(DCI)
	*dci = Make(spidrv, dc, baudrate)
	return dci
}

// Setup configures and enables SPI. Use it after reset and after any
// other use of the same SPI peripheral.
func (dci *DCI) Setup() {
	drv := dci.spi
	drv.Disable()
	drv.P.StoreFREQUENCY(dci.freq)
	drv.Enable()
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

func (dci *DCI) Err(clear bool) error {
	return nil
}
