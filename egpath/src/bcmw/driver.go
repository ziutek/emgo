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
)

//emgo:const
var errStr = [...]string{
	ErrTimeout:     "bcmw: timeout",
	ErrIOStatus:    "bcmw: IO status",
	ErrUnknownChip: "bcmw: unknown chip",
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
}

func MakeDriver(sd sdcard.Host) Driver {
	return Driver{sd: sd}
}

func NewDriver(sd sdcard.Host) *Driver {
	d := new(Driver)
	*d = MakeDriver(sd)
	return d
}

func (d *Driver) firstErr() error {
	if err := d.sd.Err(false); err != nil {
		return err
	}
	if d.ioStatus&^sdcard.IO_CURRENT_STATE != 0 {
		return ErrIOStatus
	}
	if d.err != 0 {
		return d.err
	}
	return nil
}

func (d *Driver) IOStatus() sdcard.IOStatus {
	return d.IOStatus()
}

func (d *Driver) ChipID() uint16 {
	return d.chipID
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

	r8 := d.cmd52(cia, sdio.CCCR_BUSICTRL, sdcard.Read, 0)
	d.cmd52(cia, sdio.CCCR_BUSICTRL, sdcard.Write, r8&^3|byte(sdcard.Bus4))

	// Enable High Speed if supported.

	r8 = d.cmd52(cia, sdio.CCCR_SPEEDSEL, sdcard.Read, 0)
	if false && r8&1 != 0 {
		d.cmd52(cia, sdio.CCCR_SPEEDSEL, sdcard.Write, r8|2)
		sd.SetClock(50e6, false)
	} else {
		sd.SetClock(25e6, false)
	}

	// Set block size for all functions.

	d.sdioSetBlockSize(cia, 64)
	d.sdioSetBlockSize(backplane, 64)
	d.sdioSetBlockSize(wlanData, 64)

	// TODO: Enable interrupts.

	// Enable function 1.

	d.sdioEnableFunc(backplane, 500)

	// Enable Active Low-Power clock.

	d.sdiodWrite8(
		sbsdioFunc1ChipClkCSR,
		sbsdioForceHwClkReqOff|sbsdioALPAvailReq|sbsdioForceALP,
	)
	for retry := 8; ; retry-- {
		delay.Millisec(2)
		r8 = d.sdiodRead8(sbsdioFunc1ChipClkCSR)
		if d.error() {
			return d.firstErr() // Fast return if error.
		}
		if r8&sbsdioALPAvail != 0 {
			break
		}
		if retry == 1 {
			d.err = ErrTimeout
			return d.firstErr()
		}
	}
	d.sdiodWrite8(sbsdioFunc1ChipClkCSR, 0) // Clear ALP request.

	// Disable extra SDIO pull-ups.

	d.sdiodWrite8(sbsdioFunc1SDIOPullUp, 0)

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
	case 43438:
		d.ramSize = 512 * 1024
	default:
		d.err = ErrUnknownChip
	}
	return d.firstErr()
}

func (d *Driver) UploadFirmware(r io.Reader) error {
	if d.error() {
		return d.firstErr()
	}
	d.chipCoreDisable(coreARMCM3, 0, 0)
	d.chipCoreReset(coreSOCSRAM, 0, 0, 0)

	if d.chipID == 43438 {
		// Disable remap for SRAM3 in case of 4343x
		d.backplaneWrite32(socsramBankxIndex, 3)
		d.backplaneWrite32(socsramBankxPDA, 0)
	}

	delay.Millisec(50)

	var tmp [8]uint64
	buf := sdcard.AsData(tmp[:])
	for {
		n, err := io.ReadFull(r, buf.Bytes())
		_ = n
		if err != nil {
			//...
		}
	}
	return d.firstErr()
}

/*
	/////////////////////////////////////////////////////////

	// Chip must be passive before access its cores.

	d.chipCoreDisable(coreARMCM3, 0, 0)
	d.chipCoreReset(
		coreDot11MAC, ioCtlDot11PhyReset|ioCtlDot11PhyClockEn,
		ioCtlDot11PhyClockEn, ioCtlDot11PhyClockEn,
	)
	d.chipCoreReset(coreSOCSRAM, 0, 0, 0)

	if d.chipID == 43438 {
		// Disable remap for SRAM3 in case of 4343x
		d.backplaneWrite32(socsramBankxIndex, 3)
		d.backplaneWrite32(socsramBankxPDA, 0)
	}

	// Done with backplane-dependent accesses, disable clock.

	d.sdiodWrite8(sbsdioFunc1ChipClkCSR, 0)


		// Enable interrupts from Backplane and WLAN Data functions (1<<cia is
		// Master Interrupt Enable bit).

		d.cmd52(cia, sdio.CCCR_INTEN, sdcard.Write, 1<<cia|1<<backplane|1<<wlanData)


		// Wait till the backplane is ready.

		for retry := 250; ; retry-- {
			fmt.Printf("bkpl rdy\n")
			r = d.cmd52(cia, sdio.CCCR_IORDY, sdcard.Read, 0)
			if d.error() {
				return // Fast return if error.
			}
			if r&(1<<backplane) != 0 {
				break
			}
			if retry == 1 {
				d.timeout = true
				return
			}
			delay.Millisec(2)
		}

		// Enable function 2.

		d.enableFunction(wlanData)

		// Enable out-of-band interrupts.

		if oobIntPin >= 0 {
			d.cmd52(
				cia, cccrSepIntCtl, sdcard.Write,
				sepIntCtlMask|sepIntCtlEn|sepIntCtlPol, // Active high.
			)
			switch oobIntPin {
			case 0:
				// Default pin
			case 1:
				d.cmd52(backplane, sbsdioGPIOSel, sdcard.Write, 0xF)
				d.cmd52(backplane, sbsdioGPIOOut, sdcard.Write, 0)
				d.cmd52(backplane, sbsdioGPIOEn, sdcard.Write, 2)
				d.wbr32(commonGPIOCtl, 2)
			default:
				panic("bcmw: bad IRQ pin")
			}
		}

		// Disable Backplane interrupt

		d.cmd52(cia, sdio.CCCR_INTEN, sdcard.Write, 1<<cia|1<<wlanData)
*/

/*

func (d *Driver) UploadNVRAM(r io.Reader, nvram string) {
	if d.error() {
		return
	}
	var tmp [8]uint64
	buf := sdcard.AsData(tmp[:])
	nvsiz := (len(nvram) + 63) &^ 63 // Round up to n*64 bytes.
	addr := uint32(d.chip.ramSize - 4 - nvsiz)
	for len(nvram) > 0 {
		n := copy(buf.Bytes(), nvram)
		nvram = nvram[n:]
		d.wbb(addr, buf.Words())
		addr += 64
	}
	token := uint32(nvsiz) >> 2
	token = ^token<<16 | token
	d.wbr32(addr, token)

	d.resetCore(coreARMCM3)
	if d.isCoreUp(coreARMCM3) {
		fmt.Printf("ARM up!\n")
	} else {
		fmt.Printf("ARM down!\n")
		return
	}
	fmt.Printf("ht clk:")
	for retry := 250; ; retry-- {
		r := d.cmd52(backplane, sbsdioFunc1ChipClkCSR, sdcard.Read, 0)
		if d.error() {
			return // Fast return if error.
		}
		fmt.Printf(" %x", r)
		if r&sbsdioHTAvail != 0 {
			break
		}
		if retry == 1 {
			d.timeout = true
			return
		}
		delay.Millisec(2)
	}
}
*/
