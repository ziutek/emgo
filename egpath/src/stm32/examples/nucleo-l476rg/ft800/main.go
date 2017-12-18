// This example demonstrates usage of FTDI EVE based displays.
//
// It seems that FT800CB-HY50B display is unstable with fast SPI. If you have
// problems please reduce SPI speed or better desolder U1 and U2 (74LCX125
// buffers) and short the U1:2-3,5-6,11-2, U2:2-3,5-6 traces.
package main

import (
	"delay"
	"fmt"
	"math/rand"
	"rtos"

	"display/eve"
	"display/eve/ft80"

	"stm32/evedci"

	"stm32/hal/dma"
	"stm32/hal/exti"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/spi"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var dci *evedci.SPI

func init() {
	system.Setup80(0, 0)
	systick.Setup(2e6)

	// GPIO

	gpio.A.EnableClock(true)
	spiport, sck, miso, mosi := gpio.A, gpio.Pin5, gpio.Pin6, gpio.Pin7
	pdn := gpio.A.Pin(9)

	gpio.B.EnableClock(true)
	csn := gpio.B.Pin(6)

	gpio.C.EnableClock(true)
	irqn := gpio.C.Pin(7)

	// EVE SPI

	spiport.Setup(sck|mosi, &gpio.Config{Mode: gpio.Alt, Speed: gpio.High})
	spiport.Setup(miso, &gpio.Config{Mode: gpio.AltIn})
	spiport.SetAltFunc(sck|miso|mosi, gpio.SPI1)
	d := dma.DMA1
	d.EnableClock(true)
	rxdc, txdc := d.Channel(2, 0), d.Channel(3, 0)
	rxdc.SetRequest(dma.DMA1_SPI1)
	txdc.SetRequest(dma.DMA1_SPI1)
	spidrv := spi.NewDriver(spi.SPI1, rxdc, txdc)
	spidrv.P.EnableClock(true)
	rtos.IRQ(irq.SPI1).Enable()
	rtos.IRQ(irq.DMA1_Channel2).Enable()
	rtos.IRQ(irq.DMA1_Channel3).Enable()

	// EVE control lines

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.High}
	pdn.Setup(&cfg)
	csn.Setup(&cfg)
	irqn.Setup(&gpio.Config{Mode: gpio.In})
	irqline := exti.Lines(irqn.Mask())
	irqline.Connect(irqn.Port())
	irqline.EnableFallTrig()
	irqline.EnableIRQ()
	rtos.IRQ(irq.EXTI9_5).Enable()

	dci = evedci.NewSPI(spidrv, csn, pdn)
}

func curFreq(lcd *eve.Driver) uint32 {
	clk1 := lcd.ReadUint32(ft80.REG_CLOCK)
	t1 := rtos.Nanosec()
	delay.Millisec(8)
	clk2 := lcd.ReadUint32(ft80.REG_CLOCK)
	t2 := rtos.Nanosec()
	return uint32(int64(clk2-clk1) * 1e9 / (t2 - t1))
}

func main() {
	var rnd rand.XorShift64
	rnd.Seed(1)

	spibus := dci.SPI().P.Bus()
	fmt.Printf("\nSPI on %s (%d MHz).\n", spibus, spibus.Clock()/1e6)
	fmt.Printf("SPI speed: %d bps.\n", dci.SPI().P.Baudrate(dci.SPI().P.Conf()))

	lcd := eve.NewDriver(dci, 128)
	lcd.Init(&eve.Default480x272)

	fmt.Printf("EVE clock: %d Hz.\n", curFreq(lcd))
	dci.SetBaudrate(30e6)
	fmt.Printf("SPI speed: %d bps.\n", dci.SPI().P.Baudrate(dci.SPI().P.Conf()))

	lcd.SetBacklight(64)

	lcdW, lcdH := lcd.Width(), lcd.Height()

	ge := lcd.GE(-1)
	ge.DLStart()
	ge.Clear(eve.CST)
	ge.TextString(lcdW/2, lcdH/2, 31, eve.OPT_CENTER, "Hello world!")
	ge.Display()
	ge.Swap()
	lcd.Wait(eve.INT_CMDEMPTY)

	lcd.ClearIntFlags(eve.INT_TOUCH)
	lcd.Wait(eve.INT_TOUCH)

	lcd.Write(0, gopherMask[:])
	addr := (len(gopherMask) + 3) &^ 3

	ge.DLStart()
	ge.LoadImageBytes(addr, eve.OPT_RGB565, gopher[:])
	lcd.Wait(eve.INT_CMDEMPTY) // Ensure free space in RAM_CMD.
	ge.BitmapHandle(1)
	ge.BitmapLayout(eve.L1, 216/8, 251)
	ge.BitmapSize(eve.DEFAULT, 211, 251)
	ge.Clear(eve.CST)
	ge.Gradient(0, 0, 0x001155, 0, 271, 0x772200)
	ge.Begin(eve.BITMAPS)
	ge.ColorMask(eve.A)
	ge.Clear(eve.C)
	ge.BitmapHandle(1)
	ge.Vertex2f(31*16, 21*16)
	ge.ColorMask(eve.RGBA)
	ge.BlendFunc(eve.DST_ALPHA, eve.ONE_MINUS_DST_ALPHA)
	ge.BitmapHandle(0)
	ge.Vertex2f(0, 0)
	ge.Display()
	ge.Swap()
	lcd.Wait(eve.INT_CMDEMPTY)

	delay.Millisec(1e5)

	ge.DLStart()
	ge.Clear(eve.CST)
	addr = ge.Calibrate()
	lcd.Wait(eve.INT_CMDEMPTY)
	if lcd.ReadInt(addr) == 0 {
		fmt.Printf("Touch calibration failed!\n")
	}

	const buttonTag = 1

	for n := 0; ; n++ {
		tag := lcd.TouchTag()
		x, y := lcd.TouchScreenXY()
		ge.DLStart()
		ge.ClearColorRGB(0xc3a6f4)
		ge.Clear(eve.CST)
		ge.Gradient(0, 0, 0x0004ff, 0, 271, 0xe08484)
		ge.Text(lcdW-180, 20, 26, eve.DEFAULT)
		fmt.Fprintf(&ge, "x=%d y=%d tag=%d\000", x, y, tag)
		ge.Align32()
		ge.TextString(lcdW/2, lcdH/2, 31, eve.OPT_CENTER, "Hello World!")
		ge.Begin(eve.RECTS)
		ge.ColorA(128)
		ge.ColorRGB(0xff8000)
		ge.Vertex2ii(260, 100, 0, 0)
		ge.Vertex2ii(360, 200, 0, 0)
		ge.ColorRGB(0x0080ff)
		ge.Vertex2ii(300, 160, 0, 0)
		ge.Vertex2ii(400, 260, 0, 0)
		ge.ColorRGB(0xffffff)
		ge.ColorA(200)
		ge.Clock(60, 60, 50, eve.OPT_NOBACK, 23, 49, n/60, n%60*1000/60)
		ge.ColorA(255)
		ge.Tag(buttonTag)
		buttonFont := byte(27)
		buttonStyle := uint16(eve.DEFAULT)
		if tag == buttonTag {
			buttonFont--
			buttonStyle |= eve.OPT_FLAT
			ge.TextString(20, lcdH-90, 29, eve.DEFAULT, "Thanks!")
		}
		ge.ButtonString(20, lcdH-50, 100, 32, buttonFont, buttonStyle, "Push me!")
		ge.Display()
		ge.Swap()
		lcd.Wait(eve.INT_CMDEMPTY) // Wait for end of Swap (next frame).
	}

	fmt.Printf("End.\n")
}

/*ge := lcd.GE(ft80.RAM_CMD + n)
ge.Clear(eve.CST)
ge.Calibrate()
n += ge.Close() + 4
lcd.WriteInt(ft80.REG_CMD_WRITE, n&4095)*/

func lcdSPIISR() {
	dci.SPI().ISR()
}

func lcdRxDMAISR() {
	dci.SPI().DMAISR(dci.SPI().RxDMA)
}

func lcdTxDMAISR() {
	dci.SPI().DMAISR(dci.SPI().TxDMA)
}

func exti9_5ISR() {
	exti.Pending().ClearPending()
	dci.ISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.SPI1:          lcdSPIISR,
	irq.DMA1_Channel2: lcdRxDMAISR,
	irq.DMA1_Channel3: lcdTxDMAISR,
	irq.EXTI9_5:       exti9_5ISR,
}
