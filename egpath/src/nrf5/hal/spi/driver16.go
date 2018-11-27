package spi

import (
	"bits"
	"reflect"
	"rtos"
	"sync/fence"
	"unsafe"
)

// WriteReadWord16 sends a 16-bit word and returns received 16-bit word.
func (d *Driver) WriteReadWord16(w uint16) uint16 {
	p := d.P
	ev := p.Event(READY)
	if d.swp == 0 {
		p.StoreTXD(byte(w))
		p.StoreTXD(byte(w >> 8))
		for !ev.IsSet() {
			rtos.SchedYield()
		}
		ev.Clear()
		w = uint16(p.LoadRXD())
		for !ev.IsSet() {
			rtos.SchedYield()
		}
		ev.Clear()
		return w | uint16(p.LoadRXD())<<8
	}
	p.StoreTXD(byte(w >> 8))
	p.StoreTXD(byte(w))
	for !ev.IsSet() {
		rtos.SchedYield()
	}
	ev.Clear()
	w = uint16(p.LoadRXD())
	for !ev.IsSet() {
		rtos.SchedYield()
	}
	ev.Clear()
	return w<<8 | uint16(p.LoadRXD())
}

func (d *Driver) writeSwpISR() {
	p := d.P
	ev := p.Event(READY)
	txptr := d.txptr
	txmax := d.txmax
	for ev.IsSet() {
		ev.Clear()
		p.LoadRXD()
		if txptr < txmax {
			p.StoreTXD(*(*byte)(unsafe.Pointer(txptr ^ 1)))
		} else if txptr > txmax {
			p.NVIRQ().ClearPending() // Can be edge triggered during ISR.
			d.done.Signal(1)
			return
		}
		txptr++
	}
	d.txptr = txptr
}

func (d *Driver) readSwpISR() {
	p := d.P
	ev := p.Event(READY)
	rxptr := d.rxptr
	rxmax := d.rxmax
	for ev.IsSet() {
		ev.Clear()
		*(*byte)(unsafe.Pointer(rxptr ^ 1)) = p.LoadRXD()
		rxptr++
		if rxptr < rxmax {
			p.StoreTXD(0xFF)
		} else if rxptr > rxmax {
			p.NVIRQ().ClearPending() // Can be edge triggered during ISR.
			d.done.Signal(1)
			return
		}
	}
	d.rxptr = rxptr
}

func (d *Driver) writeRead16ISR() {
	p := d.P
	ev := p.Event(READY)
	swp := uintptr(d.swp)
	txptr := d.txptr
	txmax := d.txmax
	rxptr := d.rxptr
	rxmax := d.rxmax
	for ev.IsSet() {
		ev.Clear()
		b := p.LoadRXD()
		if rxptr <= rxmax {
			*(*byte)(unsafe.Pointer(rxptr ^ swp)) = b
			rxptr++
		}
		if txptr < txmax {
			p.StoreTXD(*(*byte)(unsafe.Pointer(txptr ^ swp)))
			txptr++
		} else {
			if rxptr == rxmax {
				// There is still one byte to receive.
				continue
			} else if rxptr > rxmax {
				p.NVIRQ().ClearPending() // Can be edge triggered during ISR.
				d.done.Signal(1)
				return
			}
			p.StoreTXD(*(*byte)(unsafe.Pointer(txptr ^ swp ^ rxptr&1)))
		}
	}
	d.txptr = txptr
	d.rxptr = rxptr
}

// AsyncWriteRead16 starts SPI transaction: sending 16-bit words from out slice
// and receiving 16-bit words into in slice. It returns immediately without
// waiting for end of transaction. See Wait for more infomation. len(ou) and
// len(in) may be different. If len(out) < len(in) the last word of out is
// repeated as long as in will be filled. If len(out) == 0 and len(in) > 0 the
// 0xFFFF word is sent.
func (d *Driver) AsyncWriteRead16(out, in []uint16) {
	d.n = len(in)
	p := d.P
	if len(out) == 0 && len(in) == 0 {
		goto returnNoIRQ
	}
	if len(in) == 0 {
		// Tx-only transaction.
		if len(out) == 1 {
			d.WriteReadWord16(out[0])
			goto returnNoIRQ
		}
		txptr := (*reflect.SliceHeader)(unsafe.Pointer(&out)).Data
		if d.swp == 0 {
			p.StoreTXD(*(*byte)(unsafe.Pointer(txptr)))
			p.StoreTXD(*(*byte)(unsafe.Pointer(txptr + 1)))
			d.isr = (*Driver).writeISR
		} else {
			p.StoreTXD(*(*byte)(unsafe.Pointer(txptr + 1)))
			p.StoreTXD(*(*byte)(unsafe.Pointer(txptr)))
			d.isr = (*Driver).writeSwpISR
		}
		d.txptr = txptr + 2
		d.txmax = txptr + uintptr(len(out))*2 - 1
		goto returnIRQ
	}
	if len(out) == 0 {
		// Rx-only transaction. Send 0xFFFF words.
		if len(in) == 1 {
			in[0] = d.WriteReadWord16(0xFFFF)
			goto returnNoIRQ
		}
		p.StoreTXD(0xFF)
		p.StoreTXD(0xFF)
		d.rxptr = (*reflect.SliceHeader)(unsafe.Pointer(&in)).Data
		d.rxmax = d.rxptr + uintptr(len(in))*2 - 1
		if d.swp == 0 {
			d.isr = (*Driver).readISR
		} else {
			d.isr = (*Driver).readSwpISR
		}
		goto returnIRQ
	}
	if len(in) == 1 && len(out) == 1 {
		in[0] = d.WriteReadWord16(out[0])
		goto returnNoIRQ
	}
	// Now (len(out), len(in)) must be (>=1, >=2) or (>=2, >=1).
	{
		txptr := (*reflect.StringHeader)(unsafe.Pointer(&out)).Data
		if d.swp == 0 {
			p.StoreTXD(*(*byte)(unsafe.Pointer(txptr)))
			p.StoreTXD(*(*byte)(unsafe.Pointer(txptr + 1)))
		} else {
			p.StoreTXD(*(*byte)(unsafe.Pointer(txptr + 1)))
			p.StoreTXD(*(*byte)(unsafe.Pointer(txptr)))
		}
		if len(out) == 1 {
			d.txptr = txptr + 1
			d.txmax = d.txptr
		} else {
			d.txptr = txptr + 2
			d.txmax = txptr + uintptr(len(out))*2 - 1
		}
		d.rxptr = (*reflect.StringHeader)(unsafe.Pointer(&in)).Data
		d.rxmax = d.rxptr + uintptr(len(in))*2 - 1
		d.isr = (*Driver).writeRead16ISR
	}
returnIRQ:
	d.done.Reset(0)
	fence.W()
	p.Event(READY).EnableIRQ()
	return
returnNoIRQ:
	d.done.Reset(1)
}

// WriteRead16 calls AsyncWriteRead16 followed by Wait.
func (d *Driver) WriteRead16(out, in []uint16) int {
	d.AsyncWriteRead16(out, in)
	return d.Wait()
}

func (d *Driver) WriteReadMany16(oi ...[]uint16) int {
	var n int
	for k := 0; k < len(oi); k += 2 {
		var in []uint16
		if k+1 < len(oi) {
			in = oi[k+1]
		}
		out := oi[k]
		n += d.WriteRead16(out, in)
	}
	return n
}

func (d *Driver) repeatWord16ISR() {
	p := d.P
	ev := p.Event(READY)
	n := d.n
	w := d.w
	for ev.IsSet() {
		ev.Clear()
		p.LoadRXD()
		if n > 0 {
			p.StoreTXD(byte(w))
			w = bits.ReverseBytes16(w)
		} else if n < 0 {
			p.NVIRQ().ClearPending() // Can be edge triggered during ISR.
			d.done.Signal(1)
			return
		}
		n--
	}
	d.n = n
	d.w = w
}

func (d *Driver) AsyncRepeatWord16(w uint16, n int) {
	if n <= 1 {
		if n == 1 {
			d.WriteReadWord16(w)
		}
		d.done.Reset(1)
		return
	}
	d.n = n*2 - 2
	p := d.P
	if d.swp != 0 {
		w = bits.ReverseBytes16(w)
	}
	d.w = w
	p.StoreTXD(byte(w))
	p.StoreTXD(byte(w >> 8))
	d.isr = (*Driver).repeatWord16ISR
	d.done.Reset(0)
	fence.W()
	p.Event(READY).EnableIRQ()
}

func (d *Driver) RepeatWord16(w uint16, n int) {
	d.AsyncRepeatWord16(w, n)
	d.Wait()
}
