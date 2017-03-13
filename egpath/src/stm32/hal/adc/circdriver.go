package adc

import (
	"reflect"
	"sync/atomic"
	"sync/fence"
	"unsafe"

	"stm32/hal/dma"
)

type CircDriver struct {
	p  *Periph
	ch *dma.Channel

	hc      chan int32
	buf     []uint16
	err     uint32
	waitfor uint32
}

// NewCircDriver returns new circural driver using p and ch with internal
// buffers of total size 4*chunkLen bytes.
func NewCircDriver(p *Periph, ch *dma.Channel, chunkLen int) *CircDriver {
	d := new(CircDriver)
	d.p = p
	d.ch = ch
	d.hc = make(chan int32, 2)
	d.buf = make([]uint16, 2*chunkLen)
	return d
}

func (d *CircDriver) P() *Periph {
	return d.p
}

func (d *CircDriver) DMA() *dma.Channel {
	return d.ch
}

func (d *CircDriver) ISR() {
	waitfor := Event(atomic.LoadUint32(&d.waitfor))
	if waitfor == 0 {
		// Other ADC (shared IRQ).
		return
	}
	p := d.p
	ev, err := p.Status()
	if err == 0 && ev&waitfor == 0 {
		// Other ADC (shared IRQ).
		return
	}
	p.DisableIRQ(EvAll, ErrAll)
	if waitfor == ^EvAll {
		p.Clear(EvAll, ErrAll)
	} else {
		p.Clear(waitfor, ErrAll)
		atomic.StoreUint32(&d.waitfor, 0)
	}
	if err != 0 {
		atomic.OrUint32(&d.err, uint32(err))
	}
	select {
	case d.hc <- -1:
	default:
	}
}

func (d *CircDriver) DMAISR() {
	ch := d.ch
	ev, e := ch.Status()
	ch.Clear(dma.EvAll, dma.ErrAll)
	err := uint32(e&^dma.ErrFIFO) << 16
	var bh int32 // Buffer handle.
	switch ev & (dma.Complete | dma.HalfComplete) {
	case dma.Complete:
		bh = int32(len(d.buf) / 2)
	case dma.Complete | dma.HalfComplete:
		err |= uint32(ErrDrvOverrun) << 24
		// Calculate bh that points to more realiable half-buffer.
		if _, ws := ch.WordSize(); ch.Len()*int(ws) > len(d.buf) {
			bh = int32(len(d.buf) / 2)
		}
	}
	// Try send bh and next -1 to d.hc (use loop to save memory). Two values
	// are sent to detect an overrun. Receiver should read second value just
	// after it finish work with buffer pointed by bh. It should ignore negative
	// value.
	for {
		select {
		case d.hc <- bh:
		default:
			err |= uint32(ErrDrvOverrun) << 24
			bh = -1
		}
		if bh < 0 {
			if err != 0 {
				atomic.OrUint32(&d.err, err)
			}
			return
		}
		bh = -1
	}
}

func (d *CircDriver) watch(ev Event) {
	atomic.StoreUint32(&d.waitfor, uint32(ev))
	fence.W() // ISR must observe correct d.waitfor.
	p := d.p
	if ev == ^EvAll {
		p.Clear(EvAll, ErrAll)
		p.EnableIRQ(0, ErrAll)
	} else {
		p.Clear(ev, ErrAll)
		p.EnableIRQ(ev, ErrAll)
	}
	fence.W() // To order writes to normal and I/O memory.
}

func (d *CircDriver) Enable(calibrate bool) {
	d.enable(calibrate)
}

func (d *CircDriver) Start(wordSize, byteOffset uintptr) {
	p := d.p
	paddr := p.raw.DR.U32.Addr()
	if wordSize == 1 {
		paddr += byteOffset
	}
	enableDMA(
		d.ch, dma.Circ, dma.HalfComplete,
		unsafe.Pointer(paddr), unsafe.Pointer(&d.buf[0]),
		wordSize, len(d.buf)*2/int(wordSize),
	)
	p.EnableDMA(true)
	d.watch(^EvAll)
	acceptTrig(p)
}

func (d *CircDriver) Stop() {
	d.stopADC()
	d.p.DisableDMA()
	ch := d.ch
	ch.Disable()
	ch.DisableIRQ(dma.EvAll, dma.ErrAll)
	ch.Clear(dma.EvAll, dma.ErrAll)
}

// HandleChan returns the channel that can be used to obtain buffer handles.
func (d *CircDriver) HandleChan() <-chan int32 {
	return d.hc
}

func (d *CircDriver) Words16(bh int32) []uint16 {
	begin := int(bh)
	end := begin + len(d.buf)/2
	return d.buf[begin:end]
}

func (d *CircDriver) Bytes(bh int32) []byte {
	sli := *(*reflect.SliceHeader)(unsafe.Pointer(&d.buf))
	sli.Len *= 2
	sli.Cap *= 2
	begin := int(bh) * 2
	end := begin + len(d.buf)
	return (*(*[]byte)(unsafe.Pointer(&sli)))[begin:end]
}

func (d *CircDriver) Err() error {
	if atomic.LoadUint32(&d.err) == 0 {
		return nil
	}
	err := atomic.SwapUint32(&d.err, 0)
	if e := Error(err) & ErrAll; e != 0 {
		if atomic.LoadUint32(&d.waitfor) == uint32(^EvAll) {
			// Circural DMA. Reenable ADC IRQs to detect further errors.
			d.p.EnableIRQ(0, ErrAll)
		}
		return e
	}
	if e := dma.Error(err>>16) & dma.ErrAll; e != 0 {
		return e
	}
	if e := DriverError(err >> 24); e != 0 {
		return e
	}
	return nil
}
