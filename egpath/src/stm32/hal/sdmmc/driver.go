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

// SetBus sets data bus width and SDMMC_CK frequency <= freqhz.
func (d *Driver) SetBus(width sdcard.BusWidth, freqhz int, pwrsave bool) {
	var (
		clkdiv int
		cfg    BusClock
	)
	if width > sdcard.Bus8 {
		panic("sdcard: bus width")
	}
	if freqhz > 0 {
		// BUG: This code assumes 48 MHz SDMMCCLK.
		cfg = ClkEna | BusClock(width*3>>2)<<3
		clkdiv = (48e6+freqhz-1)/freqhz - 2
	}
	if clkdiv < 0 {
		cfg |= ClkByp
	}
	if pwrsave {
		cfg |= PwrSave
	}
	p := d.p
	p.SetBusClock(cfg, clkdiv)
	p.SetDataTimeout(uint(freqhz)) // ≈ 1s
}

func (d *Driver) ISR() {
	d.p.DisableIRQ(EvAll, ErrAll)
	d.done.Signal(1)
}

// SendCmd sends the cmd to the card and receives its response, if any. Short
// response is returned in r[0]. Long is returned in r[0:3] (r[0] contains the
// least significant bits, r[3] contains the most significant bits). If preceded
// by SetupData, SendCmd performs the data transfer.
func (d *Driver) SendCmd(cmd sdcard.Command, arg uint32) (r sdcard.Response) {
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
			r[3] = p.Resp(0) // Most significant bits.
			r[2] = p.Resp(1)
			r[1] = p.Resp(2)
			r[0] = p.Resp(3) // Least significant bits.
		} else {
			r[0] = p.Resp(0)
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
	ch := d.dma
	for {
		ev, err := ch.Status()
		if err != 0 {
			d.dmaErr = err
			break
		}
		if ev == dma.Complete {
			break
		}
	}
	ch.Disable() // Required by STM32F1 to allow setup next transfer.
	return
}

// SetupData setups the data transfer for subsequent command. On every call it
// configures DMA stream/channel completely from scratch so Driver can share its
// DMA stream/channel with other driver that do the same.
func (d *Driver) SetupData(mode sdcard.DataMode, buf sdcard.Data) {
	if d.err != 0 || d.dmaErr != 0 {
		return
	}
	d.data = mode
	dmacfg := dma.PFC | dma.IncM
	if mode&sdcard.Recv == 0 {
		dmacfg |= dma.MTP
	}
	if len(buf)&1 == 0 {
		dmacfg |= dma.FT4 | dma.PB4 | dma.MB4
	} else {
		dmacfg |= dma.FT2
	}
	ch := d.dma
	ch.Clear(dma.EvAll, dma.ErrAll)
	ch.Setup(dmacfg)
	ch.SetWordSize(4, 4)
	ch.SetAddrP(unsafe.Pointer(&d.p.raw.FIFO))
	ch.SetAddrM(unsafe.Pointer(&buf[0]))
	ch.Enable()
	p := d.p
	p.SetDataLen(len(buf) * 8)
	p.SetDataCtrl(DTEna | UseDMA | DataCtrl(mode))
}
