package spi

import (
	"reflect"
	"rtos"
	"sync/fence"
	"unsafe"
)

// WriteReadByte sends byte and returns the received byte.
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

func (d *Driver) writeISR() {
	p := d.P
	ev := p.Event(READY)
	txptr := d.txptr
	txmax := d.txmax
	for ev.IsSet() {
		ev.Clear()
		p.LoadRXD()
		if txptr < txmax {
			p.StoreTXD(*(*byte)(unsafe.Pointer(txptr)))
		} else if txptr > txmax {
			p.NVIC().ClearPending() // Can be edge triggered during ISR.
			d.done.Signal(1)
			return
		}
		txptr++
	}
	d.txptr = txptr
}

func (d *Driver) readISR() {
	p := d.P
	ev := p.Event(READY)
	rxptr := d.rxptr
	rxmax := d.rxmax
	for ev.IsSet() {
		ev.Clear()
		*(*byte)(unsafe.Pointer(rxptr)) = p.LoadRXD()
		rxptr++
		if rxptr < rxmax {
			p.StoreTXD(0xFF)
		} else if rxptr > rxmax {
			p.NVIC().ClearPending() // Can be edge triggered during ISR.
			d.done.Signal(1)
			return
		}
	}
	d.rxptr = rxptr
}

func (d *Driver) writeReadISR() {
	p := d.P
	ev := p.Event(READY)
	txptr := d.txptr
	txmax := d.txmax
	rxptr := d.rxptr
	rxmax := d.rxmax
	for ev.IsSet() {
		ev.Clear()
		b := p.LoadRXD()
		if rxptr <= rxmax {
			*(*byte)(unsafe.Pointer(rxptr)) = b
			rxptr++
		}
		if txptr < txmax {
			p.StoreTXD(*(*byte)(unsafe.Pointer(txptr)))
			txptr++
		} else {
			if rxptr == rxmax {
				// There is still one byte to receive.
				continue
			} else if rxptr > rxmax {
				p.NVIC().ClearPending() // Can be edge triggered during ISR.
				d.done.Signal(1)
				return
			}
			p.StoreTXD(*(*byte)(unsafe.Pointer(txptr)))
		}
	}
	d.txptr = txptr
	d.rxptr = rxptr
}

// AsyncWriteStringRead starts SPI transaction: sending bytes from out string
// and receiving bytes into in slice. It returns immediately without waiting
// for end of transaction. See Wait for more infomation. len(ou) and len(in) may
// be different. If len(out) < len(in) the last byte of out is repeated as long
// as in will be filled. If len(out) == 0 and len(in) > 0 the 0xFF byte is sent.
func (d *Driver) AsyncWriteStringRead(out string, in []byte) {
	d.n = len(in)
	p := d.P
	if len(out) == 0 && len(in) == 0 {
		goto returnNoIRQ
	}
	if len(in) == 0 {
		// Tx-only transaction.
		if len(out) == 1 {
			d.WriteReadByte(out[0])
			goto returnNoIRQ
		}
		txptr := (*reflect.StringHeader)(unsafe.Pointer(&out)).Data
		p.StoreTXD(*(*byte)(unsafe.Pointer(txptr)))
		p.StoreTXD(*(*byte)(unsafe.Pointer(txptr + 1)))
		if len(out) == 2 {
			d.txptr = txptr // Can be any value.
			d.txmax = txptr // Same as in d.txptr, means: no more data to send.
		} else {
			d.txptr = txptr + 2
			d.txmax = txptr + uintptr(len(out)) - 1
		}
		d.isr = (*Driver).writeISR
		goto returnIRQ
	}
	if len(out) == 0 {
		// Rx-only transaction. Send 0xFF bytes.
		if len(in) == 1 {
			in[0] = d.WriteReadByte(0xFF)
			goto returnNoIRQ
		}
		p.StoreTXD(0xFF)
		p.StoreTXD(0xFF)
		d.w = uint16(0xFF) // nRF5 is little-endian.
		d.rxptr = (*reflect.StringHeader)(unsafe.Pointer(&in)).Data
		d.rxmax = d.rxptr + uintptr(len(in)) - 1
		d.isr = (*Driver).readISR
		goto returnIRQ
	}
	if len(in) == 1 && len(out) == 1 {
		in[0] = d.WriteReadByte(out[0])
		goto returnNoIRQ
	}
	// Now (len(out), len(in)) must be (>=1, >=2) or (>=2, >=1).
	{
		txptr := (*reflect.StringHeader)(unsafe.Pointer(&out)).Data
		d.txmax = txptr + uintptr(len(out)) - 1
		p.StoreTXD(*(*byte)(unsafe.Pointer(txptr)))
		if len(out) >= 2 {
			txptr++
		}
		p.StoreTXD(*(*byte)(unsafe.Pointer(txptr)))
		if len(out) > 2 {
			txptr++
		}
		d.txptr = txptr
		d.rxptr = (*reflect.StringHeader)(unsafe.Pointer(&in)).Data
		d.rxmax = d.rxptr + uintptr(len(in)) - 1
		d.isr = (*Driver).writeReadISR
	}
returnIRQ:
	d.done.Reset(0)
	fence.W()
	p.Event(READY).EnableIRQ()
	return
returnNoIRQ:
	d.done.Reset(1)
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
			return
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
