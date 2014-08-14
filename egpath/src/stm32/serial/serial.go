// Package serial provides high level interface to STM32 USART devices.
package serial

// USART is interface need by Serial to operate on some USART device.
type USART interface {
	Store(b byte)
	Load() byte
	Ready() (tx, rx bool)
	TxIRQ(enable bool)
}

// Serial provides high level interface to send and receive data on USART device. It
// uses interrupts to avoid pulling.
type Serial struct {
	dev  USART
	tx   chan byte
	rx   chan byte
	unix bool
}

// NewSerial creates new Serial for USART device dev with Tx/Rx buffer of specified
// lengths.
func NewSerial(dev USART, txlen, rxlen int) *Serial {
	s := new(Serial)
	s.dev = dev
	s.tx = make(chan byte, txlen)
	s.rx = make(chan byte, rxlen)
	return s
}

type Error int

const (
	RxBufErr Error = iota
)

var errStr = [...]string{
	"Rx buffer is full",
}

func (e Error) Error() string {
	return errStr[e]
}

// IRQ should be called by USART interrupt handler.
func (s *Serial) IRQ() (err error) {
	tx, rx := s.dev.Ready()
	if rx {
		select {
		case s.rx <- s.dev.Load():
		default:
			err = RxBufErr
		}
	}
	if tx {
		select {
		case b := <-s.tx:
			s.dev.Store(b)
		default:
			s.dev.TxIRQ(false)
		}
	} else {
		s.dev.TxIRQ(true)
	}
	return
}

// SetUnix enabls/disables unix text mode. If enabled, every readed '\r' is
// translated to '\n' and  every writed '\n' is translated to "\r\n". This
// simple translation works well for many terminal emulators but not for all.
func (s *Serial) SetUnix(enable bool) {
	s.unix = enable
}

func (s *Serial) WriteByte(b byte) error {
	if s.unix && b == '\n' {
		s.tx <- '\r'
		s.dev.TxIRQ(true)
	}
	s.tx <- b
	s.dev.TxIRQ(true)
	return nil
}

func (s *Serial) ReadByte() (byte, error) {
	b := <-s.rx
	if s.unix && b == '\r' {
		b = '\n'
	}
	return b, nil
}

func (s *Serial) Write(buf []byte) (int, error) {
	for i, b := range buf {
		if e := s.WriteByte(b); e != nil {
			return i + 1, e
		}
	}
	return len(buf), nil
}

func (s *Serial) WriteString(str string) (int, error) {
	for i := 0; i < len(str); i++ {
		if e := s.WriteByte(str[i]); e != nil {
			return i + 1, e
		}
	}
	return len(str), nil
}
