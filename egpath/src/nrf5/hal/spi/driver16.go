package spi

import (
	"reflect"
	"rtos"
	"unsafe"
)

func (d *Driver) isr16() {
	p := d.P
	swp := int(d.swp)
	if d.txn == 0 {
		// New transaction.
		p.StoreTXD(d.txbuf[0^swp])
		p.StoreTXD(d.txbuf[1^swp])
		d.txn = 2
	}
	// SPI can generate events fast (1M event/s for max. speed) so check READY
	// event loop before return.
	ev := p.Event(READY)
	for ev.IsSet() {
		ev.Clear()
		b := p.LoadRXD()
		if n := len(d.rxbuf); n < cap(d.rxbuf) {
			d.rxbuf = d.rxbuf[:n+1]
			if &d.rxbuf[0] != nil {
				n ^= swp
				d.rxbuf[:n+1][n] = b
			}
		}
		if d.txn < len(d.txbuf) {
			p.StoreTXD(d.txbuf[d.txn^swp])
			d.txn++
		} else {
			switch len(d.rxbuf) {
			case cap(d.rxbuf) - 1:
				// There is still one byte to receive.
				continue
			case cap(d.rxbuf):
				d.done.Signal(1)
				return
			}
			p.StoreTXD(d.txbuf[(d.txn-1)^swp^len(d.rxbuf)&1])
		}
	}
}

func (d *Driver) AsyncWriteRead16(out, in []uint16) {
	txbuf := *(*reflect.StringHeader)(unsafe.Pointer(&out))
	txbuf.Len *= 2
	d.txbuf = *(*string)(unsafe.Pointer(&txbuf))
	d.txn = 0
	rxbuf := *(*reflect.SliceHeader)(unsafe.Pointer(&in))
	rxbuf.Cap = rxbuf.Len * 2
	rxbuf.Len = 0
	d.rxbuf = *(*[]byte)(unsafe.Pointer(&rxbuf))
	if len(out) == 0 {
		if len(in) == 0 {
			d.done.Reset(1)
			return
		}
		d.txbuf = "\xFF\xFF" // Rx-only mode: send 0xFF bytes.
	}
	d.done.Reset(0)
	rtos.IRQ(d.P.NVIC()).Trigger()
}

func (d *Driver) WriteRead16(out, in []uint16) int {
	d.AsyncWriteRead16(out, in)
	return d.Wait()
}

func (d *Driver) WriteReadWord16(w uint16) uint16 {
	txbuf := reflect.StringHeader{uintptr(unsafe.Pointer(&w)), 2}
	d.txbuf = *(*string)(unsafe.Pointer(&txbuf))
	d.txn = 0
	rxbuf := reflect.SliceHeader{uintptr(unsafe.Pointer(&w)), 0, 2}
	d.rxbuf = *(*[]byte)(unsafe.Pointer(&rxbuf))
	d.done.Reset(0)
	rtos.IRQ(d.P.NVIC()).Trigger()
	d.done.Wait(1, 0)
	return w
}

func (d *Driver) AsyncRepeatWord16(w uint16, n int) {
	d.rep = w
	txbuf := reflect.StringHeader{uintptr(unsafe.Pointer(&d.rep)), 2}
	d.txbuf = *(*string)(unsafe.Pointer(&txbuf))
	d.txn = 0
	rxbuf := reflect.SliceHeader{0, 0, n} // Means: discard the received bytes.
	d.rxbuf = *(*[]byte)(unsafe.Pointer(&rxbuf))
	d.done.Reset(0)
	rtos.IRQ(d.P.NVIC()).Trigger()
}

func (d *Driver) RepeatWord16(w uint16, n int) {
	d.AsyncRepeatWord16(w, n)
	d.Wait()
}
