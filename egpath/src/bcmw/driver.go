package bcmw

import (
	"delay"
	"io"

	"sdcard"
	"sdcard/sdio"
)

type Error byte

const (
	ErrTimeout Error = iota + 1
	ErrIOStatus
	ErrUnknownChip
	ErrARMIsDown
)

//emgo:const
var errStr = [...]string{
	ErrTimeout:     "bcmw: timeout",
	ErrIOStatus:    "bcmw: IO status",
	ErrUnknownChip: "bcmw: unknown chip",
	ErrARMIsDown:   "bcmw: ARM is down",
}

func (e Error) Error() string {
	return errStr[e]
}

type Driver struct {
	sd              sdcard.Host
	backplaneWindow uint32
	ramSize         uint32
	chipID          uint16
	ioStatus        sdcard.IOStatus
	err             Error

	irq bool
}

func MakeDriver(sd sdcard.Host) Driver {
	return Driver{sd: sd}
}

func NewDriver(sd sdcard.Host) *Driver {
	d := new(Driver)
	*d = MakeDriver(sd)
	return d
}

func (d *Driver) IOStatus() sdcard.IOStatus {
	return d.ioStatus
}

func (d *Driver) ChipID() uint16 {
	return d.chipID
}

func (d *Driver) ISR() {
	d.irq = true
}

func (d *Driver) Init(reset func(nrst int), oobIntPin int) error {
	d.ramSize = 0
	d.chipID = 0
	d.ioStatus = 0
	d.err = 0
	d.sd.Err(true)

	reset(0)
	sd := d.sd
	sd.SetBusWidth(sdcard.Bus4)
	sd.SetClock(400e3, true)
	delay.Millisec(1)
	reset(1)
	delay.Millisec(1)

	// Enumerate and put the card into Transfer State.

	for retry := 250; ; retry-- {
		delay.Millisec(2)
		sd.SendCmd(sdcard.CMD0())
		sd.SendCmd(sdcard.CMD5(0))
		rca, _ := sd.SendCmd(sdcard.CMD3()).R6()
		if sd.Err(true) == nil {
			sd.SendCmd(sdcard.CMD7(rca)) // Select the card.
			break
		}
		if retry == 1 {
			d.err = ErrTimeout
			return d.firstErr()
		}
	}

	// Enable 4-bit data bus.

	r8 := d.sdioRead8(cia, sdio.CCCR_BUSICTRL)
	d.sdioWrite8(cia, sdio.CCCR_BUSICTRL, r8&^3|byte(sdcard.Bus4))

	// Enable High Speed if supported.

	r8 = d.sdioRead8(cia, sdio.CCCR_SPEEDSEL)
	if false && r8&1 != 0 {
		d.sdioWrite8(cia, sdio.CCCR_SPEEDSEL, r8|2)
		sd.SetClock(50e6, false)
	} else {
		sd.SetClock(25e6, false)
	}

	// Enable function 1.

	d.sdioEnableFunc(backplane, 500)
	d.sdioSetBlockSize(backplane, 64)

	// Disable extra SDIO pull-ups.

	d.sdioWrite8(backplane, sbsdioFunc1SDIOPullUp, 0)

	// Enable OOB IRQ, active high.

	if oobIntPin >= 0 {
		d.sdioWrite8(cia, cccrSepIntCtl, sepIntCtlMask|sepIntCtlEn|sepIntCtlPol)
		switch oobIntPin {
		case 0:
			// Default GPIO_0 pin.
		case 1:
			// GPIO_1 pin.
			d.sdioWrite8(backplane, sbsdioGPIOSel, 0xF)
			d.sdioWrite8(backplane, sbsdioGPIOOut, 0)
			d.sdioWrite8(backplane, sbsdioGPIOEn, 2)
			d.backplaneWrite32(commonGPIOCtl, 2)
		default:
			panic("bcmw: bad IRQ pin")
		}
	}

	// Enable Active Low-Power clock.

	d.sdioWrite8(
		backplane, sbsdioFunc1ChipClkCSR,
		sbsdioForceHwClkReqOff|sbsdioALPAvailReq|sbsdioForceALP,
	)
	for retry := 8; ; retry-- {
		delay.Millisec(2)
		r8 = d.sdioRead8(backplane, sbsdioFunc1ChipClkCSR)
		if d.error() {
			return d.firstErr()
		}
		if r8&sbsdioALPAvail != 0 {
			break
		}
		if retry == 1 {
			d.err = ErrTimeout
			return d.err
		}
	}
	d.sdioWrite8(backplane, sbsdioFunc1ChipClkCSR, 0) // Clear ALP request.

	// Identify chip.

	r32 := d.backplaneRead32(commonEnumBase)
	d.chipID = uint16(r32)
	if r32>>28&0xF != 1 {
		d.err = ErrUnknownChip // Not AXI.
		return d.firstErr()
	}
	switch d.chipID {
	case 43362:
		d.ramSize = 240 * 1024
	case 43430:
		d.ramSize = 512 * 1024
	default:
		d.err = ErrUnknownChip
	}
	return d.firstErr()
}

func (d *Driver) UploadFirmware(firmware, nvram io.Reader, nvramSiz int) error {
	if d.error() {
		return d.firstErr()
	}

	// Disable ARMCM3 core and reset SOCSRAM.

	d.chipCoreDisable(coreARMCM3, 0, 0)
	d.chipCoreReset(coreSOCSRAM, 0, 0, 0)

	// Upload firmware.

	if d.chipID == 43430 {
		// Disable remap for SRAM3 in case of 4343x
		d.backplaneWrite32(socsramBankxIndex, 3)
		d.backplaneWrite32(socsramBankxPDA, 0)
	}
	if err := d.backplaneUpload(0, firmware); err != nil {
		return err
	}

	// Upload NVRAM.

	siz := uint32(nvramSiz+63) &^ 63
	if err := d.backplaneUpload(d.ramSize-4-siz, nvram); err != nil {
		return err
	}
	token := uint32(siz) >> 2
	token = ^token<<16 | token&0xFFFF
	d.backplaneWrite32(d.ramSize-4, token)

	// Reset ARMCM3 core.

	d.chipCoreReset(coreARMCM3, 0, 0, 0)
	up := d.chipIsCoreUp(coreARMCM3)
	if d.error() {
		return d.firstErr()
	}
	if !up {
		d.err = ErrARMIsDown
		return d.err
	}

	// Wait for High Throughput clock.

	for retry := 250; ; retry-- {
		r := d.sdioRead8(backplane, sbsdioFunc1ChipClkCSR)
		if d.error() {
			return d.firstErr()
		}
		if r&sbsdioHTAvail != 0 {
			break
		}
		if retry == 1 {
			d.err = ErrARMIsDown
			return d.err
		}
		delay.Millisec(2)
	}

	if d.chipID == 43430 {
		// TODO: Check save/restore support in firmware and configure it.
	}

	// Enable function 2.

	d.sdioEnableFunc(wlanData, 500)
	d.sdioSetBlockSize(wlanData, 512)

	// Enable interrupts from function 2 (WLAN data).

	d.backplaneWrite32(sdiodHostIntMask, intHMBSWMask)
	d.backplaneWrite32(sdiodFuncIntMask, 1<<wlanData)
	d.sdioWrite8(cia, sdio.CCCR_INTEN, 1<<cia|1<<wlanData)

	// Lower F2 watermark to avoid DMA hang in F2 when SD clock is stopped.

	d.sdioWrite8(backplane, sbsdioWatermark, 8)

	return d.firstErr()
}

func (d *Driver) StatusLoop(oobIRQ func() int) {
	for {
		d.debug(
			"%t OOB IRQ: %d, CCCR_INTPEND: %bb",
			d.irq, oobIRQ(), d.sdioRead8(cia, sdio.CCCR_INTPEND),
		)
		irqs := d.backplaneRead32(sdiodIntStatus)
		d.debug(" sdiodIntStatus: %bb\n", irqs)
		if irqs&intHMBSWMask != 0 {
			d.backplaneWrite32(sdiodIntStatus, irqs&intHMBSWMask)
		}
		if irqs&intHMBFrame == 0 {
			delay.Millisec(500)
			continue
		}
		for d.sdpcmReadFrame() {
		}
	}
}
