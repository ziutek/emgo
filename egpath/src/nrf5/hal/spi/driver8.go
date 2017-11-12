package spi

import (
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
	// event loop before return..
	ev := p.Event(READY)
	for ev.IsSet() {
		ev.Clear()
		b := p.LoadRXD()
		if n := len(d.rxbuf); n < cap(d.rxbuf) {
			d.rxbuf = d.rxbuf[:n+1]
			d.rxbuf[n] = b
		}
		if d.txn >= len(d.txbuf) {
			switch len(d.rxbuf) {
			case cap(d.rxbuf) - 1:
				// There is still one byte to receive.
				continue
			case cap(d.rxbuf):
				d.done.Signal(1)
				return
			}
		}
		if d.txn < len(d.txbuf) {
			d.txn++
		}
		p.StoreTXD(d.txbuf[d.txn-1])
	}
}

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

func (d *Driver) WriteStringRead(out string, in []byte) int {
	d.AsyncWriteStringRead(out, in)
	return d.Wait()
}

func (d *Driver) WriteReadByte(b byte) byte {
	d.txbuf = "\xFF" // Set txbuf to any, one-byte string.
	d.txn = 1        // Mark txbuf as sent.
	var buf [1]byte
	d.rxbuf = buf[0:0:1]
	d.done.Reset(0)
	fence.W()
	d.P.StoreTXD(b)
	d.done.Wait(1, 0)
	return buf[0]
}

func (d *Driver) AsyncWriteRead(out, in []byte) {
	d.AsyncWriteStringRead(*(*string)(unsafe.Pointer(&out)), in)
}

func (d *Driver) WriteRead(out, in []byte) int {
	return d.WriteStringRead(*(*string)(unsafe.Pointer(&out)), in)
}
