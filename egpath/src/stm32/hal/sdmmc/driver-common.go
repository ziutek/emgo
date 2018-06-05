package sdmmc

import (
	"rtos"
	"sync/fence"

	"sdcard"
)

type driverCommon struct {
	p    *Periph
	done rtos.EventFlag
}

func (d *driverCommon) setClock(freqhz int, pwrsave bool) {
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

func (d *driverCommon) setBusWidth(width sdcard.BusWidth) {
	if width > sdcard.Bus8 {
		panic("sdmmc: bad bus width")
	}
	p := d.p
	cfg, clkdiv := p.BusClock()
	cfg = cfg&^BusWidth | BusClock(width*3>>2)<<3
	p.SetBusClock(cfg, clkdiv)
}

func (d *driverCommon) sendCmd(cmd sdcard.Command, arg uint32, r *sdcard.Response) Error {
	var waitFor Event
	if cmd&sdcard.HasResp != 0 {
		waitFor = CmdRespOK
	} else {
		waitFor = CmdSent
	}
	d.done.Reset(0)
	p := d.p
	p.Clear(EvAll, ErrAll)
	p.SetIRQMask(waitFor, ErrAll)
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
	return err
}
