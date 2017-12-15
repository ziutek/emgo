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

	const testLen = 120

	// Create display list directly in RAM_DL.

	n := 800 // Greather numbers can exceed the render limit (8192 pixel/line).
	fmt.Printf("\n%d bitmaps:", n)

	addr := 0
	lcd.W(addr).Write(LenaFaceRGB[:])
	addr += len(LenaFaceRGB)

	t := rtos.Nanosec()
	for i := 0; i < testLen; i++ {
		dl := lcd.DL(-1)
		dl.BitmapLayout(eve.RGB565, 80, 40)
		dl.BitmapSize(eve.DEFAULT, 40, 40)
		dl.Clear(eve.CST)
		dl.Begin(eve.BITMAPS)
		for k := 0; k < n; k++ {
			v := rnd.Uint32()
			c := v&0xFFFFFF | 0x808080
			x := int(v>>12) % lcd.Width()
			y := int(v>>23) % lcd.Height()
			dl.ColorRGB(c)
			dl.Vertex2f((x-20)*16, (y-20)*16)
		}
		dl.Display()
		lcd.SwapDL()
	}
	fmt.Printf(" %.2f fps.\n", testLen*1e9/float32(rtos.Nanosec()-t))
	delay.Millisec(1000)

	// Create display list using Graphics Engine co-processor.

	ge := lcd.GE(-1)
	ge.DLStart()
	ge.BitmapHandle(1)

	fmt.Printf("Loading JPEG image...")

	t = rtos.Nanosec()
	ge.LoadImage(addr, eve.OPT_RGB565)
	img := GopherJPEG[:]
	for len(img) > 0 {
		n := lcd.CmdSpace()
		if n > len(img) {
			n = len(img)
		}
		lcd.GE(-1).Write(img[:n])
		img = img[n:]
	}
	lcd.Wait(eve.INT_CMDEMPTY)

	t = rtos.Nanosec() - t
	fmt.Printf(
		" done (%d B / %d ms = %d B/s).\n",
		len(GopherJPEG), t/1e6, int64(len(GopherJPEG))*1e9/t,
	)

	ge = lcd.GE(-1)
	ge.Clear(eve.CST)
	ge.Begin(eve.BITMAPS)
	ge.Vertex2f(0, 0)
	ge.BitmapHandle(0)
	ge.Vertex2f((lcd.Width()-40)*16, 0)
	ge.Button(300, 110, 140, 40, 23, 0, "Push me!")
	ge.Display()
	ge.Swap()
	lcd.Wait(eve.INT_CMDEMPTY)

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
