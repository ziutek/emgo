package sdmmc

import (
	"rtos"
	"sync/fence"

	"sdcard"

	"stm32/hal/dma"
)

// Driver implements sdcard.Host interface.
type Driver struct {
	p    *Periph
	dma  *dma.Channel
	done rtos.EventFlag
	err  Error
}

// MakeDriver returns initialized SPI driver that uses provided SPI peripheral
// and DMA channel.
func MakeDriver(p *Periph, dma *dma.Channel) Driver {
	return Driver{p: p, dma: dma}
}

// NewDriver provides convenient way to create heap allocated Driver struct.
func NewDriver(p *Periph, dma *dma.Channel) *Driver {
	d := new(Driver)
	*d = MakeDriver(p, dma)
	return d
}

func (d *Driver) Periph() *Periph {
	return d.p
}

func (d *Driver) DMA() *dma.Channel {
	return d.dma
}

func (d *Driver) Err(clear bool) error {
	if d.err == 0 {
		return nil
	}
	var err error
	if d.err == ErrCmdTimeout {
		err = sdcard.ErrCmdTimeout
	} else {
		err = d.err
	}
	if clear {
		d.err = 0
	}
	return err
}

// SetFreq sets SDMMCCLK divider to provide SDMMC_CK frequency <= freqhz.
func (d *Driver) SetFreq(freqhz int, pwrsave bool) {
	// BUG: This code assumes 48 MHz SDMMCCLK.
	clkdiv := (48e6+freqhz-1)/freqhz - 2
	cfg := ClkEna
	if clkdiv < 0 {
		cfg |= ClkByp
	}
	if pwrsave {
		cfg |= PwrSave
	}
	d.p.SetBusClock(cfg, clkdiv)
}

func (d *Driver) ISR() {
	d.p.DisableIRQ(EvAll, ErrAll)
	d.done.Signal(1)
}

func (d *Driver) Cmd(cmd sdcard.Command, arg uint32) (resp sdcard.Response) {
	if d.err != 0 {
		return
	}
	var waitFor Event
	if cmd&sdcard.HasResp != 0 {
		waitFor = CmdRespOK
	} else {
		waitFor = CmdSent
	}
	d.done.Reset(0)
	d.p.Clear(EvAll, ErrAll)
	d.p.EnableIRQ(waitFor, ErrAll)
	d.p.SetArg(arg)
	fence.W() // This orders writes to normal and I/O memory.
	d.p.SetCmd(CmdEna | Command(cmd)&255)
	d.done.Wait(1, 0)
	_, d.err = d.p.Status()
	if cmd&sdcard.HasResp == 0 {
		return
	}
	if d.err != 0 {
		if d.err&ErrCmdCRC == 0 {
			return
		}
		if r := cmd & sdcard.RespType; r != sdcard.R3 && r != sdcard.R4 {
			return
		}
		// Ignore CRC error for R3, R4 responses.
		d.err &^= ErrCmdCRC
	}
	if cmd&sdcard.LongResp != 0 {
		resp[3] = d.p.Resp(0) // Most significant bits.
		resp[2] = d.p.Resp(1)
		resp[1] = d.p.Resp(2)
		resp[0] = d.p.Resp(3) // Least significant bits.
	} else {
		resp[0] = d.p.Resp(0)
	}
	return
}
