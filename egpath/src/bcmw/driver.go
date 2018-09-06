package bcmw

import (
	"delay"
	"errors"
	"io"

	"sdcard"
	"sdcard/sdio"
)

var (
	ErrTimeout  = errors.New("bcmw: timeout")
	ErrIOStatus = errors.New("bcmw: IO status")
)

type Driver struct {
	sd              sdcard.Host
	chip            *Chip
	backplaneWindow uint32
	timeout         bool
	ioStatus        sdcard.IOStatus
}

func MakeDriver(sd sdcard.Host, chip *Chip) Driver {
	return Driver{sd: sd, chip: chip}
}

func NewDriver(sd sdcard.Host, chip *Chip) *Driver {
	d := new(Driver)
	*d = MakeDriver(sd, chip)
	return d
}

func (d *Driver) Err(clear bool) error {
	err := d.sd.Err(clear)
	switch {
	case err != nil:
	case d.ioStatus&^sdcard.IO_CURRENT_STATE != 0:
		err = ErrIOStatus
	case d.timeout:
		err = ErrTimeout
	default:
		return nil
	}
	if clear {
		d.ioStatus &= sdcard.IO_CURRENT_STATE
		d.timeout = false
	}
	return err
}

func (d *Driver) IOStatus() sdcard.IOStatus {
	return d.IOStatus()
}

func (d *Driver) Init(reset func(nrst int)) {
	if d.error() {
		return
	}
	reset(0)
	sd := d.sd
	sd.SetBusWidth(sdcard.Bus4)
	sd.SetClock(400e3, true)
	delay.Millisec(10)
	reset(1)
	delay.Millisec(10)

	// Enumerate.

	sd.SendCmd(sdcard.CMD0())
	sd.SendCmd(sdcard.CMD5(0))
	rca, _ := sd.SendCmd(sdcard.CMD3()).R6()

	// Select the card and put it into Transfer State.

	sd.SendCmd(sdcard.CMD7(rca))

	// Enable 4-bit data bus.

	r := d.cmd52(cia, sdio.CCCR_BUSICTRL, sdcard.Read, 0)
	d.cmd52(cia, sdio.CCCR_BUSICTRL, sdcard.Write, r&^3|byte(sdcard.Bus4))

	// Set block size to 64 bytes for all functions.

	for f := cia; f <= wlanData; f++ {
		d.cmd52(cia, f<<8+sdio.FBR_BLKSIZE0, sdcard.Write, 64)
		d.cmd52(cia, f<<8+sdio.FBR_BLKSIZE1, sdcard.Write, 0)
	}

	/*
		// Enable out-of-band interrupts.
		d.cmd52(
			cia, SEP_INT_CTL, sdcard.Write,
			SEP_INTR_CTL_MASK|SEP_INTR_CTL_EN|SEP_INTR_CTL_POL,
		)
		// EMW3165 uses default IRQ pin (Pin0). Redirection isn't needed.
	*/

	// Enable interrupts from Backplane and WLAN Data functions (1<<cia is
	// Master Interrupt Enable bit).

	d.cmd52(cia, sdio.CCCR_INTEN, sdcard.Write, 1<<cia|1<<backplane|1<<wlanData)

	// Enable High Speed if supported.

	r = d.cmd52(cia, sdio.CCCR_SPEEDSEL, sdcard.Read, 0)
	if r&1 != 0 {
		d.cmd52(cia, sdio.CCCR_SPEEDSEL, sdcard.Write, r|2)
		sd.SetClock(50e6, true)
	} else {
		sd.SetClock(25e6, true)
	}

	// Enable function 1.

	d.enableFunction(backplane)

	// Enable Active Low-Power clock.

	d.cmd52(
		backplane, sbsdioFunc1ChipClkCSR, sdcard.Write,
		sbsdioForceHwClkReqOff|sbsdioALPAvailReq|sbsdioForceALP,
	)
	for retry := 50; ; retry-- {
		delay.Millisec(2)
		r := d.cmd52(backplane, sbsdioFunc1ChipClkCSR, sdcard.Read, 0)
		if d.error() {
			return // Fast return if error.
		}
		if r&sbsdioALPAvail != 0 {
			break
		}
		if retry == 1 {
			d.timeout = true
			return
		}
	}
	// Clear the enable request.
	d.cmd52(backplane, sbsdioFunc1ChipClkCSR, sdcard.Write, 0)

	// Disable pull-ups - we use STM32 GPIO pull-ups.

	d.cmd52(backplane, sbsdioFunc1SDIOPullUp, sdcard.Write, 0)
}

func (d *Driver) UploadFirmware(r io.Reader) {
	if d.error() {
		return
	}
	d.disableCore(coreARMCM3)
	d.resetCore(coreSOCSRAM)
}
