package uart

import (
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

	P *Periph

	rxbuf   []byte
	pi, pr  int
	err     uint32
	rxready syscall.Event

	txn    int
	txdata string
	txdone syscall.Event
}

// NewDriver provides convenient way to create heap allocated Driver.
func NewDriver(p *Periph, rxbuf []byte) *Driver {
	d := new(Driver)
	d.P = p
	d.rxbuf = rxbuf
	return d
}

func (d *Driver) Enable() {
	d.P.StoreENABLE(true)
}

func (d *Driver) Disable() {
	d.P.StoreENABLE(false)
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
	p.Event(TXDRDY).NVIRQ().Disable()
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
			if nextpi == len(d.rxbuf) {
				nextpi = 0
			}
			if atomic.LoadInt(&d.pr) == nextpi {
				atomic.OrUint32(&d.err, uint32(ErrBufOverflow)<<8)
			} else {
				d.rxbuf[d.pi] = b
				fence.W_SMP() // store(d.rxbuf) must be before store(d.pi).
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
			fence.W() // clear(TXDRDY) must be observed before d.offs == 0.
			atomic.StoreInt(&d.txn, d.txn+1)
			if uint(d.txn) < uint(len(d.txdata)) {
				// Uints above allow compiler to optimize bounds checking below.
				p.StoreTXD(d.txdata[d.txn])
				again = true
			} else {
				d.txdone.Send()
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
		n += len(d.rxbuf)
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

// ReadByte reads one byte from the internal buffer. ReadByte can block only if
// the internal buffer is empty.
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
			fence.R_SMP() // Control dep. between load(d.pi) and load(d.rxbuf).
			b = d.rxbuf[pr]
			if pr++; pr == len(d.rxbuf) {
				pr = 0
			}
			fence.RW_SMP() // Ensure load(d.rxbuf) finished before store(d.pr).
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

// Read reads into s data from the internal buffer. It returns number of bytes
// read and error if occured. It can return n < len(s) even if err == nil. Read
// blocks only if the internal buffer is empty.
func (d *Driver) Read(s []byte) (n int, err error) {
	if len(s) == 0 {
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
			fence.R_SMP() // Control dep. between load(d.pi) and load(d.rxbuf).
			if pi > pr {
				n = copy(s, d.rxbuf[pr:pi])
				pr += n
			} else {
				n = copy(s, d.rxbuf[pr:])
				if n < len(s) && pi > 0 {
					n += copy(s[n:], d.rxbuf[:pi])
				}
				if pr += n; pr >= len(d.rxbuf) {
					pr -= len(d.rxbuf)
				}
			}
			fence.RW_SMP() // Ensure load(d.rxbuf) finished before store(d.pr).
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

// WaitWrite waits for the end of the previous write which must be initiated by
// AsyncWrite or AsyncWriteString. It returns number of bytes written and error.
func (d *Driver) WaitWrite() (int, error) {
	event, dl := d.txdone, d.deadlineTx
	if dl != 0 {
		event |= syscall.Alarm
	}
	for {
		txn := atomic.LoadInt(&d.txn)
		if txn == len(d.txdata) {
			return txn, nil
		}
		if dl != 0 {
			if syscall.Nanosec() >= dl {
				return txn, ErrTimeout
			}
			syscall.SetAlarm(dl)
		}
		d.txdone.Wait()
	}
}

// AsyncWriteString works like AsyncWrite.
func (d *Driver) AsyncWriteString(s string) {
	d.txdata = s
	d.txn = 0
	if len(s) == 0 {
		return
	}
	fence.W() // New d.txdata, d.txn must be observed before store(TXD).
	d.P.StoreTXD(s[0])
}

// AsyncWrite initiates UART transmision of data referenced by s. This is
// dangerous function: you must ensure that data referenced by s are alive
// until subsequent WaitWrite return. In particular, there is probably always
// bad idea to use AsyncWrite with stack allocated data.
func (d *Driver) AsyncWrite(s []byte) {
	d.AsyncWriteString(*(*string)(unsafe.Pointer(&s)))
}

// WriteString works like Write.
func (d *Driver) WriteString(s string) (int, error) {
	if len(s) == 0 {
		return 0, nil
	}
	d.txdata = s
	d.txn = 0
	fence.W() // New d.txdata, d.txn must be observed before store(TXD).
	d.P.StoreTXD(s[0])
	return d.WaitWrite()
}

// Write transmits data referenced by s.
func (d *Driver) Write(s []byte) (int, error) {
	return d.WriteString(*(*string)(unsafe.Pointer(&s)))
}

// WriteByte transmits one byte.
func (d *Driver) WriteByte(b byte) error {
	d.txdata = ""
	d.txn = -1
	fence.W() // New d.txdata, d.txn must be observed before store(TXD).
	d.P.StoreTXD(b)
	_, err := d.WaitWrite()
	return err
}
