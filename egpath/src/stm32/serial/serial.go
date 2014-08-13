// Package serial provides high level interface to STM32 USART devices.
package serial

type USART interface {
	Store(byte)
	Load() byte
	Ready() (tx, rx bool)
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

func (s *Serial) IRQ() {
	tx, rx := s.dev.Ready()
	if rx {
		select {
		case s.rx <- s.dev.Load():
		default:
		}
	}
	if tx {
		select {
		case b := <-s.tx:
			s.dev.Store(b)
		default:
		}
	}
}

func (s *Serial) WriteByte(b byte) error {
	s.tx <- b
	return nil
}
