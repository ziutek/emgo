package onedrv

import (
	"io"

	"onewire"

	"stm32/hal/usart"
)

type USARTDriver struct {
	USART *usart.Driver
}

func (d USARTDriver) Reset() error {
	d.USART.SetBaudRate(9600)
	err := d.USART.WriteByte(0xf0)
	if err != nil {
		return err
	}
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

// resetRx resets internal ring buffer in usart.Driver.
func (d USARTDriver) resetRX() {
	// Clear Rx buffer and all errors.
	d.USART.DisableRx()
	d.USART.Status()
	d.USART.Load()
	d.USART.EnableRx()
}

func (d USARTDriver) sendRecvSlot(slot byte) (byte, error) {
	if err := d.USART.WriteByte(slot); err != nil {
		d.resetRX()
		return 0, err
	}
	b, err := d.USART.ReadByte()
	if err != nil {
		d.resetRX()
	}
	return b, err
}

func (d USARTDriver) sendRecv(slots *[8]byte) error {
	if _, err := d.USART.Write(slots[:]); err != nil {
		d.resetRX()
		return err
	}
	_, err := io.ReadFull(d.USART, slots[:])
	if err != nil {
		d.resetRX()
	}
	return err
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
		d.resetRX()
		return onewire.ErrBusFault
	}
	return nil
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
			d.resetRX()
			return onewire.ErrBusFault
		}
	}
	return nil
}

/*
func printStatus(u *usart.Driver) {
	ev, err := u.Status()

	fmt.Printf("Events: ")
	if ev&usart.Idle != 0 {
		fmt.Printf("Idle ")
	}
	if ev&usart.RxNotEmpty != 0 {
		fmt.Printf("RxNotEmpty ")
	}
	if ev&usart.TxDone != 0 {
		fmt.Printf("TxDone ")
	}
	if ev&usart.TxEmpty != 0 {
		fmt.Printf("TxEmpty ")
	}
	if ev&usart.LINBreak != 0 {
		fmt.Printf("LINBreak ")
	}
	if ev&usart.CTS != 0 {
		fmt.Printf("CTS ")
	}

	fmt.Printf("Errors:")
	if err&usart.ErrParity != 0 {
		fmt.Printf(" ErrParity")
	}
	if err&usart.ErrFraming != 0 {
		fmt.Printf(" ErrFraming")
	}
	if err&usart.ErrNoise != 0 {
		fmt.Printf(" ErrNoise")
	}
	if err&usart.ErrOverrun != 0 {
		fmt.Printf(" ErrOverrun")
	}
	fmt.Println()
}
*/
