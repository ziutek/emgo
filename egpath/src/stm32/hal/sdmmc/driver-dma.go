package sdmmc

import (
	"rtos"
	"sync/fence"
	"unsafe"

	"sdcard"

	"stm32/hal/dma"
)

// DriverDMA implements sdcard.Host interface using DMA.
type DriverDMA struct {
	p      *Periph
	dma    *dma.Channel
	done   rtos.EventFlag
	err    Error
	dmaErr dma.Error
	dtc    DataCtrl
}

// MakeDriverDMA returns initialized SPI driver that uses provided SPI
// peripheral and DMA channel.
func MakeDriverDMA(p *Periph, dma *dma.Channel) DriverDMA {
	return DriverDMA{p: p, dma: dma}
}

// NewDriverDMA provides convenient way to create heap allocated Driver struct.
func NewDriverDMA(p *Periph, dma *dma.Channel) *DriverDMA {
	d := new(DriverDMA)
	*d = MakeDriverDMA(p, dma)
	return d
}

func (d *DriverDMA) Periph() *Periph {
	return d.p
}

func (d *DriverDMA) DMA() *dma.Channel {
	return d.dma
}

func (d *DriverDMA) Err(clear bool) error {
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

// SetBusClock sets SD bus clock frequency (freqhz <= 0 disables clock). If
// pwrsave is true the clock output is automatically disabled when bus is idle.
func (d *DriverDMA) SetBusClock(freqhz int, pwrsave bool) {
	var (
		clkdiv int
		cfg    BusClock
		p      = d.p
	)
	busWidth, _ := p.BusClock()
	busWidth &= BusWidth
	if freqhz > 0 {
		// BUG: This code assumes 48 MHz SDMMCCLK.
		cfg = ClkEna
		clkdiv = (48e6+freqhz-1)/freqhz - 2
	}
	if clkdiv < 0 {
		clkdiv = 0
		cfg |= ClkByp
	}
	if pwrsave {
		cfg |= PwrSave
	}
	p.SetBusClock(cfg|busWidth, clkdiv)
	p.SetDataTimeout(uint(freqhz)) // ≈ 1s
}

// SetBusWidth sets the SD bus width.
func (d *DriverDMA) SetBusWidth(width sdcard.BusWidth) {
	if width > sdcard.Bus8 {
		panic("sdmmc: bad bus width")
	}
	p := d.p
	cfg, clkdiv := p.BusClock()
	cfg = cfg&^BusWidth | BusClock(width*3>>2)<<3
	p.SetBusClock(cfg, clkdiv)
}

func (d *DriverDMA) ISR() {
	d.p.DisableIRQ(EvAll, ErrAll)
	d.done.Signal(1)
}

// SendCmd sends the cmd to the card and receives its response, if any. Short
// response is returned in r[0]. Long is returned in r[0:3] (r[0] contains the
// least significant bits, r[3] contains the most significant bits). If preceded
// by SetupData, SendCmd performs the data transfer.
func (d *DriverDMA) SendCmd(cmd sdcard.Command, arg uint32) (r sdcard.Response) {
	if uint(d.err)|uint(d.dmaErr) != 0 {
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
	_, err := p.Status()
	if cmd&sdcard.HasResp != 0 {
		if err&ErrCmdCRC != 0 {
			switch cmd & sdcard.RespType {
			case sdcard.R3, sdcard.R4:
				err &^= ErrCmdCRC
			}
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
	if err != 0 {
		d.err = err
		d.dtc = 0
		return
	}
	if d.dtc == 0 {
		return // No data transfer scheduled.
	}
	if d.dtc&Recv == 0 {
		p.SetDataCtrl(d.dtc)
	}
	waitFor = 0
	if d.dtc&Stream == 0 {
		waitFor |= DataBlkEnd
	}
	d.dtc = 0
	d.done.Reset(0)
	fence.W() // This orders writes to normal and I/O memory.
	p.EnableIRQ(DataEnd, ErrAll)
	d.done.Wait(1, 0)
	for waitFor != 0 {
		ev, err := p.Status()
		if err != 0 {
			d.err = err
			d.dma.Disable()
			p.SetDataCtrl(0)
			return
		}
		waitFor &^= ev
	}
	ch := d.dma
	for {
		ev, err := ch.Status()
		if err &^= dma.ErrFIFO; err != 0 {
			d.dmaErr = err
			break
		}
		if ev&dma.Complete != 0 {
			break
		}
		/*if !ch.Enabled() {
			break  // STM32F103 RM says about waiting until channel disabled.
		}*/
	}
	return
}

// SetupData setups the data transfer for subsequent command. Ensure len(buf) <=
// 32767. SetupData configures DMA stream/channel completely from scratch so
// Driver can share its DMA stream/channel with other driver that do the same.
func (d *DriverDMA) SetupData(mode sdcard.DataMode, buf sdcard.Data) {
	if len(buf) > 32767 {
		panic("sdio: buf too big")
	}
	if uint(d.err)|uint(d.dmaErr) != 0 {
		return
	}
	d.dtc = DTEna | UseDMA | DataCtrl(mode)
	dmacfg := dma.PFC | dma.IncM
	if d.dtc&Recv == 0 {
		dmacfg |= dma.MTP
	}
	if len(buf)&1 == 0 {
		dmacfg |= dma.FT4 | dma.PB4 | dma.MB4
	} else {
		dmacfg |= dma.FT2
	}
	ch := d.dma
	ch.Disable()
	ch.Clear(dma.EvAll, dma.ErrAll)
	ch.Setup(dmacfg)
	ch.SetAddrP(unsafe.Pointer(&d.p.raw.FIFO))
	ch.SetAddrM(unsafe.Pointer(&buf[0]))
	ch.SetWordSize(4, 4)
	//ch.SetLen(len(buf) * 2) // Does  STM32F1 require this?
	ch.Enable()
	p := d.p
	p.SetDataLen(len(buf) * 8)
	if d.dtc&Recv != 0 {
		p.SetDataCtrl(d.dtc)
	}
}
