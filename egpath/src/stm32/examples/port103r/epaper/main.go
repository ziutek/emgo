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
	"bytes"
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
	epd.spi = spi.NewDriver(spi.SPI1, d.Channel(3, 0), nil)
	epd.spi.Periph().EnableClock(true)
	epd.spi.Periph().SetConf(
		spi.Master | spi.MSBF | spi.CPOL0 | spi.CPHA0 |
			epd.spi.Periph().BR(10e6) | // 10 MHz max.
			spi.SoftSS | spi.ISSHigh,
	)
	epd.spi.Periph().Enable()
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
	delay.Millisec(100) // For SWO output.

	spip := epd.spi.Periph()
	fmt.Printf("\nSPI on %s (%d MHz).\n", spip.Bus(), spip.Bus().Clock()/1e6)
	fmt.Printf("SPI speed: %d bps.\n", spip.Baudrate(spip.Conf()))

	led1.Set()

	epd.Reset()

	epd.Begin()

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

	const (
		w = 200
		h = 200
	)

	black := make([]byte, w*h/4)
	bytes.Fill(black, 0xFF)
	red := make([]byte, w*h/8)
	bytes.Fill(red, 0xFF)

	for i := 0; i < w; i++ {
		x, y := i, i
		o := y*(w/4) + x/4
		black[o] ^= (3 << 6) >> uint(x&3*2)
	}
	for i := 0; i < w-10; i++ {
		x, y := i+10, i
		o := y*(w/8) + x/8
		red[o] ^= (1 << 7) >> uint(x&7)
	}

	epd.Wait()
	epd.Cmd(DisplayStartTx)
	epd.Write(black)
	epd.Cmd(DisplayStartTxRed)
	epd.Write(red)
	epd.Cmd(DisplayRefresh)

	epd.End()

	epd.Wait()

	led3.Set()
}

func epdSPIISR() {
	epd.spi.ISR()
}

func epdTxDMAISR() {
	epd.spi.DMAISR(epd.spi.TxDMA())
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTCAlarm: rtcst.ISR,

	irq.SPI1:          epdSPIISR,
	irq.DMA1_Channel3: epdTxDMAISR,
}
