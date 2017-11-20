// Package onedci usese STM32 UART/USART peripheral to implement onewire.DCI.
package onedci

import (
	"io"
	"onewire"
	"rtos"

	"stm32/hal/usart"
)

// USART wraps *usart.Driver to implement 1-Wire low-level signaling.
type USART struct {
	Drv *usart.Driver
}

func (dci USART) setReadTimeout() {
	dci.Drv.SetReadDeadline(rtos.Nanosec() + 1e9)
}

func (dci USART) Reset() error {
	dci.Drv.P.SetBaudRate(9600)
	err := dci.Drv.WriteByte(0xf0)
	if err != nil {
		return err
	}
	dci.setReadTimeout()
	r, err := dci.Drv.ReadByte()
	if err != nil {
		return err
	}
	if r == 0xf0 {
		return onewire.ErrNoResponse
	}
	dci.Drv.P.SetBaudRate(115200)
	return nil
}

// resetRx resets internal ring buffer in usart.Driver.
func (dci USART) resetRX() {
	// Clear Rx buffer and all errors.
	dci.Drv.DisableRx()
	dci.Drv.P.Status()
	dci.Drv.P.Load()
	dci.Drv.EnableRx()
}

func (dci USART) sendRecvSlot(slot byte) (byte, error) {
	if err := dci.Drv.WriteByte(slot); err != nil {
		dci.resetRX()
		return 0, err
	}
	dci.setReadTimeout()
	b, err := dci.Drv.ReadByte()
	if err != nil {
		dci.resetRX()
	}
	return b, err
}

func (dci USART) sendRecv(slots *[8]byte) error {
	if _, err := dci.Drv.Write(slots[:]); err != nil {
		dci.resetRX()
		return err
	}
	dci.setReadTimeout()
	_, err := io.ReadFull(dci.Drv, slots[:])
	if err != nil {
		dci.resetRX()
	}
	return err
}

func (dci USART) ReadBit() (byte, error) {
	slot, err := dci.sendRecvSlot(0xff)
	if err != nil {
		return 0, err
	}
	return slot & 1, nil
}
func (dci USART) WriteBit(bit byte) error {
	if bit != 0 {
		bit = 0xff
	}
	slot, err := dci.sendRecvSlot(bit)
	if err != nil {
		return err
	}
	if slot != bit {
		dci.resetRX()
		return onewire.ErrBusFault
	}
	return nil
}

func (dci USART) ReadByte() (byte, error) {
	slots := [8]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	err := dci.sendRecv(&slots)
	var v int
	for i, slot := range slots {
		v += int(slot&1) << uint(i)
	}
	return byte(v), err
}

func (dci USART) WriteByte(b byte) error {
	var slots [8]byte
	v := int(b)
	for i := range slots {
		if v&1 != 0 {
			slots[i] = 0xff
		}
		v >>= 1
	}
	if err := dci.sendRecv(&slots); err != nil {
		return err
	}
	v = int(b)
	for i, slot := range slots {
		r := v & (1 << uint(i))
		if r != 0 {
			r = 0xff
		}
		if int(slot) != r {
			dci.resetRX()
			return onewire.ErrBusFault
		}
	}
	return nil
}
