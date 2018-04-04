package usart

import (
	"reflect"
	"rtos"
	"sync/atomic"
	"sync/fence"
	"unsafe"

	"stm32/hal/dma"
)

type DriverError byte

const (
	ErrBufOverflow DriverError = iota + 1
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

	P     *Periph
	RxDMA *dma.Channel
	TxDMA *dma.Channel
	RxBuf []byte // Rx ring buffer for RxDMA.

	txdone   rtos.EventFlag
	rxready  rtos.EventFlag
	rxH, rxL uint32
	dmaH     uint32
}

// NewDriver provides convenient way to create heap allocated Driver struct.
func NewDriver(p *Periph, rxdma, txdma *dma.Channel, rxbuf []byte) *Driver {
	d := new(Driver)
	d.P = p
	d.RxDMA = rxdma
	d.TxDMA = txdma
	d.RxBuf = rxbuf
	return d
}

func (d *Driver) setupDMA(ch *dma.Channel, mode dma.Mode, addr uintptr) {
	ch.Setup(mode)
	ch.SetWordSize(1, 1)
	ch.SetAddrP(unsafe.Pointer(addr))
}

func startDMA(ch *dma.Channel, maddr unsafe.Pointer, mlen int) {
	ch.SetAddrM(maddr)
	ch.SetLen(mlen)
	ch.Clear(dma.EvAll, dma.ErrAll)
	ch.EnableIRQ(dma.Complete, dma.ErrAll&^dma.ErrFIFO /* Ignore FIFO error */)
	fence.W() // This orders writes to normal and I/O memory.
	ch.Enable()
}

// EnableRx enables Rx part of P, setups RxDMA in circular mode and enables it
// to continuous reception of data. Driver assumes that it has exclusive access
// to P and RxDMA between EnableRx and DisableRx.
func (d *Driver) EnableRx() {
	p := &d.P.raw
	ch := d.RxDMA
	p.RE().Set()
	p.DMAR().Set()
	d.setupDMA(ch, dma.PTM|dma.IncM|dma.Circ, d.P.rdAddr())
	startDMA(ch, unsafe.Pointer(&d.RxBuf[0]), len(d.RxBuf))
}

// DisableRx disables recieve of data and resets the state of internal ring
// buffer.
func (d *Driver) DisableRx() {
	p := &d.P.raw
	ch := d.RxDMA
	ch.Disable()
	ch.DisableIRQ(dma.EvAll, dma.ErrAll)
	p.RE().Clear()
	p.DMAR().Clear()
	d.rxH = 0
	d.rxL = 0
	// Wait fo DMA really stops.
	for ch.Enabled() {
		rtos.SchedYield()
	}
	d.dmaH = 0
}

func (d *Driver) RxDMAISR() {
	ch := d.RxDMA
	ev, err := ch.Status()
	if err != 0 {
		d.rxready.Signal(1)
		return
	}
	if ev&dma.Complete != 0 {
		ch.Clear(dma.EvAll, dma.ErrAll)
		atomic.StoreUint32(&d.dmaH, d.dmaH+1)
	}
}

func (d *Driver) dmaHL() (h, l uint32) {
	ch := d.RxDMA
	h = atomic.LoadUint32(&d.dmaH)
	for {
		fence.R() // First read of dmaH must be executed before read of ch.Len.
		cl := ch.Len()
		fence.R() // Second read of dmaH must be executed after read of ch.Len.
		nh := atomic.LoadUint32(&d.dmaH)
		if h == nh {
			return h, uint32(len(d.RxBuf) - cl)
		}
		h = nh
	}
}

func (d *Driver) rxHLadd(n int) {
	d.rxL += uint32(n)
	if d.rxL >= uint32(len(d.RxBuf)) {
		d.rxL -= uint32(len(d.RxBuf))
		d.rxH++
	}
}

func (d *Driver) disableRxIRQ() {
	d.P.DisableIRQ(RxNotEmpty)
	d.P.Clear(RxNotEmpty, 0)
	d.P.DisableErrorIRQ()

}

func (d *Driver) ISR() {
	d.disableRxIRQ()
	d.rxready.Signal(1)
}

func (d *Driver) Read(buf []byte) (int, error) {
start:
	dmaH, dmaL := d.dmaHL()
	switch dmaH - d.rxH {
	case 0:
		if dmaL == d.rxL {
			// No data in ring buffer. Wait for RxNotEmpty or error (DMA IRQs
			// are useless because they can only signal half or full transfer.
			d.rxready.Reset(0)
			fence.W()
			d.P.EnableIRQ(RxNotEmpty)
			d.P.EnableErrorIRQ()
			dmaH, dmaL = d.dmaHL()
			if dmaL != d.rxL || dmaH != d.rxH {
				d.disableRxIRQ()
				goto start
			}
			if !d.rxready.Wait(1, d.deadlineRx) {
				return 0, ErrTimeout
			}
			if _, e := d.P.Status(); e != 0 {
				// Clear errors
				d.P.Load()      // For older MCUs (complete read SR then DR seq)
				d.P.Clear(0, e) // For newer MCUs.
				return 0, e
			}
			if _, e := d.RxDMA.Status(); e != 0 {
				return 0, e
			}
			goto start
		}
		if dmaL == 0 {
			// Belated RxDMAISR: dmaHL read NDTR just after it was reloaded and
			// before TC interrupt was taken.
			dmaL = uint32(len(d.RxBuf))
		}
		n := copy(buf, d.RxBuf[d.rxL:dmaL])
		d.rxHLadd(n)
		return n, nil
	case 1:
		if dmaL > d.rxL {
			break
		}
		n := copy(buf, d.RxBuf[d.rxL:])
		if n < len(buf) {
			n += copy(buf[n:], d.RxBuf[:dmaL])
		}
		dmaH, dmaL = d.dmaHL()
		if dmaH-d.rxH != 1 || dmaL > d.rxL {
			// There is no certainty that we managed to copy before overwriting.
			break
		}
		d.rxHLadd(n)
		return n, nil
	}
	d.rxH = dmaH
	d.rxL = dmaL
	return 0, ErrBufOverflow
}

func (d *Driver) ReadByte() (byte, error) {
	var buf [1]byte
	_, err := d.Read(buf[:])
	return buf[0], err
}

// EnableTx enables Tx part of P and setups TxDMA. Driver assumes that it has
// exclusive access to P and TxDMA between EnableTx and DisableTx.
func (d *Driver) EnableTx() {
	p := &d.P.raw
	p.TE().Set()
	p.DMAT().Set()
	d.setupDMA(d.TxDMA, dma.MTP|dma.IncM|dma.FIFO_4_4, d.P.tdAddr())
}

func (d *Driver) DisableTx() {
	p := &d.P.raw
	p.TE().Clear()
}

func (d *Driver) WriteString(s string) (int, error) {
	ch := d.TxDMA
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
		d.txdone.Reset(0)
		d.P.clear(TxDone, 0) // Clear TC.
		startDMA(ch, unsafe.Pointer(sh.Data+uintptr(n)), m)
		n += m
		done := d.txdone.Wait(1, d.deadlineTx)
		ch.Disable() // To be compatible with STM32F1.
		if !done {
			ch.DisableIRQ(dma.EvAll, dma.ErrAll)
			return n - ch.Len(), ErrTimeout
		}
		_, e := ch.Status()
		if e &^= dma.ErrFIFO; e != 0 {
			return n - ch.Len(), e
		}
	}
	return n, nil
}

func (d *Driver) Write(buf []byte) (int, error) {
	return d.WriteString(*(*string)(unsafe.Pointer(&buf)))
}

func (d *Driver) WriteByte(b byte) error {
	_, err := d.Write([]byte{b})
	return err
}

func (d *Driver) TxDMAISR() {
	ch := d.TxDMA
	if ev, err := ch.Status(); ev&dma.Complete != 0 || err != 0 {
		ch.DisableIRQ(dma.EvAll, dma.ErrAll)
		d.txdone.Signal(1)
	}
}

func (d *Driver) SetReadDeadline(t int64) {
	d.deadlineRx = t
}

func (d *Driver) SetWriteDeadline(t int64) {
	d.deadlineTx = t
}
