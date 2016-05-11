package usart

import (
	"rtos"
	"unsafe"

	"stm32/hal/dma"

	"stm32/hal/raw/usart"
)

type Error byte

const (
	noErr    Error = 0
	ErrUSART Error = 1
	ErrDMA   Error = 2
)

func (e Error) Error() string {
	switch e {
	case ErrUSART:
		return "USART error"
	case ErrDMA:
		return "DMA error"
	case ErrUSART | ErrDMA:
		return "USART and DMA error"
	default:
		return ""
	}
}

type Driver struct {
	*Periph
	RxRing []byte // Rx ring buffer required for RxDMA.
	RxDMA  *dma.Channel
	TxDMA  *dma.Channel

	rxdone, txdone rtos.EventFlag
	rxbuf, txbuf   []byte
	rxn, txn       int
	rxe, txe       Error
}

const dmaErrMask = dma.ERR &^ dma.FFERR // Ignore FIFO error.

func (d *Driver) disableDMA(ch *dma.Channel, enbit usart.CR3_Bits) {
	ch.Disable()
	ch.DisableInt(dma.EV | dma.ERR)
	d.Periph.raw.CR3.ClearBits(enbit)
}

func (d *Driver) setupDMA(ch *dma.Channel, mode dma.Mode) {
	ch.ClearEvents(dma.EV | dma.ERR)
	ch.Setup(mode)
	ch.SetWordSize(1, 1)
	ch.SetAddrP(unsafe.Pointer(d.Periph.raw.DR.U16.Addr()))
}

func (d *Driver) startDMA(ch *dma.Channel, enbit usart.CR3_Bits, maddr unsafe.Pointer, mlen int) {
	d.Periph.raw.CR3.SetBits(enbit)
	ch.SetAddrM(maddr)
	ch.SetLen(mlen)
	ch.Enable()
	ch.EnableInt(dma.TRCE | dmaErrMask)
}

func (d *Driver) Write(buf []byte) (int, error) {
	ch := d.TxDMA
	d.disableDMA(ch, usart.DMAT)
	d.setupDMA(ch, dma.MTP|dma.IncM|dma.FIFO_4_4)
	m := len(buf)
	if m > 0xffff {
		m = 0xffff
	}
	d.txbuf = buf
	d.txn = m
	d.txe = noErr
	d.startDMA(ch, usart.DMAT, unsafe.Pointer(&buf[0]), m)
	d.txdone.Wait(0)
	d.txdone.Clear()
	if d.txe != 0 {
		return d.txn, d.txe
	}
	return d.txn, nil
}

func (d *Driver) TxDMAISR() {
	ch := d.TxDMA
	d.disableDMA(ch, usart.DMAT)
	if ch.Events()&dmaErrMask != 0 {
		d.txn -= ch.Len()
		d.txe |= ErrDMA
		d.txdone.Set()
		return
	}
	if d.txn == len(d.txbuf) {
		d.txdone.Set()
		return
	}
	m := len(d.txbuf) - d.txn
	if m > 0xffff {
		m = 0xffff
	}
	d.startDMA(ch, usart.DMAT, unsafe.Pointer(&d.txbuf[d.txn]), m)
	d.txn += m
}
