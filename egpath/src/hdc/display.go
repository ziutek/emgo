package hdc

import (
	"delay"
	"io"
)

// Display allows to send commands and data to HD44780 LCD controller in 4-bit
// mode.
//
// Every byte is sent as two 4-bit nibbles. One nibble is sent by writing three
// bytes to provided ReadWriter:
//
//	- first:  with E bit unset, need >= 40 ns
//	- second: with E bit set,   need >= 230 ns
//	- thrid:  with E bit unset, need >= 10 ns
//
// Full E cycle needs >= 500 ns. Display doesn't control proper nible timings,
// instead ReadWriter implementation must do it.
//
// Display reads busy flag before executing a command. Read command is written
// with all data bits set to 1.
type Display struct {
	Cols       byte
	Rows       byte
	ReadWriter io.ReadWriter
	DS         byte // Data shift: 4-bit nible uses bits DS to DS+3.
	E          byte // E line.
	RW         byte // R/W line.
	RS         byte // RS line.
	AUX        byte // AUX line (typically used for backlight).

	auxst byte
}

func (d *Display) SetAUX() error {
	d.auxst = d.AUX
	_, err := d.ReadWriter.Write([]byte{d.AUX, d.AUX, d.AUX})
	return err
}

func (d *Display) ClearAUX() error {
	d.auxst = 0
	_, err := d.ReadWriter.Write([]byte{0, 0, 0})
	return err
}

func writeByte(d *Display, rs, b byte) error {
	h := (b>>4)<<d.DS | rs | d.auxst
	l := (b&0xf)<<d.DS | rs | d.auxst
	buf := []byte{
		h, h | d.E, h,
		l, l | d.E, l,
	}
	if err := waitBusy(d); err != nil {
		return err
	}
	_, err := d.ReadWriter.Write(buf)
	return err
}

func writeCmd(d *Display, b byte) error {
	return writeByte(d, 0, b)
}

func writeData(d *Display, b byte) error {
	return writeByte(d, d.RS, b)
}

func waitBusy(d *Display) error {
	rba := d.RW | 0xf<<d.DS | d.auxst
	out := []byte{rba, rba | d.E, rba}
	var in [1]byte
	for {
		if _, err := d.ReadWriter.Write(out[:2]); err != nil {
			return err
		}
		if _, err := d.ReadWriter.Read(in[:]); err != nil {
			return err
		}
		if _, err := d.ReadWriter.Write(out[2:]); err != nil {
			return err
		}
		if _, err := d.ReadWriter.Write(out); err != nil {
			return err
		}
		if (in[0]>>d.DS)&0x8 == 0 {
			return nil
		}
	}
}

// Init initializes the display to the following state:
// - 4-bit mode,
// - one line display if Rows == 1, two line display otherwise,
// - 5x8 font,
// - display off, cursor off, blink off,
// - increment mode,
// - display cleared and cursor at home position.
func (d *Display) Init() error {
	// Set 8-bit mode
	//
	// Controller can be in 8-bit mode or in 4-bit mode (with upper nibble
	// received or not). So we should properly handle all three cases. We send
	// (multiple times) a command that enables 8-bit mode and works in both
	// modes when only 4 (upper) data pins are used.

	set8bit := byte(3) << d.DS
	buf := []byte{set8bit, set8bit | d.E, set8bit}

	delay.Millisec(40)

	// If in 4-bit mode this may be lower nibble of some previous command.
	if _, err := d.ReadWriter.Write(buf); err != nil {
		return err
	}
	delay.Millisec(5)

	// Now we are in 8-bit mode or this is upper nibble after previous command.
	if _, err := d.ReadWriter.Write(buf); err != nil {
		return err
	}
	delay.Millisec(1)

	// One more time.
	if _, err := d.ReadWriter.Write(buf); err != nil {
		return err
	}
	delay.Millisec(1)

	// Now we are certainly in 8-bit mode so set 4-bit mode and initialise.

	set4bit := byte(2) << d.DS
	buf = []byte{set4bit, set4bit | d.E, set4bit}

	if _, err := d.ReadWriter.Write(buf); err != nil {
		return err
	}
	delay.Millisec(1)

	// Now we are in 4-bit mode.

	// Some controller models may require to use SetFunction before any other
	// instuction.
	f := Font5x8
	if d.Rows > 1 {
		f |= TwoLines
	}
	if err := d.SetFunction(f); err != nil {
		return err
	}
	if err := d.SetDisplayMode(DisplayOff | CursorOff | BlinkOff); err != nil {
		return err
	}
	if err := d.SetEntryMode(Incr); err != nil {
		return err
	}
	if err := d.ReturnHome(); err != nil {
		return err
	}
	return d.ClearDisplay()
}

/*
func readBusyAddr(d *Display) (byte, error) {
	var in [2]byte
	out := []byte{d.RW, d.RW | d.E}

	if _, err := d.ReadWriter.Write(out); err != nil {
		return 0, err
	}
	if _, err := d.ReadWriter.Read(in[:1]); err != nil {
		return 0, err
	}
	if _, err := d.ReadWriter.Write(out[:1]); err != nil {
		return 0, err
	}
	if _, err := d.ReadWriter.Write(out); err != nil {
		return 0, err
	}
	if _, err := d.ReadWriter.Read(in[1:]); err != nil {
		return 0, err
	}
	if _, err := d.ReadWriter.Write(out[:1]); err != nil {
		return 0, err
	}
	return (in[0]>>d.DS)<<4 | (in[1]>>d.DS)&0xf, nil
}
*/
