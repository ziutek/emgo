package main

import (
	"fmt"
	"rtos"

	"nrf5/hal/clock"
	"nrf5/hal/gpio"
	"nrf5/hal/irq"
	"nrf5/hal/rtc"
	"nrf5/hal/spi"
	"nrf5/hal/system"
	"nrf5/hal/system/timer/rtcst"
	"nrf5/hal/uart"
)

var (
	u *uart.Driver
	s *spi.Driver
)

func init() {
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC0, 1)

	p0 := gpio.P0
	s = spi.NewDriver(spi.SPI0)
	s.P.StorePSEL(spi.SCK, p0.Pin(6))
	s.P.StorePSEL(spi.MISO, p0.Pin(8))
	s.P.StorePSEL(spi.MOSI, p0.Pin(10))
	s.P.StoreFREQUENCY(spi.Freq125k)
	s.Enable()
	rtos.IRQ(s.P.NVIC()).Enable()

	u = uart.NewDriver(uart.UART0, make([]byte, 80))
	u.P.StorePSEL(uart.RXD, p0.Pin(11))
	u.P.StorePSEL(uart.TXD, p0.Pin(9))
	u.P.StoreBAUDRATE(uart.Baud115200)
	u.Enable()
	//u.EnableRx()
	u.EnableTx()
	rtos.IRQ(u.P.NVIC()).Enable()
	fmt.DefaultWriter = u
}

func main() {
	in := make([]byte, 5)
	fmt.Printf(
		"\r\n\r\nn=%d\r\n",
		s.WriteRead([]byte{0x01, 0x23, 0x45, 0x67, 0x89}, nil),
	)
	fmt.Printf(
		"n=%d %02x\r\n",
		s.WriteRead(nil, in),
		in,
	)
	fmt.Printf(
		"n=%d %02x\r\n",
		s.WriteRead([]byte{0x01}, in),
		in,
	)
	fmt.Printf(
		"n=%d %02x\r\n",
		s.WriteRead([]byte{0x01, 0x23}, in),
		in,
	)
	fmt.Printf(
		"n=%d %02x\r\n",
		s.WriteRead([]byte{0x01, 0x23, 0x45}, in),
		in,
	)
	fmt.Printf(
		"n=%d %02x\r\n",
		s.WriteRead([]byte{0x01, 0x23, 0x45, 0x67}, in),
		in,
	)
	fmt.Printf(
		"n=%d %02x\r\n",
		s.WriteRead([]byte{0x01, 0x23, 0x45, 0x67, 0x89}, in),
		in,
	)
	fmt.Printf("%x\r\n", s.WriteReadByte(0xAB))
	s.AsyncRepeatByte(0x12, 100)
	fmt.Printf("n=%d\r\n", s.Wait())

	in16 := make([]uint16, 4)
	fmt.Printf(
		"\r\nn=%d\r\n",
		s.WriteRead16([]uint16{0x0123, 0x4567, 0x89AB, 0xCDEF}, nil),
	)
	fmt.Printf(
		"n=%d %04x\r\n",
		s.WriteRead16(nil, in16),
		in16,
	)
	fmt.Printf(
		"n=%d %04x\r\n",
		s.WriteRead16([]uint16{0x0123}, in16),
		in16,
	)
	fmt.Printf(
		"n=%d %04x\r\n",
		s.WriteRead16([]uint16{0x0123, 0x4567}, in16),
		in16,
	)
	fmt.Printf(
		"n=%d %04x\r\n",
		s.WriteRead16([]uint16{0x0123, 0x4567, 0x89AB}, in16),
		in16,
	)
	fmt.Printf(
		"n=%d %04x\r\n",
		s.WriteRead16([]uint16{0x0123, 0x4567, 0x89AB, 0xCDEF}, in16),
		in16,
	)
	fmt.Printf("%x\r\n", s.WriteReadWord16(0xABCD))
}

func spiISR() {
	s.ISR()
}

func uartISR() {
	u.ISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTC0:      rtcst.ISR,
	irq.SPI0_TWI0: spiISR,
	irq.UART0:     uartISR,
}
