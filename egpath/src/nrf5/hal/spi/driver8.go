package spi

import (
	"reflect"
	"rtos"
	"sync/fence"
	"unsafe"
)

func (d *Driver) isr8() {
	p := d.P
	if d.txn == 0 {
		// New transaction.
		if len(d.txbuf) == 1 {
			b := d.txbuf[0]
			p.StoreTXD(b)
			d.txn = 1
			if cap(d.rxbuf) > 1 {
				p.StoreTXD(b)
			}
		} else {
			p.StoreTXD(d.txbuf[0])
			p.StoreTXD(d.txbuf[1])
			d.txn = 2
		}
	}
	// SPI can generate events fast (1M event/s for max. speed) so check READY
	// event in loop before return.
	ev := p.Event(READY)
	for ev.IsSet() {
		ev.Clear()
		b := p.LoadRXD()
		if n := len(d.rxbuf); n < cap(d.rxbuf) {
			d.rxbuf = d.rxbuf[:n+1]
			if &d.rxbuf[0] != nil {
				d.rxbuf[n] = b
			}
		}
		if d.txn < len(d.txbuf) {
			d.txn++
		} else {
			switch len(d.rxbuf) {
			case cap(d.rxbuf) - 1:
				// There is still one byte to receive.
				continue
			case cap(d.rxbuf):
				p.NVIC().ClearPending() // Can be edge triggered during ISR.
				d.done.Signal(1)
				return
			}
		}
		p.StoreTXD(d.txbuf[d.txn-1])
	}
}

// AsyncWriteStringRead starts SPI transaction: sending bytes from out string
// and receiving bytes into in slice. It returns immediately without waiting
// for end of transaction. See Wait for more infomation.
func (d *Driver) AsyncWriteStringRead(out string, in []byte) {
	d.txbuf = out
	d.txn = 0
	d.rxbuf = in[0:0:len(in)]
	if len(out) == 0 {
		if len(in) == 0 {
			d.done.Reset(1)
			return
		}
		d.txbuf = "\xFF" // Rx-only mode: send 0xFF bytes.
	}
	d.done.Reset(0)
	rtos.IRQ(d.P.NVIC()).Trigger()
}

// WriteStringRead calls AsyncWriteStringRead followed by Wait.
func (d *Driver) WriteStringRead(out string, in []byte) int {
	d.AsyncWriteStringRead(out, in)
	return d.Wait()
}

// WriteReadByte writes and reads byte.
func (d *Driver) WriteReadByte(b byte) byte {
	d.txbuf = "\xFF" // Set txbuf to any, one-byte string.
	d.txn = 1        // Mark txbuf as sent.
	buf := [1]byte{b}
	d.rxbuf = buf[0:0:1]
	d.done.Reset(0)
	fence.W()
	d.P.StoreTXD(b)
	d.done.Wait(1, 0)
	return buf[0]
}

// Clean and safe code ended. Magical and unsafe code begins.

// AsyncWriteRead starts SPI transaction: sending bytes from out slice and
// receiving bytes into in slice. It returns immediately without waiting for
// end of transaction. See Wait for more infomation.
func (d *Driver) AsyncWriteRead(out, in []byte) {
	d.AsyncWriteStringRead(*(*string)(unsafe.Pointer(&out)), in)
}

// WriteRead calls AsyncWriteRead followed by Wait.
func (d *Driver) WriteRead(out, in []byte) int {
	return d.WriteStringRead(*(*string)(unsafe.Pointer(&out)), in)
}

// AsyncRepeatByte starts SPI transaction that sends byte n times. It returns
// immediately without waiting for end of transaction. See Wait for more
// infomation.
func (d *Driver) AsyncRepeatByte(b byte, n int) {
	if n == 0 {
		d.done.Reset(1)
		return
	}
	(*[2]byte)(unsafe.Pointer(&d.rep))[0] = b
	txbuf := reflect.StringHeader{uintptr(unsafe.Pointer(&d.rep)), 1}
	d.txbuf = *(*string)(unsafe.Pointer(&txbuf))
	d.txn = 0
	rxbuf := reflect.SliceHeader{0, 0, n} // Means: discard the received bytes.
	d.rxbuf = *(*[]byte)(unsafe.Pointer(&rxbuf))
	d.done.Reset(0)
	rtos.IRQ(d.P.NVIC()).Trigger()
}

// RepeatByte calls AsyncRepeatByte followed by Wait.
func (d *Driver) RepeatByte(b byte, n int) {
	d.AsyncRepeatByte(b, n)
	d.Wait()
}
