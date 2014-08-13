// Package serial provides high level interface to STM32 USART devices.
package serial

type USART interface {
	Store(b byte)
	Load() byte
	Ready() (tx, rx bool)
	TxIRQ(enable bool)
}

type Serial struct {
	dev USART
	tx  chan byte
	rx  chan byte
}

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

func (s *Serial) WriteByte(b byte) error {
	s.tx <- b
	s.IRQ()
	return nil
}

func (s *Serial) ReadByte() (byte, error) {
	return <-s.rx, nil
}
