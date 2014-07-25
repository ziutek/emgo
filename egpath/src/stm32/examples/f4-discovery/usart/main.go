package main

import (
	"stm32/f4/gpio"
	"stm32/f4/periph"
	"stm32/f4/setup"
	"stm32/f4/usart"
)

var sp = usart.USART2

func init() {
	setup.Performance168(8)

	periph.AHB1ClockEnable(periph.GPIOA)
	periph.AHB1Reset(periph.GPIOA)
	periph.APB1ClockEnable(periph.USART2)
	periph.APB1Reset(periph.USART2)

	io, tx, rx := gpio.A, 2, 3

	io.SetMode(tx, gpio.Alt)
	io.SetOutType(tx, gpio.PushPullOut)
	io.SetPull(tx, gpio.PullUp)
	io.SetOutSpeed(tx, gpio.Fast)
	io.SetAltFunc(tx, gpio.USART2)
	io.SetMode(rx, gpio.Alt)

	sp.SetBaudRate(115200)
	sp.SetWordLen(usart.Bits8)
	sp.SetParity(usart.None)
	sp.SetStopBits(usart.Stop1b)
	sp.Enable()
	sp.EnableTx()
	sp.EnableRx()
}

type Serial struct {
	dev *usart.Dev
	tx  chan byte
	rx  chan byte
}

func NewSerial(dev *usart.Dev) *Serial {
	s := new(Serial)
	s.dev = dev
	s.tx = make(chan byte, 1)
	s.rx = make(chan byte, 1)
	return s
}

func (s *Serial) WriteByte(b byte) error {
	sp.Store(b)
	for sp.Status()&usart.TxEmpty == 0 {
	}
	return nil
}

func main() {
	s := NewSerial(sp)
	for {
		s.WriteByte('H')
		s.WriteByte('i')
		s.WriteByte('!')
		s.WriteByte('\r')
		s.WriteByte('\n')
	}
}
