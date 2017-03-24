package uart

import (
	"reflect"
	"sync/atomic"
	"sync/fence"
	"syscall"
	"unsafe"
)

type DriverError byte

const (
	ErrBufOverflow DriverError = iota + 1
	ErrTimeout
)

func (e DriverError) Error() string {
	switch e {
	case ErrBufOverflow:
		return "buffer overflow"
	case ErrTimeout:
		return "timeout"
	default:
		return ""
	}
}

// Driver is interrupt driven driver to UART peripheral.
type Driver struct {
	deadlineRx int64
	deadlineTx int64

	P     *Periph
	RxBuf []byte

	pi, pr  int
	err     uint32
	rxready syscall.Event

	offs   int
	txend  uintptr
	txdone syscall.Event
}

// NewDriver provides convenient way to create heap allocated Driver.
func NewDriver(p *Periph, rxbuf []byte) *Driver {
	d := new(Driver)
	d.P = p
	d.RxBuf = rxbuf
	return d
}

func (d *Driver) EnableRx() {
	if d.rxready == 0 {
		d.rxready = syscall.AssignEvent()
		fence.W() // Ensure rxready is stored before enable IRQ.
	}
	p := d.P
	p.Event(RXDRDY).Clear()
	p.Event(ERROR).Clear()
	p.ClearERRORSRC(ErrAll)
	p.EnableIRQ(1<<ERROR | 1<<RXDRDY)
	p.Task(STARTRX).Trigger()
}

func (d *Driver) DisableRx() {
	p := d.P
	p.Task(STOPRX).Trigger()
	p.DisableIRQ(1<<ERROR | 1<<RXDRDY)
}

func (d *Driver) EnableTx() {
	if d.txdone == 0 {
		d.txdone = syscall.AssignEvent()
		fence.W() // Ensure txdone is stored before enable IRQ.
	}
	p := d.P
	p.Event(TXDRDY).Clear()
	p.Event(TXDRDY).EnableIRQ()
	p.Task(STARTTX).Trigger()
}

func (d *Driver) DisableTx() {
	p := d.P
	p.Task(STOPTX).Trigger()
	p.Event(TXDRDY).IRQ().Disable()
}

func (d *Driver) ISR() {
	p := d.P
	for p.Event(RXDRDY).IsSet() {
		p.Event(RXDRDY).Clear()
		b := p.RXD() // Always read RXD to do not block RXDRDY event.
		nextpi := d.pi + 1
		if nextpi == len(d.RxBuf) {
			nextpi = 0
		}
		if nextpi == atomic.LoadInt(&d.pr) {
			atomic.OrUint32(&d.err, uint32(ErrBufOverflow)<<8)
			break
		}
		d.RxBuf[d.pi] = b
		fence.W_SMP() // store(d.RxBuf) must be observed before store(d.pi).
		atomic.StoreInt(&d.pi, nextpi)
	}
	for p.Event(ERROR).IsSet() {
		p.Event(ERROR).Clear()
		err := p.ERRORSRC()
		p.ClearERRORSRC(err)
		atomic.OrUint32(&d.err, uint32(err))
	}
	if atomic.LoadUint32(&d.err) != 0 || d.pi != atomic.LoadInt(&d.pr) {
		d.rxready.Send()
	}
	if p.Event(TXDRDY).IsSet() {
		p.Event(TXDRDY).Clear()
		newo := d.offs + 1
		if newo == 0 {
			fence.W() // clear(TXDRDY) must be observed before d.offs == 0.
			atomic.StoreInt(&d.offs, newo)
			d.txdone.Send()
		} else {
			p.SetTXD(*(*byte)(unsafe.Pointer(d.txend + uintptr(newo))))
			atomic.StoreInt(&d.offs, newo)
		}
	}
}

func (d *Driver) SetReadDeadline(t int64) {
	d.deadlineRx = t
}

func (d *Driver) SetWriteDeadline(t int64) {
	d.deadlineTx = t
}

func (d *Driver) clearError() error {
	err := atomic.SwapUint32(&d.err, 0)
	if pe := Error(err); pe != 0 {
		return pe
	}
	return DriverError(err >> 8)
}

// Len returns number of bytes that are ready to read from internal Rx buffer.
func (d *Driver) Len() int {
	n := atomic.LoadInt(&d.pi) - d.pr
	if n < 0 {
		n += len(d.RxBuf)
	}
	return n
}

func (d *Driver) ReadByte() (byte, error) {
	event := d.rxready
	if d.deadlineRx != 0 {
		event |= syscall.Alarm
	}
	for {
		if atomic.LoadUint32(&d.err) != 0 {
			return 0, d.clearError()
		}
		if pr := d.pr; atomic.LoadInt(&d.pi) != pr {
			fence.R_SMP() // Control dep. between load(d.pi) and load(d.RxBuf).
			b := d.RxBuf[pr]
			if pr++; pr == len(d.RxBuf) {
				pr = 0
			}
			fence.RW_SMP() // Ensure load(d.RxBuf) is before store(d.pr).
			atomic.StoreInt(&d.pr, pr)
			return b, nil
		}
		if dl := d.deadlineRx; dl != 0 {
			if syscall.Nanosec() >= dl {
				return 0, ErrTimeout
			}
			syscall.SetAlarm(dl)
		}
		event.Wait()
	}
}

/*
func (d *Driver) Read(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}
	if !d.rxready.Wait(1, d.deadlineRx) {
		return 0, ErrTimeout
	}
	d.rxready.Reset(0)
	pi := atomic.LoadInt(&d.pi)
}
*/

func (d *Driver) WriteByte(b byte) error {
	d.offs = -1
	fence.W() // store(d.offs) must be observed before p.SetTX.
	d.P.SetTXD(b)
	for atomic.LoadInt(&d.offs) != 0 {
		d.txdone.Wait()
	}
	return nil
}

func (d *Driver) WriteString(s string) (int, error) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	if sh.Len == 0 {
		return 0, nil
	}
	d.txend = sh.Data + uintptr(sh.Len)
	d.offs = -sh.Len
	fence.W() // store(d.offs) must be observed before p.SetTX.
	d.P.SetTXD(s[0])
	for atomic.LoadInt(&d.offs) != 0 {
		d.txdone.Wait()
	}
	return sh.Len, nil
}

func (d *Driver) Write(b []byte) (int, error) {
	return d.WriteString(*(*string)(unsafe.Pointer(&b)))
}
