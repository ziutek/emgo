package onewire

import (
	"errors"

	"stm32/serial"
)

// ROM commands
const (
	SearchROM   = 0xf0
	ReadROM     = 0x33
	MatchROM    = 0x55
	SkipROM     = 0xCC
	AlarmSearch = 0xec
)

// DS18B20 function commands
const (
	ConvertT = 0x44
)

var (
	ErrNoResponse = errors.New("no response")
	ErrBadEcho    = errors.New("bad echo")
)

type uartDriver struct {
	s    *serial.Dev
	pclk uint
}

func (d *uartDriver) reset() error {
	d.s.USART().SetBaudRate(9600, d.pclk)
	d.s.WriteByte(0xf0)
	r, err := d.s.ReadByte()
	if err != nil {
		return err
	}
	if r == 0xf0 {
		return ErrNoResponse
	}
	d.s.USART().SetBaudRate(115200, d.pclk)
	return nil
}

func (d *uartDriver) sendRecv(slot byte) (byte, error) {
	d.s.WriteByte(byte(slot))
	return d.s.ReadByte()
}

func (d *uartDriver) sendBit(bit byte) error {
	if bit != 0 {
		bit = 0xff
	}
	r, err := d.sendRecv(bit)
	if err != nil {
		return err
	}
	if r != bit {
		return ErrBadEcho
	}
	return nil
}

func (d *uartDriver) readBit() (bool, error) {
	r, err := d.sendRecv(0xff)
	if err != nil {
		return false, err
	}
	return r == 0xff, nil
}

type Master struct {
	d uartDriver
}

func NewMasterSerial(s *serial.Dev, pclk uint) *Master {
	m := new(Master)
	m.d = uartDriver{s, pclk}
	return m
}

// Reset resets 1-wire bus.
func (m *Master) Reset() error {
	return m.d.reset()
}

// SendByte sends byte b on 1-wire bus.
func (m *Master) SendByte(b byte) error {
	for i := 0; i < 8; i++ {
		if err := m.d.sendBit(b & 1); err != nil {
			return err
		}
		b >>= 1
	}
	return nil
}

// ReadByte reads byte from 1-wire bus.
func (m *Master) ReadByte() (byte, error) {
	var b int
	for p := 1; p != 0; p <<= 1 {
		bit, err := m.d.readBit()
		if err != nil {
			return 0, err
		}
		if bit {
			b += p
		}
	}
	return byte(b), nil
}
