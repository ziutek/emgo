package sdmmc

import (
	"rtos"
	"sync/fence"
	"unsafe"

	"sdcard"

	"stm32/hal/dma"
	"stm32/hal/exti"
	"stm32/hal/gpio"
)

// DriverDMA implements sdcard.Host interface using DMA.
type DriverDMA struct {
	p      *Periph
	dma    *dma.Channel
	d0     gpio.Pin
	done   rtos.EventFlag
	dmaErr dma.Error
	err    Error
	dtc    DataCtrl
}

// MakeDriverDMA returns initialized driver that uses provided SDMMC peripheral
// and DMA channel. If d0 is valid it also configures EXTI to detect rising edge
// on d0 pin.
func MakeDriverDMA(p *Periph, dma *dma.Channel, d0 gpio.Pin) DriverDMA {
	if d0.IsValid() {
		setupEXTI(d0)
	}
	return DriverDMA{p: p, dma: dma, d0: d0}
}

// NewDriverDMA provides convenient way to create heap allocated Driver struct.
func NewDriverDMA(p *Periph, dma *dma.Channel, d0 gpio.Pin) *DriverDMA {
	d := new(DriverDMA)
	*d = MakeDriverDMA(p, dma, d0)
	return d
}

func (d *DriverDMA) Periph() *Periph {
	return d.p
}

func (d *DriverDMA) DMA() *dma.Channel {
	return d.dma
}

// SetBusClock sets SD bus clock frequency (freqhz <= 0 disables clock).
func (d *DriverDMA) SetClock(freqhz int, pwrsave bool) {
	setClock(d.p, freqhz, pwrsave)
}

// SetBusWidth sets the SD bus width. It returns sdcard.SDBus1|sdcard.SDBus4.
func (d *DriverDMA) SetBusWidth(width sdcard.BusWidth) sdcard.BusWidths {
	return setBusWidth(d.p, width)
}

// Wait waits for deassertion of busy signal on DATA0 line. It returns false if
// the deadline has passed. Zero deadline means no deadline.
func (d *DriverDMA) Wait(deadline int64) bool {
	return wait(d.d0, &d.done, deadline)
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
		return nil
	}
	if clear {
		d.err = 0
		d.dmaErr = 0
		d.dtc = 0
	}
	return err
}

// BusyLine returns EXTI line used to detect end of busy state.
func (d *DriverDMA) BusyLine() exti.Lines {
	return exti.Lines(d.d0.Mask())
}

// BusyISR handles EXTI IRQ that detects end of busy state.
func (d *DriverDMA) BusyISR() {
	busyISR(d.d0, &d.done)
}

func (d *DriverDMA) ISR() {
	p := d.p
	p.SetIRQMask(0, 0)
	dtc := d.dtc
	if _, err := p.Status(); err != 0 || dtc&DTEna == 0 {
		d.done.Signal(1)
		return
	}
	if dtc&Recv == 0 {
		p.SetDataCtrl(d.dtc) // Start sending.
	}
	d.dtc = dtc &^ DTEna
	p.SetIRQMask(DataEnd, ErrAll)
}

// SetupData setups the data transfer for subsequent command. Data will be read
// from / write to buf. Ensure nbytes < 1<<24 and len(buf)*8 >= nbytes.
// SetupData configures DMA stream/channel completely from scratch so the Driver
// can share its DMA stream/channel with other driver that do the same.
func (d *DriverDMA) SetupData(mode sdcard.DataMode, buf []uint64, nbytes int) {
	if len(buf) == 0 {
		panicShortBuf()
	}
	if len(buf)*8 < nbytes {
		panic("sdmmc: buf too short")
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
	//ch.SetLen(len(buf) * 2) // Does STM32F1 require this? Use nbytes?
	ch.Enable()
	p := d.p
	p.SetDataLen(nbytes)
	if d.dtc&Recv != 0 {
		p.SetDataCtrl(d.dtc)
	}
}

// SendCmd sends the cmd to the card and receives its response, if any. Short
// response is returned in r[0]. Long is returned in r[0:3] (r[0] contains the
// least significant bits, r[3] contains the most significant bits). If preceded
// by SetupData, SendCmd performs the data transfer.
func (d *DriverDMA) SendCmd(cmd sdcard.Command, arg uint32) (r sdcard.Response) {
	if uint(d.err)|uint(d.dmaErr) != 0 {
		return
	}
	cmdEnd := CmdSent
	if cmd&sdcard.HasResp != 0 {
		cmdEnd = CmdRespOK
	}
	d.done.Reset(0)
	p := d.p
	p.Clear(EvAll, ErrAll)
	p.SetArg(arg)
	p.SetCmd(CmdEna | Command(cmd)&255)
	fence.W()                    // Orders writes to normal and IO memory.
	p.SetIRQMask(cmdEnd, ErrAll) // After SetCmd because of spurious IRQs.
	d.done.Wait(1, 0)
	if _, err := p.Status(); err != 0 {
		if rt := cmd & sdcard.RespType; rt == sdcard.R3 || rt == sdcard.R4 {
			err &^= ErrCmdCRC // Ignore CRC error for R3 and R4 responses.
		}
		if err != 0 {
			d.err = err
			return
		}
	}
	if cmd&sdcard.HasResp != 0 {
		if cmd&sdcard.LongResp != 0 {
			r[3] = p.Resp(0) // Most significant bits.
			r[2] = p.Resp(1)
			r[1] = p.Resp(2)
			r[0] = p.Resp(3) // Least significant bits.
		} else {
			r[0] = p.Resp(0)
		}
	}
	if d.dtc == 0 {
		return // No data transfer scheduled.
	}
	if d.dtc&Stream == 0 {
		// Wait for data CRC. It should be received and checked already so use
		// simple pooling.
		for {
			ev, err := p.Status()
			if err != 0 {
				d.err = err
				d.dma.Disable()
				return
			}
			if ev&DataBlkEnd != 0 {
				break
			}
		}
	}
	// Wait for end of DMA transfer. It should be ended already so use simple
	// pooling.
	ch := d.dma
	for {
		ev, err := ch.Status()
		if err &^= dma.ErrFIFO; err != 0 {
			d.dmaErr = err
			return
		}
		if ev&dma.Complete != 0 {
			break
		}
		/*if !ch.Enabled() {
			break  // STM32F103 RM says about waiting until channel disabled.
		}*/
	}
	d.dtc = 0
	return
}
