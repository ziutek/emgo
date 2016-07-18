package onedrv

import (
	"onewire"

	"stm32/hal/usart"
)

type USARTDriver struct {
	USART *usart.Driver
}

func (d USARTDriver) Reset() error {
	d.USART.SetBaudRate(9600)
	d.USART.WriteByte(0xf0)
	r, err := d.USART.ReadByte()
	if err != nil {
		return err
	}
	if r == 0xf0 {
		return onewire.ErrNoResponse
	}
	d.USART.SetBaudRate(115200)
	return nil
}

func (d USARTDriver) sendRecvSlot(slot byte) (byte, error) {
	d.USART.WriteByte(slot)
	return d.USART.ReadByte()
}

func (d USARTDriver) ReadBit() (byte, error) {
	slot, err := d.sendRecvSlot(0xff)
	if err != nil {
		return 0, err
	}
	return slot & 1, nil
}
func (d USARTDriver) WriteBit(bit byte) error {
	if bit != 0 {
		bit = 0xff
	}
	slot, err := d.sendRecvSlot(bit)
	if err != nil {
		return err
	}
	if slot != bit {
		return onewire.ErrBusFault
	}
	return nil
}

func (d USARTDriver) sendRecv(slots *[8]byte) error {
	d.USART.Write(slots[:])
	var n int
	for {
		m, err := d.USART.Read(slots[n:])
		if n += m; n == len(slots) {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

func (d USARTDriver) ReadByte() (byte, error) {
	slots := [8]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	err := d.sendRecv(&slots)
	var v int
	for i, slot := range slots {
		v += int(slot&1) << uint(i)
	}
	return byte(v), err
}

func (d USARTDriver) WriteByte(b byte) error {
	var slots [8]byte
	v := int(b)
	for i := range slots {
		if v&1 != 0 {
			slots[i] = 0xff
		}
		v >>= 1
	}
	if err := d.sendRecv(&slots); err != nil {
		return err
	}
	v = int(b)
	for i, slot := range slots {
		r := v & (1 << uint(i))
		if r != 0 {
			r = 0xff
		}
		if int(slot) != r {
			return onewire.ErrBusFault
		}
	}
	return nil
}
