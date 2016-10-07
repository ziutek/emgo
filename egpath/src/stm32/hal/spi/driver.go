package spi

import (
	"reflect"
	"rtos"
	"sync/atomic"
	"unsafe"

	"stm32/hal/dma"
)

type DriverError byte

const ErrTimeout DriverError = 1

func (e DriverError) Error() string {
	switch e {
	case ErrTimeout:
		return "timeout"
	default:
		return ""
	}
}

type Driver struct {
	deadline int64

	P     *Periph
	RxDMA *dma.Channel
	TxDMA *dma.Channel

	dmacnt int32
	done   rtos.EventFlag
	err    uint
}

// NewDriver provides convenient way to create heap allocated Driver struct.
func NewDriver(p *Periph, rxdma, txdma *dma.Channel) *Driver {
	d := new(Driver)
	d.P = p
	d.RxDMA = rxdma
	d.TxDMA = txdma
	return d
}

func (d *Driver) setupDMA(ch *dma.Channel, mode dma.Mode, wordSize uintptr) {
	ch.Setup(mode)
	ch.SetWordSize(wordSize, wordSize)
	ch.SetAddrP(unsafe.Pointer(d.P.raw.DR.U16.Addr()))
}

func (d *Driver) DMAISR(ch *dma.Channel) {
	ch.Disable()
	ch.DisableInt(dma.EvAll, dma.ErrAll)
	_, e := ch.Status()
	if e&^dma.ErrFIFO != 0 || atomic.AddInt32(&d.dmacnt, -1) == 0 {
		d.done.Set()
	}
}

func (d *Driver) ISR() {
	d.P.DisableErrorIRQ()
	d.done.Set()
}

func startDMA(ch *dma.Channel, addr uintptr, n int) {
	ch.SetAddrM(unsafe.Pointer(addr))
	ch.SetLen(n)
	ch.Clear(dma.EvAll, dma.ErrAll)
	ch.Enable()
	ch.EnableInt(dma.Complete, dma.ErrAll&^dma.ErrFIFO) // Ignore FIFO error.
}

func (d *Driver) waitDone() bool {
	if !d.done.Wait(d.deadline) {
		d.TxDMA.Disable()
		d.RxDMA.Disable()
		d.done.Clear()
		d.err = uint(ErrTimeout) << 8
		return false
	}
	d.done.Clear()
	if _, e := d.P.Status(); e != 0 {
		d.err = uint(e)
		return false
	}
	_, rxe := d.RxDMA.Status()
	_, txe := d.TxDMA.Status()
	if e := (rxe | txe) &^ dma.ErrFIFO; e != 0 {
		d.err = uint(e) << 16
		return false
	}
	return true
}

func (d *Driver) writeRead(out, in uintptr, olen, ilen int) int {
	txdmacfg := dma.MTP | dma.FIFO_4_4
	if olen > 1 {
		txdmacfg |= dma.IncM
	}
	d.setupDMA(d.TxDMA, txdmacfg, 1)
	d.setupDMA(d.RxDMA, dma.PTM|dma.IncM|dma.FIFO_1_4, 1)
	d.P.SetDuplex(Full)
	d.P.SetDMA(RxNotEmpty | TxEmpty)
	d.P.EnableErrorIRQ()
	var n int
	for {
		m := ilen - n
		if m == 0 {
			return n
		}
		if m > 0xffff {
			m = 0xffff
		}
		atomic.StoreInt32(&d.dmacnt, 2)
		startDMA(d.RxDMA, in, m)
		startDMA(d.TxDMA, out, m)
		if olen > 1 {
			out += uintptr(m)
		}
		in += uintptr(m)
		n += m
		if !d.waitDone() {
			return n - d.RxDMA.Len()
		}
	}
}

func (d *Driver) write(out uintptr, olen int) {
	d.setupDMA(d.TxDMA, dma.MTP|dma.IncM|dma.FIFO_4_4, 1)
	d.P.SetDuplex(HalfOut) // Avoid ErrOverflow.
	d.P.SetDMA(TxEmpty)
	d.P.EnableErrorIRQ()
	var n int
	for {
		m := olen - n
		if m == 0 {
			return
		}
		if m > 0xffff {
			m = 0xffff
		}
		atomic.StoreInt32(&d.dmacnt, 1)
		startDMA(d.TxDMA, out+uintptr(n), m)
		n += m
		if !d.waitDone() {
			return
		}
	}
}

// Err returns the first error that was encountered by the Driver. It also
// clears internal error flags so subsequent Err call returns nil or next error.
func (d *Driver) Err() error {
	e := d.err
	if e == 0 {
		return nil
	}
	d.err = 0
	if err := e >> 16; err != 0 {
		return dma.Error(err)
	}
	if err := e >> 8; err != 0 {
		return DriverError(e)
	}
	err := Error(e)
	if err&ErrOverrun != 0 {
		d.P.LoadByte()
		d.P.Status()
	}
	return err
}

var ffff uint16 = 0xffff

func (d *Driver) WriteStringRead(out string, in []byte) int {
	olen := len(out)
	ilen := len(in)
	if d.err != 0 || olen == 0 && ilen == 0 {
		return 0
	}
	oaddr := (*reflect.StringHeader)(unsafe.Pointer(&out)).Data
	iaddr := (*reflect.SliceHeader)(unsafe.Pointer(&in)).Data
	if olen > ilen {
		var n int
		if ilen > 0 {
			n = d.writeRead(oaddr, iaddr, ilen, ilen)
			if d.err != 0 {
				return n
			}
			olen -= ilen
			oaddr += uintptr(ilen)
		}
		d.write(oaddr, olen)
		return n
	}
	if ilen < olen {
		var n int
		if olen > 0 {
			n = d.writeRead(oaddr, iaddr, olen, olen)
			if d.err != 0 {
				return n
			}
			ilen -= olen
			iaddr += uintptr(olen)
			oaddr += uintptr(olen - 1)
		} else {
			oaddr = uintptr(unsafe.Pointer(&ffff))
		}
		return n + d.writeRead(oaddr, iaddr, 1, ilen)
	}
	return d.writeRead(oaddr, iaddr, ilen, ilen)
}

func (d *Driver) WriteRead(out, in []byte) int {
	return d.WriteStringRead(*(*string)(unsafe.Pointer(&out)), in)
}

/*
func (d *Driver) WriteRead(oi ...[]byte) (n int, err error) {
	for k := 0; k < len(oi); k += 2 {
		var in []byte
		if k+1 < len(oi) {
			in = oi[k+1]
		}
		out := oi[k]
		if m := len(out); m == len(in) {
			if m == 0 {
				continue
			}
			m, err = d.writeRead(
				unsafe.Pointer(&out[0]), unsafe.Pointer(&in[0]),
				m, m,
			)
			n += m
			if err != nil {
				return n, err
			}
		} else {
			panic("nwww")
		}
	}
	return n, err
}
*/
