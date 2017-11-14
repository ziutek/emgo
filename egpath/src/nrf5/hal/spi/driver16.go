package spi

import (
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
				p.NVIC().ClearPending() // Can be edge triggered during ISR.
				d.done.Signal(1)
				break
			}
			p.StoreTXD(*(*byte)(unsafe.Pointer(txptr ^ swp ^ rxptr&1)))
		}
	}
	d.txptr = txptr
	d.rxptr = rxptr
}

// AsyncWriteRead16 starts SPI transaction: sending 16-bit words from out slice
// and receiving 16-bit words into in slice. It returns immediately without
// waiting for end of transaction. See Wait for more infomation.
func (d *Driver) AsyncWriteRead16(out, in []uint16) {
	d.n = len(in)
	if len(out) == 0 && len(in) == 0 {
		d.done.Reset(1)
		return
	}
	// Now we are sure that len(out) >= 1 or len(in) >= 1 so at least two bytes
	// can be written to TXD register.
	p := d.P
	var w uint16
	if len(out) >= 1 {
		txptr := (*reflect.SliceHeader)(unsafe.Pointer(&out)).Data
		w = *(*uint16)(unsafe.Pointer(txptr))
		if len(out) > 1 {
			d.txptr = txptr + 2
			d.txmax = txptr + uintptr(len(out))*2 - 1
		} else {
			d.txptr = txptr + 1
			d.txmax = d.txptr
		}
	} else {
		// Rx-only transaction. Send 0xFF bytes.
		d.txptr = uintptr(unsafe.Pointer(&d.w)) + 1
		d.txmax = d.txptr
		w = 0xFFFF
		d.w = w
	}
	if d.swp == 0 {
		p.StoreTXD(byte(w))
		p.StoreTXD(byte(w >> 8))
	} else {
		p.StoreTXD(byte(w >> 8))
		p.StoreTXD(byte(w))
	}
	if len(in) == 0 {
		d.rxptr = 0
		d.rxmax = 0
	} else {
		d.rxptr = (*reflect.StringHeader)(unsafe.Pointer(&in)).Data
		d.rxmax = d.rxptr + uintptr(len(in))*2 - 1
	}
	d.isr = (*Driver).writeRead16ISR
	d.done.Reset(0)
	fence.W()
	p.Event(READY).EnableIRQ()
}

// WriteRead16 calls AsyncWriteRead16 followed by Wait.
func (d *Driver) WriteRead16(out, in []uint16) int {
	d.AsyncWriteRead16(out, in)
	return d.Wait()
}
