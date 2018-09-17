package bcmw

import (
	"delay"
	"errors"
	"fmt"
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

func (d *Driver) Init(reset func(nrst int), oobIntPin int) {
	if d.error() {
		return
	}
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
			d.timeout = true
			return
		}
	}

	// Enable function 1.

	d.enableFunction(backplane)

	// Enable 4-bit data bus.

	r := d.cmd52(cia, sdio.CCCR_BUSICTRL, sdcard.Read, 0)
	d.cmd52(cia, sdio.CCCR_BUSICTRL, sdcard.Write, r&^3|byte(sdcard.Bus4))

	// Set block size to 64 bytes for all functions.

	for retry := 250; ; retry-- {
		fmt.Printf("blk siz\n")
		delay.Millisec(2)
		r := d.cmd52(cia, sdio.CCCR_BLKSIZE0, sdcard.WriteRead, 64)
		if d.error() {
			return // Fast return if error.
		}
		if r == 64 {
			break
		}
		if retry == 1 {
			d.timeout = true
			return
		}
	}
	for f := cia; f <= wlanData; f++ {
		d.cmd52(cia, f<<8+sdio.FBR_BLKSIZE0, sdcard.Write, 64)
		d.cmd52(cia, f<<8+sdio.FBR_BLKSIZE1, sdcard.Write, 0)
	}

	// Enable interrupts from Backplane and WLAN Data functions (1<<cia is
	// Master Interrupt Enable bit).

	d.cmd52(cia, sdio.CCCR_INTEN, sdcard.Write, 1<<cia|1<<backplane|1<<wlanData)

	// Enable High Speed if supported.

	r = d.cmd52(cia, sdio.CCCR_SPEEDSEL, sdcard.Read, 0)
	if false && r&1 != 0 {
		d.cmd52(cia, sdio.CCCR_SPEEDSEL, sdcard.Write, r|2)
		sd.SetClock(50e6, false)
	} else {
		sd.SetClock(25e6, false)
	}

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

	// Enable Active Low-Power clock.

	d.cmd52(
		backplane, sbsdioFunc1ChipClkCSR, sdcard.Write,
		sbsdioForceHwClkReqOff|sbsdioALPAvailReq|sbsdioForceALP,
	)
	for retry := 50; ; retry-- {
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
		delay.Millisec(2)
	}
	// Clear the enable request.
	d.cmd52(backplane, sbsdioFunc1ChipClkCSR, sdcard.Write, 0)

	// Disable pull-ups - we use STM32 GPIO pull-ups.

	d.cmd52(backplane, sbsdioFunc1SDIOPullUp, sdcard.Write, 0)

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
}

func (d *Driver) UploadFirmware(r io.Reader, firmware []uint64) {
	if d.error() {
		return
	}
	d.disableCore(coreARMCM3)
	d.resetCore(coreSOCSRAM)

	if d.chip == &chip43438 {
		// Disable remap for SRAM3 in case of 4343x
		d.wbr32(socsramBankxIndex, 3)
		d.wbr32(socsramBankxPDA, 0)
	}
	d.wbb(0, firmware)
}

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
