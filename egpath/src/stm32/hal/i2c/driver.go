package i2c

import (
	"rtos"
	"sync"
	"unsafe"

	"arch/cortexm"

	"stm32/hal/dma"

	"stm32/hal/raw/i2c"
)

// Driver implements interrupt driven driver to I2C peripheral. To use the
// Driver the Periph field must be set and the I2CISR method must be setup as
// I2C interrupt handler for both (event and error) Periph's IRQs.
//
// Setting the RxDMA or/and TxDMA fields enables using DMA for Rx or/and Tx data
// transfer. If DMA is enabled for some direction the DMAISR method must be
// setup as DMA interrupt handler for this direction.
type Driver struct {
	P     *Periph
	RxDMA *dma.Channel
	TxDMA *dma.Channel

	mutex sync.Mutex
	done  rtos.EventFlag
	buf   []byte
	n     int
	addr  int16
	stop  bool
	state state
}

// NewDriver provides convenient way to create heap allocated Driver struct.
func NewDriver(p *Periph, rxdma, txdma *dma.Channel) *Driver {
	d := new(Driver)
	d.P = p
	d.RxDMA = rxdma
	d.TxDMA = txdma
	return d
}

// MasterConn returns initialized MasterConn struct that can be used to
// communicate with the slave device. Addr is the I2C address of the slave.
func (d *Driver) MasterConn(addr int16, stopm StopMode) MasterConn {
	return MasterConn{d: d, addr: addr << 1, stopm: stopm}
	// TODO: Add support for 10-bit addr.
}

// NewMasterConn is like MasterConn but returns pointer to heap allocated
// MasterConn struct.
func (d *Driver) NewMasterConn(addr int16, stopm StopMode) *MasterConn {
	mc := new(MasterConn)
	*mc = d.MasterConn(addr, stopm)
	return mc
}

// Unlock must be used after recovering from error.
func (d *Driver) Unlock() {
	d.mutex.Unlock()
}

// I2CISR is I2C (event and error) interrupt handler.
func (d *Driver) I2CISR() {
	sr1 := d.P.raw.SR1.Load()
	if e := getError(sr1); e != 0 {
		d.handleErrors()
		return
	}
	if d.state > stateReadN {
		d.state = stateBadEvent
		d.handleErrors()
		return
	}
	eventHandlers[d.state](d, sr1)
}

type state byte

const (
	stateIdle state = iota
	stateStart
	stateAddr

	stateWriteWait
	stateWrite

	stateReadWait
	stateRead
	stateRead1
	stateRead2
	stateReadN // Must be the last number of defined eventHandler (see I2CISR).

	stateWriteDMA
	stateReadDMA

	stateError
	stateBelatedStop
	stateBadEvent
	stateTimeout
)

//emgo:const
var eventHandlers = [...]func(d *Driver, sr1 i2c.SR1){
	stateIdle:  (*Driver).idleEH,
	stateStart: (*Driver).startISR,
	stateAddr:  (*Driver).addrISR,

	stateWrite:     (*Driver).writeISR,
	stateWriteWait: (*Driver).writeWaitEH,

	stateReadWait: (*Driver).readWaitEH,
	stateRead:     (*Driver).readISR,
	stateRead1:    (*Driver).read1ISR,
	stateRead2:    (*Driver).read2ISR,
	stateReadN:    (*Driver).readNISR,
}

func (d *Driver) idleEH(sr1 i2c.SR1) {
	p := &d.P.raw
	if d.addr < 0 {
		// Slave - not implemented.
	} else {
		// Master
		bits := i2c.START
		if d.addr&1 != 0 {
			bits |= i2c.ACK
		}
		p.CR1.SetBits(bits)
		d.state = stateStart
	}
	if sr1&i2c.BTF != 0 {
		// Repeated start (most likely).
		// Eensure that BTF was cleared before enable interrupts.
		maxrep := (d.P.Bus().Clock() + 16) / 32 // Timeout: 1/32 s.
		for {
			if getError(sr1) != 0 {
				d.handleErrors()
				return
			}
			sr1 = p.SR1.Load()
			if sr1&i2c.BTF == 0 {
				break
			}
			if maxrep == 0 {
				d.state = stateTimeout
				d.handleErrors()
				return
			}
			maxrep--
		}
	}
	d.enableIntI2C(i2c.ITEVTEN | i2c.ITERREN)
}

func (d *Driver) startISR(sr1 i2c.SR1) {
	if sr1&i2c.SB == 0 {
		d.badEvent(sr1)
		return
	}
	d.state = stateAddr
	p := &d.P.raw
	p.DR.Store(i2c.DR(d.addr))
	p.SR1.Load() // Ensure that SB was cleared before return.
}

func (d *Driver) addrISR(sr1 i2c.SR1) {
	if sr1&i2c.ADDR == 0 {
		d.badEvent(sr1)
		return
	}
	p := &d.P.raw
	if d.addr&1 == 0 {
		// Write.
		p.SR2.Load()
		d.write()
		return
	}
	// Read
	if d.RxDMA != nil && len(d.buf) > 1 {
		d.state = stateReadDMA
		p.SR2.Load()
		d.setupDMA(d.RxDMA, dma.PTM|dma.IncM|dma.FIFO_1_4)
		d.startDMA(d.RxDMA)
		return
	}
	if !d.stop {
		d.state = stateRead
		p.SR2.Load()
		return
	}
	switch len(d.buf) {
	case 1:
		d.state = stateRead1
		p.ACK().Clear()
		cortexm.SetPRIMASK()
		p.SR2.Load()
		p.STOP().Set()
		cortexm.ClearPRIMASK()
		d.enableIntI2C(i2c.ITEVTEN | i2c.ITERREN | i2c.ITBUFEN)
	case 2:
		d.state = stateRead2
		p.POS().Set()
		cortexm.SetPRIMASK()
		p.SR2.Load()
		p.ACK().Clear()
		cortexm.ClearPRIMASK()
	default:
		d.state = stateReadN
		p.SR2.Load()
	}
}

func (d *Driver) readISR(sr1 i2c.SR1) {
	if sr1&i2c.BTF == 0 {
		d.badEvent(sr1)
		return
	}
	d.buf[d.n] = byte(d.P.raw.DR.Load())
	d.n++
	if d.n == len(d.buf) {
		d.disableIntI2C(i2c.ITEVTEN | i2c.ITERREN | i2c.ITBUFEN)
		d.state = stateReadWait
		d.done.Signal(1)
	}
}

func (d *Driver) read1ISR(sr1 i2c.SR1) {
	if sr1&i2c.RXNE == 0 {
		d.badEvent(sr1)
		return
	}
	d.disableIntI2C(i2c.ITEVTEN | i2c.ITERREN | i2c.ITBUFEN)
	n := d.n
	d.buf[n] = byte(d.P.raw.DR.Load())
	d.n = n + 1
	d.state = stateIdle
	d.done.Signal(1)
}

func (d *Driver) read2ISR(sr1 i2c.SR1) {
	if sr1&i2c.BTF == 0 {
		d.badEvent(sr1)
		return
	}
	p := &d.P.raw
	cortexm.SetPRIMASK()
	p.STOP().Set()
	dr := p.DR.Load()
	cortexm.ClearPRIMASK()
	p.POS().Clear()
	d.disableIntI2C(i2c.ITEVTEN | i2c.ITERREN | i2c.ITBUFEN)
	d.buf[0] = byte(dr)
	d.buf[1] = byte(p.DR.Load())
	d.n = 2
	d.state = stateIdle
	d.done.Signal(1)
}

func (d *Driver) readNISR(sr1 i2c.SR1) {
	if sr1&i2c.BTF == 0 {
		d.badEvent(sr1)
		return
	}
	p := &d.P.raw
	n := d.n
	m := len(d.buf) - n
	if m < 3 {
		d.disableIntI2C(i2c.ITEVTEN | i2c.ITERREN | i2c.ITBUFEN)
		d.state = stateBelatedStop
		d.done.Signal(1)
		return
	}
	var dr i2c.DR
	if m == 3 {
		p.ACK().Clear()
		cortexm.SetPRIMASK()
		dr = p.DR.Load()
		p.STOP().Set()
		cortexm.ClearPRIMASK()
		d.buf[n] = byte(dr)
		n++
		dr = p.DR.Load()
		d.state = stateRead1
		// ITBUFEN must be set after second DR read.
		d.enableIntI2C(i2c.ITEVTEN | i2c.ITERREN | i2c.ITBUFEN)
	} else {
		dr = p.DR.Load()
	}
	d.buf[n] = byte(dr)
	d.n = n + 1
}

func (d *Driver) readWaitEH(_ i2c.SR1) {
	if d.RxDMA != nil && len(d.buf) > 1 {
		d.state = stateReadDMA
		if d.stop {
			d.P.raw.LAST().Set()
		}
		d.setupDMA(d.RxDMA, dma.PTM|dma.IncM|dma.FIFO_1_4)
		d.startDMA(d.RxDMA)
		return
	}
	if d.stop {
		d.state = stateReadN
	} else {
		d.state = stateRead
	}
	d.enableIntI2C(i2c.ITEVTEN | i2c.ITERREN)
}

func (d *Driver) writeISR(sr1 i2c.SR1) {
	if sr1&i2c.BTF == 0 {
		d.badEvent(sr1)
		return
	}
	n := d.n
	if n == len(d.buf) {
		d.disableIntI2C(i2c.ITEVTEN | i2c.ITERREN | i2c.ITBUFEN)
		d.state = stateWriteWait
		d.done.Signal(1)
		return
	}
	d.P.raw.DR.Store(i2c.DR(d.buf[n]))
	d.n = n + 1
}

func (d *Driver) writeWaitEH(_ i2c.SR1) {
	d.write()
	if d.state == stateWrite {
		d.enableIntI2C(i2c.ITEVTEN | i2c.ITERREN)
	}
}

func (d *Driver) write() {
	if d.TxDMA != nil && len(d.buf) > 1 {
		d.state = stateWriteDMA
		d.n = 0
		d.setupDMA(d.TxDMA, dma.MTP|dma.IncM|dma.FIFO_4_4)
		d.startDMA(d.TxDMA)
	} else {
		d.state = stateWrite
		d.n = 1
		d.P.raw.DR.Store(i2c.DR(d.buf[0]))
	}
}

func (d *Driver) setupDMA(ch *dma.Channel, mode dma.Mode) {
	d.disableDMA(ch)
	ch.Setup(mode)
	ch.SetWordSize(1, 1)
	ch.SetAddrP(unsafe.Pointer(d.P.raw.DR.U16.Addr()))
}

func (d *Driver) startDMA(ch *dma.Channel) {
	d.disableIntI2C(i2c.ITEVTEN | i2c.ITBUFEN)
	n := d.n
	m := len(d.buf) - n
	if m > 0xffff {
		m = 0xffff
	}
	if len(d.buf)-(n+m) == 1 {
		m-- // Avoid last transfer size 1.
	}
	d.n = n + m
	dmabits := i2c.DMAEN
	if d.stop && d.n == len(d.buf) {
		dmabits |= i2c.LAST
	}
	d.P.raw.CR2.SetBits(dmabits)
	ch.SetAddrM(unsafe.Pointer(&d.buf[n]))
	ch.SetLen(m)
	ch.Clear(dma.EvAll, dma.ErrAll)
	ch.EnableIRQ(dma.Complete, dmaErrMask)
	ch.Enable()
}

func (d *Driver) disableDMA(ch *dma.Channel) {
	ch.Disable()
	ch.DisableIRQ(dma.EvAll, dma.ErrAll)
	d.P.raw.CR2.ClearBits(i2c.DMAEN | i2c.LAST)
}

func (d *Driver) enableIntI2C(m i2c.CR2) {
	d.P.raw.CR2.SetBits(m)
}

func (d *Driver) disableIntI2C(m i2c.CR2) {
	d.P.raw.CR2.ClearBits(m)
}

func (d *Driver) handleErrors() {
	d.disableIntI2C(i2c.ITEVTEN | i2c.ITERREN | i2c.ITBUFEN)
	if d.RxDMA != nil {
		d.disableDMA(d.RxDMA)
	}
	if d.TxDMA != nil {
		d.disableDMA(d.TxDMA)
	}
	if d.state < stateError {
		d.state = stateError
	}
	d.done.Signal(1)
}

func (d *Driver) badEvent(sr1 i2c.SR1) {
	d.state = stateBadEvent
	d.handleErrors()
}

func (d *Driver) DMAISR(ch *dma.Channel) {
	d.disableDMA(ch)
	tx := d.addr&1 == 0
	if tx && d.state != stateWriteDMA || !tx && d.state != stateReadDMA {
		d.badEvent(0)
		return
	}
	if _, err := ch.Status(); err&dmaErrMask != 0 {
		d.n -= ch.Len()
		d.handleErrors()
		return
	}
	ch.Clear(dma.EvAll, dma.ErrAll)
	if len(d.buf)-d.n != 0 {
		d.startDMA(ch)
		return
	}
	if tx {
		d.state = stateWrite
		d.enableIntI2C(i2c.ITEVTEN | i2c.ITERREN)
		return
	}
	if d.stop {
		d.state = stateIdle
	} else {
		d.state = stateReadWait
	}
	d.done.Signal(1)
}

func (d *Driver) waitDone(ch *dma.Channel) (e Error) {
	timeout := byteTimeout + 2*9e9*int64(len(d.buf)+1)/int64(d.P.Speed())
	if !d.done.Wait(1, rtos.Nanosec() + timeout) {
		e = SoftTimeout
		d.disableIntI2C(i2c.ITEVTEN | i2c.ITERREN | i2c.ITBUFEN)
		if ch != nil {
			d.disableDMA(ch)
		}
		if d.state < stateError {
			d.state = stateError
		}
	}
	p := &d.P.raw
	if d.state >= stateError {
		switch d.state {
		case stateBelatedStop:
			e |= BelatedStop
		case stateBadEvent:
			e |= BadEvent
		case stateTimeout:
			e |= SoftTimeout
		}
		e |= getError(p.SR1.Load())
		if e&Timeout == 0 {
			p.STOP().Set()
		}
		p.SR1.Store(0) // Clear error flags.
		if ch != nil {
			if _, err := ch.Status(); err&dmaErrMask != 0 {
				e |= DMAErr
			}
			ch.Clear(dma.EvAll, dma.ErrAll)
		}
		d.state = stateIdle
	}
	return
}
