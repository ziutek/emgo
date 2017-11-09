package main

import (
	"delay"
	"fmt"
	"rtos"

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/spi"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var lcd EVE

func init() {
	system.Setup96(8)
	systick.Setup(2e6)

	// GPIO

	gpio.A.EnableClock(true)
	spiport, sck, miso, mosi := gpio.A, gpio.Pin5, gpio.Pin6, gpio.Pin7
	lcd.pdn = gpio.A.Pin(9)

	gpio.B.EnableClock(true)
	lcd.csn = gpio.B.Pin(6)

	gpio.C.EnableClock(true)
	lcd.irq = gpio.C.Pin(7)

	// SPI

	spiport.Setup(sck|mosi, &gpio.Config{Mode: gpio.Alt, Speed: gpio.High})
	spiport.Setup(miso, &gpio.Config{Mode: gpio.AltIn})
	spiport.SetAltFunc(sck|miso|mosi, gpio.SPI1)
	d := dma.DMA2
	d.EnableClock(true)
	lcd.spi = spi.NewDriver(spi.SPI1, d.Channel(2, 3), d.Channel(3, 3))
	lcd.spi.P.EnableClock(true)
	lcd.spi.P.SetConf(
		spi.Master | spi.MSBF | spi.CPOL0 | spi.CPHA0 |
			lcd.spi.P.BR(11e6) | // Max 11 MHz before configure PCLK.
			spi.SoftSS | spi.ISSHigh,
	)
	lcd.spi.P.SetWordSize(8)
	lcd.spi.P.Enable()
	rtos.IRQ(irq.SPI1).Enable()
	rtos.IRQ(irq.DMA2_Stream2).Enable()
	rtos.IRQ(irq.DMA2_Stream3).Enable()

	// Controll

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.High}
	lcd.csn.Setup(&cfg)
	lcd.csn.Set()
	lcd.pdn.Setup(&cfg)
	lcd.irq.Setup(&gpio.Config{Mode: gpio.In})
}

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

func main() {
	delay.Millisec(200)
	spibus := lcd.spi.P.Bus()
	baudrate := lcd.spi.P.Baudrate(lcd.spi.P.Conf())
	fmt.Printf(
		"\nSPI on %s (%d MHz).\nSPI speed: %d bps.\n",
		spibus, spibus.Clock()/1e6, baudrate,
	)

	// Wakeup from POWERDOWN to STANDBY (PDn must be low min. 20 ms).
	lcd.pdn.Set()
	delay.Millisec(20) // Wait 20 ms for internal oscilator and PLL.

	// Wakeup from STANDBY to ACTIVE.
	lcd.Cmd(FT800_ACTIVE)

	// Select external 12 MHz oscilator as clock source..
	lcd.Cmd(FT800_CLKEXT)

	lcd.spi.P.SetConf(lcd.spi.P.Conf()&^spi.BR256 | lcd.spi.P.BR(30e6))

	fmt.Printf("SPI set to %d MHz\n", lcd.spi.P.Baudrate(lcd.spi.P.Conf()))
	fmt.Printf("REGID=0x%X\n", lcd.Read8(REG_ID))
}

func lcdSPIISR() {
	lcd.spi.ISR()
}

func lcdRxDMAISR() {
	lcd.spi.DMAISR(lcd.spi.RxDMA)
}

func lcdTxDMAISR() {
	lcd.spi.DMAISR(lcd.spi.TxDMA)
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.SPI1:         lcdSPIISR,
	irq.DMA2_Stream2: lcdRxDMAISR,
	irq.DMA2_Stream3: lcdTxDMAISR,
}
