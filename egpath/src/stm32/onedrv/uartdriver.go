package onedrv

import (
	"stm32/serial"

	"onewire"
)

type UARTDriver struct {
	Serial *serial.Dev
	Clock  uint
}

func (d *UARTDriver) Reset() error {
	d.Serial.USART().SetBaudRate(9600, d.Clock)
	d.Serial.WriteByte(0xf0)
	r, err := d.Serial.ReadByte()
	if err != nil {
		return err
	}
	if r == 0xf0 {
		return onewire.ErrNoResponse
	}
	d.Serial.USART().SetBaudRate(115200, d.Clock)
	return nil
}

func (d *UARTDriver) sendRecv(slot byte) (byte, error) {
	d.Serial.WriteByte(byte(slot))
	return d.Serial.ReadByte()
}

func (d *UARTDriver) SendBit(bit byte) error {
	if bit != 0 {
		bit = 0xff
	}
	r, err := d.sendRecv(bit)
	if err != nil {
		return err
	}
	if r != bit {
		return onewire.ErrBusFault
	}
	return nil
}

func (d *UARTDriver) RecvBit() (byte, error) {
	r, err := d.sendRecv(0xff)
	if err != nil {
		return 0, err
	}
	return r & 1, nil
}
