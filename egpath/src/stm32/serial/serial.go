// Package serial provides high level interface to STM32 USART devices.
package serial

import (
	"stm32/hal/usart"
)

// Dev wraps usart.Dev to provide high level interface to send and receive
// data using standard Read*/Write* methods family.
//
// It expects that provided usart.Dev is properly configured, has RxNotEmpty
// interrupt enabled and any other interrupts disabled. On its own Serial uses
// Status, Load, Store, EnableIRQs, DisableIRQs methods (last two are used to
// enable/disable TxEmpty and TxDone interrupts).
type Dev struct {
	dev    *usart.USART
	rx     chan uint16
	tx     chan uint16
	txdone chan struct{}
	unix   bool
	flush  bool
	saverr Error
	//fllf   bool
}

// New creates new Dev for USART device with Rx/Tx buffer of specified
// lengths in 9-bit words.
func New(dev *usart.USART, rxlen, txlen int) *Dev {
	s := new(Dev)
	s.dev = dev
	s.rx = make(chan uint16, rxlen)
	s.tx = make(chan uint16, txlen)
	s.txdone = make(chan struct{}, 1)
	return s
}

type Error byte

const (
	ErrParity Error = 1 << iota
	ErrFraming
	ErrNoise
	ErrOverrun
	ErrBufferFull
)

func (e Error) Error() string {
	switch e {
	case 0:
		return "no error"
	case ErrParity:
		return "parity error"
	case ErrFraming:
		return "framing error"
	case ErrNoise:
		return "noisy signal"
	case ErrOverrun:
		return "hardware buffer overrun"
	case ErrBufferFull:
		return "software buffer full"
	default:
		return "two or more errors from: Parity,Framing,Noise,Overrun,BufferFull"
	}
}

const flushReq = 1 << 15

// IRQ should be called by USART interrupt handler.
func (s *Dev) IRQ() {
	st := s.dev.Status()
	if st&usart.RxNotEmpty != 0 {
		err := Error(st&0xf) | s.saverr
		select {
		case s.rx <- uint16(s.dev.Load()&0x1ff | int(err)<<10):
			s.saverr = 0
		default:
			// Rx channel is full.
			s.saverr = err | ErrBufferFull
		}
	}
	if s.flush {
		if st&usart.TxDone == 0 {
			return
		}
		s.dev.DisableIRQs(usart.TxDoneIRQ)
		s.flush = false
		s.txdone <- struct{}{}
	}
	if st&usart.TxEmpty != 0 {
		select {
		case d := <-s.tx:
			if d == flushReq {
				if st&usart.TxDone != 0 {
					// Fast path.
					s.txdone <- struct{}{}
					break
				}
				s.flush = true
				s.dev.EnableIRQs(usart.TxDoneIRQ)
				break
			}
			s.dev.Store(int(d))
		default:
			// Tx channel is empty.
			s.dev.DisableIRQs(usart.TxEmptyIRQ)
		}
	} else {
		s.dev.EnableIRQs(usart.TxEmptyIRQ)
	}
}

// Flush waits for complete transmission of last word (including its stop bits)
// written to s.
func (s *Dev) Flush() error {
	s.tx <- flushReq
	<-s.txdone
	return nil
}

// SetUnix enabls/disables unix text mode. If enabled, every readed '\r' is
// translated to '\n' and  every writed '\n' is translated to "\r\n". This
// simple translation works well for many terminal emulators but not for all.
func (s *Dev) SetUnix(enable bool) {
	s.unix = enable
}

/*
// SetFlushLF enables/disables automatic flush after every '\n'.
func (s *Dev) SetFlushLF(enable bool) {
	s.fllf = enable
}
*/

// WriteBits can write 9-bit words to s. Text mode doesn't affect written data.
func (s *Dev) WriteBits(d uint16) error {
	s.tx <- d & 0x1ff
	s.dev.EnableIRQs(usart.TxEmptyIRQ)
	return nil
}

func (s *Dev) WriteByte(b byte) error {
	if s.unix && b == '\n' {
		s.WriteBits('\r')
	}
	s.WriteBits(uint16(b))
	/*if s.fllf && b == '\n' {
		s.Flush()
	}*/
	return nil
}

func split(d16 uint16) (d9 uint16, err error) {
	if e := Error(d16 >> 10); e != 0 {
		err = e
	}
	d9 = d16 & 0x1ff
	return
}

// ReadBits can read 9-bit words from s. Text mode doesn't affect read data.
func (s *Dev) ReadBits() (uint16, error) {
	return split(<-s.rx)
}

func (s *Dev) byte(d uint16) byte {
	b := byte(d)
	if s.unix && b == '\r' {
		b = '\n'
	}
	return b
}

func (s *Dev) ReadByte() (byte, error) {
	d, err := s.ReadBits()
	return s.byte(d), err
}

func (s *Dev) Write(buf []byte) (int, error) {
	for i, b := range buf {
		if err := s.WriteByte(b); err != nil {
			return i + 1, err
		}
	}
	return len(buf), nil
}

func (s *Dev) WriteString(str string) (int, error) {
	for i := 0; i < len(str); i++ {
		if e := s.WriteByte(str[i]); e != nil {
			return i + 1, e
		}
	}
	return len(str), nil
}

func (s *Dev) Read(buf []byte) (n int, err error) {
	if len(buf) == 0 {
		return
	}
	// Need to read at least one byte.
	buf[n], err = s.ReadByte()
	n++
	if err != nil {
		return
	}
	// Read next bytes until rx channel is empty.
	for n < len(buf) {
		select {
		case d := <-s.rx:
			d, err = split(d)
			buf[n] = s.byte(d)
			n++
			if err != nil {
				return
			}
		default:
			return
		}
	}
	return
}

// USART returns wrapped USART device.
func (s *Dev) USART() *usart.USART {
	return s.dev
}
