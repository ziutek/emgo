package main

import (
	"delay"
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

type EVE struct {
	spi           *spi.Driver
	pdn, csn, irq gpio.Pin
}

func (lcd *EVE) Cmd(cmd HostCmd) {
	lcd.csn.Clear()
	lcd.spi.WriteRead([]byte{byte(cmd), 0, 0}, nil)
	lcd.csn.Set()
}

func (lcd *EVE) Read8(addr uint32) byte {
	lcd.csn.Clear()
	buf := []byte{byte(addr >> 16), byte(addr >> 8), byte(addr), 0, 0}
	lcd.spi.WriteRead(buf, buf)
	lcd.csn.Set()
	return buf[4]
}

var (
	leds [5]gpio.Pin
	lcd  EVE
	u    *uart.Driver
)

func init() {
	system.Setup(clock.XTAL, clock.XTAL, true)
	rtcst.Setup(rtc.RTC0, 1)

	p0 := gpio.P0
	for i := range leds {
		led := p0.Pin(18 + i)
		led.Setup(gpio.ModeOut)
		leds[i] = led
	}
	lcd.csn = p0.Pin(12)
	lcd.irq = p0.Pin(14)
	lcd.pdn = p0.Pin(16)

	lcd.pdn.Setup(gpio.ModeOut)
	lcd.csn.Setup(gpio.ModeOut)
	lcd.csn.Set()
	lcd.irq.Setup(gpio.ModeIn)

	lcd.spi = spi.NewDriver(spi.SPI0)
	lcd.spi.P.StorePSEL(spi.SCK, p0.Pin(6))
	lcd.spi.P.StorePSEL(spi.MISO, p0.Pin(8))
	lcd.spi.P.StorePSEL(spi.MOSI, p0.Pin(10))
	lcd.spi.P.StoreFREQUENCY(spi.Freq8M)
	lcd.spi.Enable()
	rtos.IRQ(lcd.spi.P.NVIC()).Enable()

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
	/*
		in := make([]byte, 5)
		lcd.spi.WriteRead(nil, in)
		fmt.Printf("%d\r\n", in)
		lcd.spi.WriteRead([]byte{1, 2, 3}, in)
		fmt.Printf("%d\r\n", in)
		lcd.spi.WriteRead([]byte{1, 2, 3, 4, 5, 6, 7}, in)
		fmt.Printf("%d\r\n", in)
		fmt.Printf("%x\r\n", lcd.spi.WriteReadByte(0xAB))
	*/
	// Wakeup from POWERDOWN to STANDBY (PDn must be low min. 20 ms).
	delay.Millisec(20)
	lcd.pdn.Set()
	delay.Millisec(20) // Wait 20 ms for internal oscilator and PLL.

	// Wakeup from STANDBY to ACTIVE.
	lcd.Cmd(FT800_ACTIVE)

	// Select external 12 MHz oscilator as clock source.
	lcd.Cmd(FT800_CLKEXT)

	fmt.Printf("REG_ID=0x%x\r\n", lcd.Read8(REG_ID))
}

func spiISR() {
	lcd.spi.ISR()
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
