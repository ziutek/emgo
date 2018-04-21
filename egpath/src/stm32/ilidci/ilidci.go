package ilidci

import (
	"stm32/hal/gpio"
	"stm32/hal/spi"
)

// DCI implements ili9341.DCI interface.
type DCI struct {
	spi  *spi.Driver
	dc   gpio.Pin
	brws uint // Bit 0 describes word size. Other bits describe baudrate.
}

// Make returns initialized DCI struct. It only initializes DCI type and does
// not configure SPI peripheral. Use DCI.Setup to configure SPI.
func Make(spidrv *spi.Driver, dc gpio.Pin, baudrate int) DCI {
	return DCI{spidrv, dc, uint(baudrate) << 1}
}

// New provides convenient way to create heap allocated DCI. See Make.
func New(spidrv *spi.Driver, dc gpio.Pin, baudrate int) *DCI {
	dci := new(DCI)
	*dci = Make(spidrv, dc, baudrate)
	return dci
}

// Setup configures and enables SPI peripheral.
func (dci *DCI) Setup() {
	p := dci.spi.Periph()
	p.EnableClock(true)
	p.Disable()
	p.SetConf(
		spi.Master | spi.MSBF | spi.CPOL0 | spi.CPHA0 |
			p.BR(int(dci.brws>>1)) |
			spi.SoftSS | spi.ISSHigh,
	)
	p.SetWordSize(8 * int(1+dci.brws&1))
	p.Enable()
}

func (dci *DCI) SetWordSize(size int) {
	dci.spi.Periph().SetWordSize(size)
	dci.brws = dci.brws&^1 | uint(size/8)&1
}

func (dci *DCI) SPI() *spi.Driver {
	return dci.spi
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
	return dci.spi.Err(clear)
}
