// Package serial provides high level interface to STM32 USART devices.
package serial

import "delay"

type USART interface {
	Store(byte)
	Load() byte
	SetTxEmptyIRQ(bool)
	SetRxNotEmptyIRQ(bool)
}

type Serial struct {
	dev USART
	tx  chan byte
	rx  chan byte
}

func NewSerial(dev USART) *Serial {
	s := new(Serial)
	s.dev = dev
	s.tx = make(chan byte, 1)
	s.rx = make(chan byte, 1)
	return s
}

func (s *Serial) WriteByte(b byte) error {
	s.dev.Store(b)
	delay.Millisec(20)
	return nil
}
