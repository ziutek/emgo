package usart

import (
	"reflect"
	"rtos"
	"sync/atomic"
	"unsafe"

	"stm32/hal/dma"
)

type DriverError byte

const (
	ErrBufOverflow DriverError = iota
	ErrTimeout
)

func (e DriverError) Error() string {
	switch e {
	case ErrBufOverflow:
		return "buffer overflow"
	case ErrTimeout:
		return "timeout"
	default:
		return ""
	}
}

type Driver struct {
	deadlineRx int64
	deadlineTx int64

	*Periph
	RxDMA *dma.Channel
	TxDMA *dma.Channel
	RxBuf []byte // Rx ring buffer for RxDMA.

	txdone   rtos.EventFlag
	rxready  rtos.EventFlag
	rxM, rxN uint32
	dmaN     uint32
}

// NewDriver provides convenient way to create heap allocated Driver struct.
func NewDriver(p *Periph, rxdma, txdma *dma.Channel, rxbuf []byte) *Driver {
	d := new(Driver)
	d.Periph = p
	d.RxDMA = rxdma
	d.TxDMA = txdma
	d.RxBuf = rxbuf
	return d
}
func disableDMA(ch *dma.Channel) {
	ch.Disable()
	ch.DisableInt(dma.EvAll, dma.ErrAll)
}

func (d *Driver) setupDMA(ch *dma.Channel, mode dma.Mode) {
	ch.Setup(mode)
	ch.SetWordSize(1, 1)
	ch.SetAddrP(unsafe.Pointer(d.Periph.raw.DR.U16.Addr()))
}

func startDMA(ch *dma.Channel, maddr unsafe.Pointer, mlen int) {
	ch.SetAddrM(maddr)
	ch.SetLen(mlen)
	ch.Clear(dma.EvAll, dma.ErrAll)
	ch.Enable()
	ch.EnableInt(dma.Complete, dma.ErrAll&^dma.ErrFIFO) // Ignore FIFO error.
}

// EnableRx
func (d *Driver) EnableRx() {
	p := &d.Periph.raw
	ch := d.RxDMA
	p.RE().Set()
	p.DMAR().Set()
	d.setupDMA(ch, dma.PTM|dma.IncM|dma.Circ)
	startDMA(ch, unsafe.Pointer(&d.RxBuf[0]), len(d.RxBuf))
}

// DisableRx disables recieve of data and resets the state of internal ring
// buffer.
func (d *Driver) DisableRx() {
	p := &d.Periph.raw
	ch := d.RxDMA
	disableDMA(ch)
	p.RE().Clear()
	p.DMAR().Clear()
	d.rxN = 0
	d.rxM = 0
	for ch.Enabled() {
		// Wait dma really stops.
	}
	d.dmaN = 0
}

func (d *Driver) dmaNM() (n, m uint32) {
	ch := d.RxDMA
	n = atomic.LoadUint32(&d.dmaN)
	for {
		cl := ch.Len()
		nn := atomic.LoadUint32(&d.dmaN)
		if n == nn {
			return n, uint32(len(d.RxBuf) - cl)
		}
		n = nn
	}
}

func (d *Driver) rxNMadd(m int) {
	d.rxM += uint32(m)
	if d.rxM >= uint32(len(d.RxBuf)) {
		d.rxM -= uint32(len(d.RxBuf))
		d.rxN++
	}
}

func (d *Driver) Read(buf []byte) (int, error) {
	for {
		dmaN, dmaM := d.dmaNM()
		switch dmaN - d.rxN {
		case 0:
			if dmaM == d.rxM {
				d.Periph.EnableIRQ(RxNotEmpty)
				d.Periph.EnableErrorIRQ()
				if !d.rxready.Wait(d.deadlineRx) {
					return 0, ErrTimeout
				}
				d.rxready.Clear()
				if _, e := d.Periph.Status(); e != 0 {
					// Clear errors (complete "read SR then DR" sequence).
					d.Periph.Load()
					return 0, e
				}
				if _, e := d.RxDMA.Status(); e != 0 {
					return 0, e
				}
				continue
			}
			n := copy(buf, d.RxBuf[d.rxM:dmaM])
			d.rxNMadd(n)
			return n, nil
		case 1:
			if dmaM > d.rxM {
				break
			}
			n := copy(buf, d.RxBuf[d.rxM:])
			if n < len(buf) {
				n += copy(buf[n:], d.RxBuf[:dmaM])
			}
			dmaN, dmaM = d.dmaNM()
			if dmaN-d.rxN != 1 || dmaM > d.rxM {
				break
			}
			d.rxNMadd(n)
			return n, nil
		}
		d.rxN = dmaN
		d.rxM = dmaM
		return 0, ErrBufOverflow
	}
}

func (d *Driver) ReadByte() (byte, error) {
	var buf [1]byte
	_, err := d.Read(buf[:])
	return buf[0], err
}

func (d *Driver) ISR() {
	d.Periph.DisableIRQ(RxNotEmpty)
	d.Periph.Clear(RxNotEmpty)
	d.Periph.DisableErrorIRQ()
	d.rxready.Set()
}

func (d *Driver) RxDMAISR() {
	if _, e := d.RxDMA.Status(); e != 0 {
		d.rxready.Set()
		return
	}
	d.RxDMA.Clear(dma.EvAll, dma.ErrAll)
	atomic.AddUint32(&d.dmaN, 1)
}

func (d *Driver) EnableTx() {
	p := &d.Periph.raw
	p.TE().Set()
	p.DMAT().Set()
	d.setupDMA(d.TxDMA, dma.MTP|dma.IncM|dma.FIFO_4_4)
}

func (d *Driver) DisableTx() {
	p := &d.Periph.raw
	p.TE().Clear()
}

func (d *Driver) WriteString(s string) (int, error) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	var n int
	for {
		m := sh.Len - n
		if m == 0 {
			break
		}
		if m > 0xffff {
			m = 0xffff
		}
		d.Periph.raw.SR.Store(0) // Clear TC.
		startDMA(d.TxDMA, unsafe.Pointer(sh.Data+uintptr(n)), m)
		n += m
		if !d.txdone.Wait(d.deadlineTx) {
			return n - d.TxDMA.Len(), ErrTimeout
		}
		d.txdone.Clear()
		if _, e := d.RxDMA.Status(); e&^dma.ErrFIFO != 0 {
			return n - d.RxDMA.Len(), e
		}
	}
	return n, nil
}

func (d *Driver) Write(buf []byte) (int, error) {
	return d.WriteString(*(*string)(unsafe.Pointer(&buf)))
}

func (d *Driver) WriteByte(b byte) error {
	buf := [1]byte{b}
	_, err := d.Write(buf[:])
	return err
}

func (d *Driver) TxDMAISR() {
	disableDMA(d.TxDMA)
	d.txdone.Set()
}

func (d *Driver) SetReadDeadline(t int64) {
	d.deadlineRx = t
}

func (d *Driver) SetWriteDeadline(t int64) {
	d.deadlineTx = t
}

/*
func (d *Driver) NM() (dmaN, dmaM, rxN, rxM uint32) {
	ch := d.RxDMA
	n := atomic.LoadUint32(&d.dmaN)
	for {
		cl := ch.Len()
		nn := atomic.LoadUint32(&d.dmaN)
		if n == nn {
			return n, uint32(len(d.RxBuf) - cl), d.rxN, d.rxM
		}
		n = nn
	}

}
*/
