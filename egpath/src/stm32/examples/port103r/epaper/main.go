// WaveShare 1.54inch e-Paper Module (B) / Good Dispay GxGDEW0154Z04
//
// Connections:
//  DIN (blue)   <- PA7 (SPI1 MOSI)
//  CLK (yellow) <- PA5 (SPI1 SCK)
//  CS  (orange) <- PA6
//  DC  (green)  <- PA4
//  RST (white)  <- PA3
//  BUSY(violet) -> PA1
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
	"stm32/hal/system/timer/rtcst"
)

var (
	epd EPD

	led1, led2, led3 gpio.Pin
)

func init() {
	system.SetupPLL(8, 1, 72/8)
	rtcst.Setup(32768)

	// GPIO

	gpio.A.EnableClock(true)
	spiport, sck, mosi := gpio.A, gpio.Pin5, gpio.Pin7
	epd.cs = gpio.A.Pin(6)
	epd.dc = gpio.A.Pin(4)
	epd.rst = gpio.A.Pin(3)
	epd.busy = gpio.A.Pin(1)

	gpio.B.EnableClock(false)
	led1 = gpio.B.Pin(7)
	led2 = gpio.B.Pin(6)
	led3 = gpio.B.Pin(5)

	// SPI

	spiport.Setup(sck|mosi, &gpio.Config{Mode: gpio.Alt, Speed: gpio.High})
	d := dma.DMA1
	d.EnableClock(true)
	epd.spi = spi.NewDriver(spi.SPI1, nil, d.Channel(3, 0))
	epd.spi.P.EnableClock(true)
	epd.spi.P.SetConf(
		spi.Master | spi.MSBF | spi.CPOL0 | spi.CPHA0 |
			epd.spi.P.BR(10e6) | // 10 MHz max.
			spi.SoftSS | spi.ISSHigh,
	)
	epd.spi.P.SetWordSize(8)
	epd.spi.P.Enable()
	rtos.IRQ(irq.SPI1).Enable()
	rtos.IRQ(irq.DMA1_Channel3).Enable()

	// Controll

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.High}
	epd.cs.Setup(&cfg)
	epd.cs.Set()
	epd.dc.Setup(&cfg)
	cfg.Speed = gpio.Low
	epd.rst.Setup(&cfg)
	epd.busy.Setup(&gpio.Config{Mode: gpio.In})

	// LEDs
	cfg = gpio.Config{Mode: gpio.Out, Speed: gpio.Low}
	led1.Setup(&cfg)
	led2.Setup(&cfg)
	led3.Setup(&cfg)
}

func main() {
	delay.Millisec(100)
	spibus := epd.spi.P.Bus()
	baudrate := epd.spi.P.Baudrate(epd.spi.P.Conf())
	fmt.Printf(
		"\nSPI on %s (%d MHz). SPI speed: %d bps.\n\n",
		spibus, spibus.Clock()/1e6, baudrate,
	)

	led1.Set()

	epd.Reset()
	epd.Cmd(SetPower)
	epd.Write([]byte{7, 0, 8, 0})
	epd.Cmd(BoosterSoftStart)
	epd.Write([]byte{7, 7, 7})
	epd.Cmd(PowerOn)
	epd.Wait()

	epd.Cmd(SetPanel)
	epd.WriteByte(0xCF)
	epd.Cmd(SetVcomAndDataInt)
	epd.WriteByte(0x17)
	epd.Cmd(SetPLLControl)
	epd.WriteByte(0x39)
	epd.Cmd(SetResolution)
	epd.Write([]byte{0xC8, 0, 0xC8})
	epd.Cmd(SetVCMDC)
	epd.WriteByte(0xE)

	epd.Cmd(SetVcomLUT)
	epd.Write(lut_vcom0[:])
	epd.Cmd(SetWhiteLUT)
	epd.Write(lut_w[:])
	epd.Cmd(SetBlackLUT)
	epd.Write(lut_b[:])
	epd.Cmd(SetGray1LUT)
	epd.Write(lut_g1[:])
	epd.Cmd(SetGray2LUT)
	epd.Write(lut_g2[:])

	epd.Cmd(SetVcomRedLUT)
	epd.Write(lut_vcom1[:])
	epd.Cmd(SetRed0LUT)
	epd.Write(lut_red0[:])
	epd.Cmd(SetRed1LUT)
	epd.Write(lut_red1[:])

	led2.Set()

	epd.Wait()
	epd.Cmd(DisplayStartTx)
	for i, n := 0, 200*200/4; i < n; i++ {
		if i > n/3 {
			epd.WriteByte(0xFF)
		} else {
			epd.WriteByte(0)
		}
	}
	epd.Cmd(DisplayStartTxRed)
	for i, n := 0, 200*200/8; i < n; i++ {
		if i > n*2/3 {
			epd.WriteByte(0)
		} else {
			epd.WriteByte(0xFF)
		}
	}
	epd.Cmd(DisplayRefresh)
	epd.Wait()

	led3.Set()
}

func epdSPIISR() {
	epd.spi.ISR()
}

func epdTxDMAISR() {
	epd.spi.DMAISR(epd.spi.TxDMA)
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTCAlarm: rtcst.ISR,

	irq.SPI1:          epdSPIISR,
	irq.DMA1_Channel3: epdTxDMAISR,
}
