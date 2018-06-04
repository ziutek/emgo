package sdmmc

import (
	"unsafe"

	"sdcard"
)

// Driver implements sdcard.Host interface.
type Driver struct {
	driverCommon
	addr uintptr
	n    int
	err  Error
	dtc  DataCtrl
}

// MakeDriver returns initialized SPI driver that uses provided SPI peripheral.
func MakeDriver(p *Periph) Driver {
	return Driver{driverCommon: driverCommon{p: p}}
}

// NewDriver provides convenient way to create heap allocated Driver struct.
func NewDriver(p *Periph) *Driver {
	d := new(Driver)
	*d = MakeDriver(p)
	return d
}

func (d *Driver) Periph() *Periph {
	return d.p
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

// SetBusClock sets SD bus clock frequency (freqhz <= 0 disables clock). If
// pwrsave is true the clock output is automatically disabled when bus is idle.
func (d *Driver) SetBusClock(freqhz int, pwrsave bool) {
	d.setBusClock(freqhz, pwrsave)
}

// SetBusWidth sets the SD bus width.
func (d *Driver) SetBusWidth(width sdcard.BusWidth) {
	d.setBusWidth(width)
}

func (d *Driver) ISR() {
	p := d.p
	p.DisableIRQ(EvAll, ErrAll)
	addr := d.addr
	n := d.n
	if n == 0 {
		goto done
	}
	if d.dtc&Recv != 0 {
		for n >= 16 {
			ev, err := p.Status()
			_ = ev
			if err != 0 {
				goto done
			}
			addr = burstCopyPTM(p.raw.FIFO.Addr(), addr)
			n -= 16
		}
		for d.n > 0 {
			*(*uint32)(unsafe.Pointer(addr)) = p.Load()
			addr += 4
			d.n--
		}
	} else {

	}
done:
	d.done.Signal(1)
}

// SendCmd sends the cmd to the card and receives its response, if any. Short
// response is returned in r[0]. Long is returned in r[0:3] (r[0] contains the
// least significant bits, r[3] contains the most significant bits). If preceded
// by SetupData, SendCmd performs the data transfer.
func (d *Driver) SendCmd(cmd sdcard.Command, arg uint32) (r sdcard.Response) {
	if uint(d.err) != 0 {
		return
	}
	err := d.sendCmd(cmd, arg, &r)
	if err != 0 {
		d.err = err
		d.dtc = 0
		return
	}
	if d.dtc == 0 {
		return // No data transfer scheduled.
	}
	if d.dtc&Recv == 0 {
		d.p.SetDataCtrl(d.dtc)
	}

	return
}

// SetupData setups the data transfer for subsequent command.
func (d *Driver) SetupData(mode sdcard.DataMode, buf sdcard.Data) {
	if uint(d.err) != 0 {
		return
	}
	d.addr = uintptr(unsafe.Pointer(&buf[0]))
	d.n = len(buf) * 2
	p := d.p
	p.SetDataLen(len(buf) * 8)
	if d.dtc&Recv != 0 {
		p.SetDataCtrl(d.dtc)
	}
}
