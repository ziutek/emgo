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

// Driver is interrupt based driver to UART peripheral.
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

// EnableRx enables UART receiver. EnableRx must be called before any of Read*
// methods.
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

// DisableRx disables UART receiver.
func (d *Driver) DisableRx() {
	p := d.P
	p.Task(STOPRX).Trigger()
	p.DisableIRQ(1<<ERROR | 1<<RXDRDY)
}

// EnableTx enables UART transmitter. EnableTx must be called before any of
// Write* methods.
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

// DisableTx disables UART transmitter.
func (d *Driver) DisableTx() {
	p := d.P
	p.Task(STOPTX).Trigger()
	p.Event(TXDRDY).IRQ().Disable()
}

func (d *Driver) SetReadDeadline(t int64) {
	d.deadlineRx = t
}

func (d *Driver) SetWriteDeadline(t int64) {
	d.deadlineTx = t
}

// ISR should be used as UART interrupt handler.
func (d *Driver) ISR() {
	p := d.P
	for {
		again := false
		if p.Event(RXDRDY).IsSet() {
			p.Event(RXDRDY).Clear()
			b := p.LoadRXD() // Always read RXD to do not block RXDRDY event.
			nextpi := d.pi + 1
			if nextpi == len(d.RxBuf) {
				nextpi = 0
			}
			if atomic.LoadInt(&d.pr) == nextpi {
				atomic.OrUint32(&d.err, uint32(ErrBufOverflow)<<8)
			} else {
				d.RxBuf[d.pi] = b
				fence.W_SMP() // store(d.RxBuf) must be before store(d.pi).
				atomic.StoreInt(&d.pi, nextpi)
			}
			again = true
		}
		if p.Event(ERROR).IsSet() {
			p.Event(ERROR).Clear()
			err := p.LoadERRORSRC()
			p.ClearERRORSRC(err)
			atomic.OrUint32(&d.err, uint32(err))
			again = true
		}
		if again {
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
				p.StoreTXD(*(*byte)(unsafe.Pointer(d.txend + uintptr(newo))))
				atomic.StoreInt(&d.offs, newo)
				again = true
			}
		}
		if !again {
			break
		}
	}
}

// Len returns number of bytes that are ready to read from internal Rx buffer.
func (d *Driver) Len() int {
	n := atomic.LoadInt(&d.pi) - d.pr
	if n < 0 {
		n += len(d.RxBuf)
	}
	return n
}

func (d *Driver) clearError() error {
	err := atomic.SwapUint32(&d.err, 0)
	if pe := Error(err); pe != 0 {
		return pe
	}
	return DriverError(err >> 8)
}

func (d *Driver) ReadByte() (b byte, err error) {
	event := d.rxready
	if d.deadlineRx != 0 {
		event |= syscall.Alarm
	}
	for {
		if atomic.LoadUint32(&d.err) != 0 {
			err = d.clearError()
		}
		if pr := d.pr; atomic.LoadInt(&d.pi) != pr {
			fence.R_SMP() // Control dep. between load(d.pi) and load(d.RxBuf).
			b = d.RxBuf[pr]
			if pr++; pr == len(d.RxBuf) {
				pr = 0
			}
			fence.RW_SMP() // Ensure load(d.RxBuf) finished before store(d.pr).
			atomic.StoreInt(&d.pr, pr)
			return
		}
		if err != nil {
			return
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

func (d *Driver) Read(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, nil
	}
	event := d.rxready
	if d.deadlineRx != 0 {
		event |= syscall.Alarm
	}
	for {
		if atomic.LoadUint32(&d.err) != 0 {
			err = d.clearError()
		}
		if pr, pi := d.pr, atomic.LoadInt(&d.pi); pr != pi {
			fence.R_SMP() // Control dep. between load(d.pi) and load(d.RxBuf).
			if pi > pr {
				n = copy(b, d.RxBuf[pr:pi])
				pr += n
			} else {
				n = copy(b, d.RxBuf[pr:])
				if n < len(b) && pi != 0 {
					n += copy(b[n:], d.RxBuf[:pi])
				}
				if pr += n; pr >= len(d.RxBuf) {
					pr -= len(d.RxBuf)
				}
			}
			fence.RW_SMP() // Ensure load(d.RxBuf) finished before store(d.pr).
			atomic.StoreInt(&d.pr, pr)
			return
		}
		if err != nil {
			return
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

func (d *Driver) waitWrite() (int, error) {
	event, dl := d.txdone, d.deadlineTx
	if dl != 0 {
		event |= syscall.Alarm
	}
	for {
		offs := atomic.LoadInt(&d.offs)
		if offs == 0 {
			return 0, nil
		}
		if dl != 0 {
			if syscall.Nanosec() >= dl {
				return offs, ErrTimeout
			}
			syscall.SetAlarm(dl)
		}
		d.txdone.Wait()
	}
}

func (d *Driver) WriteByte(b byte) error {
	d.offs = -1
	fence.W() // store(d.offs) must be observed before p.SetTX.
	d.P.StoreTXD(b)
	_, err := d.waitWrite()
	return err
}

func (d *Driver) WriteString(s string) (int, error) {
	h := (*reflect.StringHeader)(unsafe.Pointer(&s))
	if h.Len == 0 {
		return 0, nil
	}
	d.txend = h.Data + uintptr(h.Len)
	d.offs = -h.Len
	fence.W() // store(d.offs) must be observed before p.SetTX.
	d.P.StoreTXD(s[0])
	return d.waitWrite()
}

func (d *Driver) Write(b []byte) (int, error) {
	return d.WriteString(*(*string)(unsafe.Pointer(&b)))
}
