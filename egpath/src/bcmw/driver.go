package bcmw

import (
	"delay"
	"errors"
	"io"

	"sdcard"
	"sdcard/sdio"
)

var ErrTimeout = errors.New("bcmw: timeout")

type Driver struct {
	sd              sdcard.Host
	chip            *Chip
	backplaneWindow uint32
}

func MakeDriver(sd sdcard.Host, chip *Chip) Driver {
	return Driver{sd: sd, chip: chip}
}

func NewDriver(sd sdcard.Host, chip *Chip) *Driver {
	d := new(Driver)
	*d = MakeDriver(sd, chip)
	return d
}

func (d *Driver) Init(reset func(nrst int)) error {
	sd := d.sd

	reset(0)
	sd.SetBusWidth(sdcard.Bus4)
	sd.SetClock(400e3, true)
	delay.Millisec(10)
	reset(1)
	delay.Millisec(10)

	// Enumerate.

	sd.SendCmd(sdcard.CMD0())
	sd.SendCmd(sdcard.CMD5(0))
	rca, _ := sd.SendCmd(sdcard.CMD3()).R6()
	if err := sd.Err(true); err != nil {
		return err
	}

	// Select the card and put it into Transfer State.

	sd.SendCmd(sdcard.CMD7(rca))
	if err := sd.Err(true); err != nil {
		return err
	}

	// Enable 4-bit data bus.

	r := cmd52(sd, cia, sdio.CCCR_BUSICTRL, sdcard.Read, 0)
	cmd52(sd, cia, sdio.CCCR_BUSICTRL, sdcard.Write, r&^3|byte(sdcard.Bus4))

	// Set block size to 64 bytes for all functions.

	for f := cia; f <= wlanData; f++ {
		cmd52(sd, cia, f<<8+sdio.FBR_BLKSIZE0, sdcard.Write, 64)
		cmd52(sd, cia, f<<8+sdio.FBR_BLKSIZE1, sdcard.Write, 0)
	}

	/*
		// Enable out-of-band interrupts.
		cmd52(
			sd, cia, SEP_INT_CTL, sdcard.Write,
			SEP_INTR_CTL_MASK|SEP_INTR_CTL_EN|SEP_INTR_CTL_POL,
		)
		// EMW3165 uses default IRQ pin (Pin0). Redirection isn't needed.
	*/

	// Enable interrupts from Backplane and WLAN Data functions (bit 0 is
	// Master Interrupt Enable bit).

	cmd52(sd, cia, sdio.CCCR_INTEN, sdcard.Write, 1|1<<backplane|1<<wlanData)

	// Enable High Speed if supported.

	r = cmd52(sd, cia, sdio.CCCR_SPEEDSEL, sdcard.Read, 0)
	if r&1 != 0 {
		cmd52(sd, cia, sdio.CCCR_SPEEDSEL, sdcard.Write, r|2)
		sd.SetClock(50e6, true)
	} else {
		sd.SetClock(25e6, true)
	}

	// Enable function 1.

	if enableFunction(sd, backplane) {
		return ErrTimeout
	}

	// Enable Active Low-Power clock.

	cmd52(
		sd, backplane, sbsdioFunc1ChipClkCSR, sdcard.Write,
		sbsdioForceHwClkReqOff|sbsdioALPAvailReq|sbsdioForceALP,
	)
	for retry := 50; ; retry-- {
		delay.Millisec(2)
		r := cmd52(sd, backplane, sbsdioFunc1ChipClkCSR, sdcard.Read, 0)
		if err := sd.Err(true); err != nil {
			return err
		}
		if r&sbsdioALPAvail != 0 {
			break
		}
		if retry == 1 {
			return ErrTimeout
		}
	}
	// Clear the enable request.
	cmd52(sd, backplane, sbsdioFunc1ChipClkCSR, sdcard.Write, 0)

	// Disable pull-ups - we use STM32 GPIO pull-ups.

	cmd52(sd, backplane, sbsdioFunc1SDIOPullUp, sdcard.Write, 0)

	return sd.Err(true)
}

func (d *Driver) disableCore(core int) {
	d.setBackplaneWindow(d.chip.baseAddr[core])
	sd := d.sd
	if cmd52(sd, backplane, ssbResetCtl, sdcard.Read, 0)&1 != 0 {
		return // Already in reset state.
	}
	delay.Millisec(10)
	cmd52(sd, backplane, ssbResetCtl, sdcard.Write, 1)
	delay.Millisec(1)
	cmd52(sd, backplane, ssbIOCtl, sdcard.Write, 0)
	cmd52(sd, backplane, ssbIOCtl, sdcard.Read, 0)
	delay.Millisec(1)
}

func (d *Driver) resetCore(core int) {
	d.disableCore(core)

	// Initialization sequence.

	cmd52(d.sd, backplane, ssbIOCtl, sdcard.Write, ioCtlClk|ioCtlFGC)

}

func (d *Driver) UploadFirmware(r io.Reader) error {
	d.disableCore(coreARMCM3)

	return nil
}
