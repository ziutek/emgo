package sdmmc

import (
	"rtos"
	"sync/fence"
	"unsafe"

	"sdcard"

	"stm32/hal/dma"
)

// Driver implements sdcard.Host interface.
type Driver struct {
	p      *Periph
	dma    *dma.Channel
	done   rtos.EventFlag
	err    Error
	dmaErr dma.Error
	data   sdcard.DataMode
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
	var err error
	switch {
	case d.err != 0:
		if d.err == ErrCmdTimeout {
			err = sdcard.ErrCmdTimeout
		} else {
			err = d.err
		}
	case d.dmaErr != 0:
		err = d.dmaErr
	default:
		goto end
	}
	if clear {
		d.err = 0
		d.dmaErr = 0
	}
end:
	return err
}

// SetFreq sets SDMMCCLK divider to provide SDMMC_CK frequency <= freqhz.
func (d *Driver) SetFreq(freqhz int, pwrsave bool) {
	var (
		clkdiv int
		cfg    BusClock
	)
	if freqhz > 0 {
		// BUG: This code assumes 48 MHz SDMMCCLK.
		cfg = ClkEna
		clkdiv = (48e6+freqhz-1)/freqhz - 2
	}
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
	if d.err != 0 || d.dmaErr != 0 {
		return
	}
	var waitFor Event
	if cmd&sdcard.HasResp != 0 {
		waitFor = CmdRespOK
	} else {
		waitFor = CmdSent
	}
	d.done.Reset(0)
	p := d.p
	p.Clear(EvAll, ErrAll)
	p.EnableIRQ(waitFor, ErrAll)
	p.SetArg(arg)
	fence.W() // This orders writes to normal and I/O memory.
	p.SetCmd(CmdEna | Command(cmd)&255)
	d.done.Wait(1, 0)
	_, d.err = p.Status()
	if cmd&sdcard.HasResp != 0 {
		if d.err&ErrCmdCRC != 0 {
			switch cmd & sdcard.RespType {
			case sdcard.R3, sdcard.R4:
				d.err &^= ErrCmdCRC
			}
			if r := cmd & sdcard.RespType; r == sdcard.R3 || r == sdcard.R4 {
				d.err &^= ErrCmdCRC
			}
		}
		if d.err != 0 {
			return
		}
		if cmd&sdcard.LongResp != 0 {
			resp[3] = p.Resp(0) // Most significant bits.
			resp[2] = p.Resp(1)
			resp[1] = p.Resp(2)
			resp[0] = p.Resp(3) // Least significant bits.
		} else {
			resp[0] = p.Resp(0)
		}
	}
	switch {
	case d.data == 0:
		return
	case d.data&sdcard.Stream == 0:
		waitFor = DataBlkEnd
	default:
		waitFor = DataEnd
	}
	d.data = 0
	d.done.Reset(0)
	p.EnableIRQ(waitFor, ErrAll)
	d.done.Wait(1, 0)
	_, d.err = p.Status()
	// Ensure DMA transfer has been completed (it should be).
	for {
		ev, err := d.dma.Status()
		if err != 0 {
			d.dmaErr = err
			break
		}
		if ev == dma.Complete {
			break
		}
	}
	return
}

func (d *Driver) Data(mode sdcard.DataMode, buf sdcard.Data) {
	if d.err != 0 || d.dmaErr != 0 {
		return
	}
	d.data = mode
	dir := dma.PTM
	if mode&sdcard.Recv == 0 {
		dir = dma.MTP
	}
	ch := d.dma
	ch.Clear(dma.EvAll, dma.ErrAll)
	ch.Setup(dir | dma.PFC | dma.IncM | dma.FT4 | dma.PB4 | dma.MB4)
	ch.SetWordSize(4, 4)
	ch.SetAddrP(unsafe.Pointer(&d.p.raw.FIFO))
	ch.SetAddrM(unsafe.Pointer(&buf[0]))
	ch.Enable()
	p := d.p
	p.SetDataTimeout(191 << 18) // ≈ 1s at high speed, const. fits in mov.w.
	p.SetDataLen(len(buf) * 8)
	p.SetDataCtrl(DTEna | UseDMA | DataCtrl(mode))
}
