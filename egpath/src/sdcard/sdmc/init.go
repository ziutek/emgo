package sdmc

import (
	"delay"
	"encoding/binary/be"

	"sdcard"
)

// Init initializes the card and reports wheter the initialization was
// successful.
func (c *Card) Init(freqhz int, bw sdcard.BusWidth, ocr sdcard.OCR) (sdcard.CID, error) {
	c.status = 0
	h := c.host

	// Set initial clock and bus width.
	h.SetClock(400e3)
	hostBusWidths := h.SetBusWidth(sdcard.Bus1)

	// BUG: Use hostBusWidths to detect SPI host.
	// BUG: SPI host not supported.

	// SD card power-up takes maximum of 1 ms or 74 SDIO_CK cycles.
	delay.Millisec(1)

	// Reset.
	h.SendCmd(sdcard.CMD0())

	// CMD0 may require up to 8 SDIO_CK cycles to reset the card.
	delay.Millisec(1)

	if ocr&sdcard.HCXC != 0 {
		// Verify card interface operating condition.
		// BUG: LVR not supported
		vhs, pattern := h.SendCmd(sdcard.CMD8(sdcard.V27_36, 0xAC)).R7()
		if err := h.Err(true); err != nil {
			if err != sdcard.ErrCmdTimeout {
				h.SetClock(0)
				return sdcard.CID{}, err
			}
			ocr &^= sdcard.HCXC
		} else if vhs != sdcard.V27_36 || pattern != 0xAC {
			h.SetClock(0)
			return sdcard.CID{}, ErrInitIC
		}
	}

	// Initializing SD card.
	for i := 0; i < 20; i++ {
		h.SendCmd(sdcard.CMD55(0))
		c.oca = h.SendCmd(sdcard.ACMD41(ocr)).R3()
		if err := c.checkErr(); err != nil {
			return sdcard.CID{}, err
		}
		if c.oca&sdcard.PWUP != 0 {
			break
		}
		delay.Millisec(50)
	}
	if c.oca&sdcard.PWUP == 0 {
		return sdcard.CID{}, ErrInitOC
	}

	// Read Card Identification Register.
	cid := h.SendCmd(sdcard.CMD2()).R2CID()

	// Generate new Relative Card Address.
	c.rca, _ = h.SendCmd(sdcard.CMD3()).R6()
	if err := c.checkErr(); err != nil {
		return cid, err
	}

	// After CMD3 card is in Data Transfer Mode (Standby State) and SDIO_CK can
	// be set to no more than 25 MHz (max. push-pull freq).
	f := freqhz
	if f > 25e6 {
		f = 25e6
	}
	h.SetClock(f)

	// Select card (put into Transfer State).
	c.statusCmd(sdcard.CMD7(c.rca))
	if err := c.checkErr(); err != nil {
		return cid, err
	}

	var buf [8]uint64
	bytes := sdcard.AsData(buf[:]).Bytes()

	// Read SD Configuration Register.
	h.SendCmd(sdcard.CMD55(c.rca))
	h.SetupData(sdcard.Recv|sdcard.Block8, buf[:1])
	c.statusCmd(sdcard.ACMD51())
	if err := c.checkErr(); err != nil {
		return cid, err
	}

	// Disable 50k pull-up resistor on D3/CD.
	h.SendCmd(sdcard.CMD55(c.rca))
	c.statusCmd(sdcard.ACMD42(false))
	if err := c.checkErr(); err != nil {
		return cid, err
	}

	scr := sdcard.SCR(be.Decode64(bytes))

	if bw == sdcard.Bus4 && hostBusWidths&(1<<bw) != 0 &&
		scr.SD_BUS_WIDTHS()&sdcard.SDBus4 != 0 {

		// Enable 4-bit data bus
		h.SendCmd(sdcard.CMD55(c.rca))
		c.statusCmd(sdcard.ACMD6(sdcard.Bus4))
		if err := c.checkErr(); err != nil {
			return cid, err
		}
		h.SetBusWidth(sdcard.Bus4)
	}

	if freqhz > 25e6 && scr.SD_SPEC() > 0 {
		// Check and enable High Speed.
		h.SetupData(sdcard.Recv|sdcard.Block64, buf[:])
		c.statusCmd(sdcard.CMD6(sdcard.ModeSwitch | sdcard.HighSpeed))
		if err := c.checkErr(); err != nil {
			return cid, err
		}
		sel := sdcard.SwitchFunc(be.Decode32(bytes[13:17]) & 0xFFFFFF)
		if sel&sdcard.AccessMode == sdcard.HighSpeed {
			delay.Millisec(1) // Function switch takes max. 8 SDIO_CK cycles.
			h.SetClock(freqhz)
		}
	}

	// Set block size to 512 B (required for protocol version < 2 or SDSC card).
	if c.oca&sdcard.HCXC == 0 {
		c.statusCmd(sdcard.CMD16(512))
		if err := c.checkErr(); err != nil {
			return cid, err
		}
	}

	return cid, nil
}
