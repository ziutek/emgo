package bcmw

import (
	"delay"
	"errors"
	"io"

	"sdcard"
)

var ErrTimeout = errors.New("bcmw: timeout")

type Driver struct {
	sd   sdcard.Host
	chip *Chip
	rca  uint16
}

func MakeDriver(sd sdcard.Host, chip *Chip) Driver {
	return Driver{sd: sd, chip: chip}
}

func NewDriver(sd sdcard.Host, chip *Chip) *Driver {
	d := new(Driver)
	*d = MakeDriver(sd, chip)
	return d
}

func (d *Driver) Init(reset func(nrst int)) error {
	sd := d.sd

	reset(0)
	sd.SetBusWidth(sdcard.Bus1)
	sd.SetClock(400e3, true)
	delay.Millisec(2)
	reset(1)

	var (
		retry int
		rca   uint16
	)

	// Enumerate

	for retry = 250; retry > 0; retry-- {
		delay.Millisec(2)
		sd.SendCmd(sdcard.CMD0())
		sd.SendCmd(sdcard.CMD5(0))
		rca, _ = sd.SendCmd(sdcard.CMD3()).R6()
		if sd.Err(true) == nil {
			break
		}
	}
	if retry == 0 {
		return ErrTimeout
	}

	// Select the card and put it in Transfer State.

	sd.SendCmd(sdcard.CMD7(rca))

	return sd.Err(true)
}

func (d *Driver) LoadFirmware(r io.Reader) error {
	return nil
}
