package spi

import (
	"reflect"
	"rtos"
	"sync/fence"
	"unsafe"
)

// WriteReadByte sends a byte and returns received byte.
func (d *Driver) WriteReadByte(b byte) byte {
	p := d.P
	p.StoreTXD(b)
	ev := p.Event(READY)
	for !ev.IsSet() {
		rtos.SchedYield()
	}
	ev.Clear()
	return p.LoadRXD()
}

func (d *Driver) writeReadByteISR() {
	p := d.P
	ev := p.Event(READY)
	for ev.IsSet() {
		ev.Clear()
		b := p.LoadRXD()
		if rxptr := d.rxptr; rxptr <= d.rxmax {
			*(*byte)(unsafe.Pointer(rxptr)) = b
			d.rxptr = rxptr + 1
		}
		if txptr := d.txptr; txptr < d.txmax {
			p.StoreTXD(*(*byte)(unsafe.Pointer(txptr)))
			d.txptr = txptr + 1
		} else {
			if d.rxptr == d.rxmax {
				// There is still one byte to receive.
				continue
			} else if d.rxptr > d.rxmax {
				p.NVIC().ClearPending() // Can be edge triggered during ISR.
				d.done.Signal(1)
				return
			}
			p.StoreTXD(*(*byte)(unsafe.Pointer(txptr)))
		}
	}
}

// AsyncWriteStringRead starts SPI transaction: sending bytes from out string
// and receiving bytes into in slice. It returns immediately without waiting
// for end of transaction. See Wait for more infomation.
func (d *Driver) AsyncWriteStringRead(out string, in []byte) {
	d.n = len(in)
	if len(out) == 0 && len(in) == 0 {
		d.done.Reset(1)
		return
	}
	if len(out) <= 1 && len(in) <= 1 {
		b := byte(0xFF)
		if len(out) == 1 {
			b = out[0]
		}
		b = d.WriteReadByte(b)
		if len(in) == 1 {
			in[0] = b
		}
		d.done.Reset(1)
		return
	}
	// Now we are sure that len(out) >= 2 or len(in) >= 2 so at least two bytes
	// can be written to TXD register.
	p := d.P
	txptr := (*reflect.StringHeader)(unsafe.Pointer(&out)).Data
	if len(out) >= 2 {
		d.txmax = txptr + uintptr(len(out)) - 1
		p.StoreTXD(*(*byte)(unsafe.Pointer(txptr)))
		txptr++
		p.StoreTXD(*(*byte)(unsafe.Pointer(txptr)))
		if len(out) > 2 {
			txptr++
		}
		d.txptr = txptr
	} else {
		if len(out) == 0 {
			// Rx-only transaction. Send 0xFF bytes.
			d.w = 0xFF // nRF5 is little-endian.
			txptr = uintptr(unsafe.Pointer(&d.w))
		}
		d.txmax = txptr
		d.txptr = txptr
		p.StoreTXD(*(*byte)(unsafe.Pointer(txptr)))
		p.StoreTXD(*(*byte)(unsafe.Pointer(txptr)))
	}
	d.rxptr = (*reflect.StringHeader)(unsafe.Pointer(&in)).Data
	d.rxmax = d.rxptr + uintptr(len(in)) - 1
	d.isr = (*Driver).writeReadByteISR
	d.done.Reset(0)
	fence.W()
	p.Event(READY).EnableIRQ()
}

// WriteStringRead calls AsyncWriteStringRead followed by Wait.
func (d *Driver) WriteStringRead(out string, in []byte) int {
	d.AsyncWriteStringRead(out, in)
	return d.Wait()
}

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

func (d *Driver) repeatByteISR() {
	p := d.P
	ev := p.Event(READY)
	n := d.n
	b := byte(d.w)
	for ev.IsSet() {
		ev.Clear()
		p.LoadRXD()
		if n > 0 {
			p.StoreTXD(b)
		} else if n < 0 {
			p.NVIC().ClearPending() // Can be edge triggered during ISR.
			d.done.Signal(1)
			break
		}
		n--
	}
	d.n = n
}

func (d *Driver) AsyncRepeatByte(b byte, n int) {
	if n <= 1 {
		if n == 1 {
			d.WriteReadByte(b)
		}
		d.done.Reset(1)
		return
	}
	d.n = n - 2
	p := d.P
	p.StoreTXD(b)
	p.StoreTXD(b)
	d.w = uint16(b)
	d.isr = (*Driver).repeatByteISR
	d.done.Reset(0)
	fence.W()
	p.Event(READY).EnableIRQ()
}

func (d *Driver) RepeatByte(b byte, n int) {
	d.AsyncRepeatByte(b, n)
	d.Wait()
}
