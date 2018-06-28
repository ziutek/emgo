package sdmc

import (
	"delay"
	"encoding/binary/be"
	"errors"

	"sdcard"
)

var (
	ErrInitV2 = errors.New("sdmc: init V2")
	ErrInitOC = errors.New("sdmc: init OC")
	ErrStatus = errors.New("sdmc: status")
)

type Card struct {
	host   sdcard.Host
	oca    sdcard.OCR
	status sdcard.CardStatus
	rca    uint16
}

func MakeCard(host sdcard.Host) Card {
	return Card{host: host}
}

func NewCard(host sdcard.Host) *Card {
	c := new(Card)
	*c = MakeCard(host)
	return c
}

// Status returns last received card status.
func (c *Card) Status() sdcard.CardStatus {
	return c.status
}

func (c *Card) statusCmd(cmd sdcard.Command, arg uint32) {
	c.status = c.host.SendCmd(cmd, arg).R1()
}

func (c *Card) checkErr() error {
	if err := c.host.Err(true); err != nil {
		c.host.SetClock(0)
		return err
	}
	errFlags := sdcard.OUT_OF_RANGE |
		sdcard.ADDRESS_ERROR |
		sdcard.BLOCK_LEN_ERROR |
		sdcard.ERASE_SEQ_ERROR |
		sdcard.ERASE_PARAM |
		sdcard.WP_VIOLATION |
		sdcard.LOCK_UNLOCK_FAILED |
		sdcard.COM_CRC_ERROR |
		sdcard.ILLEGAL_COMMAND |
		sdcard.CARD_ECC_FAILED |
		sdcard.CC_ERROR |
		sdcard.ERROR |
		sdcard.CSD_OVERWRITE |
		sdcard.WP_ERASE_SKIP |
		sdcard.CARD_ECC_DISABLED |
		sdcard.ERASE_RESET |
		sdcard.AKE_SEQ_ERROR
	if c.status&errFlags != 0 {
		return ErrStatus
	}
	return nil
}

// Init initializes the card and reports wheter the initialization was
// successful.
func (c *Card) Init(freqhz int, ocr sdcard.OCR) error {
	c.status = 0
	h := c.host

	// Set initial clock and bus width.
	h.SetClock(400e3)
	h.SetBusWidth(sdcard.Bus1)

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
				return err
			}
			ocr &^= sdcard.HCXC
		} else if vhs != sdcard.V27_36 || pattern != 0xAC {
			h.SetClock(0)
			return ErrInitV2
		}
	}

	// Initializing SD card.
	for i := 0; c.oca&sdcard.PWUP == 0 && i < 20; i++ {
		h.SendCmd(sdcard.CMD55(0))
		c.oca = h.SendCmd(sdcard.ACMD41(ocr)).R3()
		if err := c.checkErr(); err != nil {
			return err
		}
		delay.Millisec(50)
	}
	if c.oca&sdcard.PWUP == 0 {
		return ErrInitOC
	}

	// Generate new Relative Card Address.
	c.rca, _ = h.SendCmd(sdcard.CMD3()).R6()
	if err := c.checkErr(); err != nil {
		return err
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
		return err
	}

	var data [8]uint64
	buf := sdcard.Data(data[:])

	// Read SD Configuration Register.
	h.SendCmd(sdcard.CMD55(c.rca))
	h.SetupData(sdcard.Recv|sdcard.Block8, buf[:1])
	c.statusCmd(sdcard.ACMD51())
	if err := c.checkErr(); err != nil {
		return err
	}

	// Disable 50k pull-up resistor on D3/CD.
	h.SendCmd(sdcard.CMD55(c.rca))
	c.statusCmd(sdcard.ACMD42(false))
	if err := c.checkErr(); err != nil {
		return err
	}

	scr := sdcard.SCR(be.Decode64(buf.Bytes()))

	if scr.SD_BUS_WIDTHS()&sdcard.SDBus4 != 0 {
		// Enable 4-bit data bus
		h.SendCmd(sdcard.CMD55(c.rca))
		c.statusCmd(sdcard.ACMD6(sdcard.Bus4))
		if err := c.checkErr(); err != nil {
			return err
		}
		h.SetBusWidth(sdcard.Bus4)
	}

	if freqhz > 25e6 && scr.SD_SPEC() > 0 {
		// Check and enable High Speed.
		h.SetupData(sdcard.Recv|sdcard.Block64, buf[:])
		c.statusCmd(sdcard.CMD6(sdcard.ModeSwitch | sdcard.HighSpeed))
		if err := c.checkErr(); err != nil {
			return err
		}
		sel := sdcard.SwitchFunc(be.Decode32(buf.Bytes()[13:17]) & 0xFFFFFF)
		if sel&sdcard.AccessMode == sdcard.HighSpeed {
			delay.Millisec(1) // Function switch takes max. 8 SDIO_CK cycles.
			h.SetClock(freqhz)
		}
	}

	// Set block size to 512 B (required for protocol version < 2 or SDSC card).
	if c.oca&sdcard.HCXC == 0 {
		c.statusCmd(sdcard.CMD16(512))
		if err := c.checkErr(); err != nil {
			return err
		}
	}

	return nil
}
