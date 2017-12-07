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
	system.SetupPLL(8, 1, 72/8)
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
	spidrv := spi.NewDriver(spi.SPI1, d.Channel(2, 0), d.Channel(3, 0))
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
	rtos.IRQ(irq.EXTI9_5).Enable()

	dci = evedci.NewSPI(spidrv, csn, pdn, irqline)
}

func curFreq(lcd *eve.Driver) uint32 {
	clk1 := lcd.ReadUint32(ft80.REG_CLOCK)
	t1 := rtos.Nanosec()
	delay.Millisec(8)
	clk2 := lcd.ReadUint32(ft80.REG_CLOCK)
	t2 := rtos.Nanosec()
	return uint32(int64(clk2-clk1) * 1e9 / (t2 - t1))
}

func printFreq(lcd *eve.Driver) {
	fmt.Printf("FT800 clock: %d Hz\n", curFreq(lcd))
}

func main() {
	delay.Millisec(200)
	spibus := dci.SPI().P.Bus()
	baudrate := dci.SPI().P.Baudrate(dci.SPI().P.Conf())
	fmt.Printf(
		"\nSPI on %s (%d MHz).\nSPI speed: %d bps.\n",
		spibus, spibus.Clock()/1e6, baudrate,
	)

	// Wakeup from POWERDOWN to STANDBY.
	dci.PDN().Clear()
	delay.Millisec(20)
	dci.PDN().Set()
	delay.Millisec(20) // Wait 20 ms for internal oscilator and PLL.

	lcd := eve.NewDriver(dci, 128)

	fmt.Print("Init:")

	// Wakeup from STANDBY to ACTIVE.
	lcd.Cmd(ft80.ACTIVE, 0)

	/*
		// Simple triming algorithm if internal oscilator is used.
		for trim := uint32(0); trim <= 31; trim++ {
			lcd.W(ft80.REG_TRIM).W32(trim)
			if f := curFreq(lcd); f > 47040000 {
				lcd.W(ft80.REG_FREQUENCY).W32(f)
				break
			}
		}
	*/

	// Select external 12 MHz oscilator as clock source.
	lcd.Cmd(ft80.CLKEXT, 0)

	if lcd.ReadByte(ft80.REG_ID) != 0x7c {
		fmt.Printf("Not EVE controller.\n")
		return
	}
	if lcd.ReadUint32(ft80.ROM_CHIPID) != 0x10008 {
		fmt.Printf("Not FT800 controller.\n")
		return
	}
	check(lcd.Err(false))

	printFreq(lcd)

	fmt.Print("Configure WQVGA (480x272) display:")

	lcd.WriteByte(ft80.REG_PWM_DUTY, 0)

	const pclkDiv = 5 // Pixel Clock divider: pclk = mainClk / pclkDiv.

	// Refresh rate: pclk/(hcycle*vcycle) = 48 MHz/5/(548*292) = 59.99 Hz.

	lcd.W(ft80.REG_CSPREAD).Write32(
		1, // REG_CSPREAD (color signals spread, reduces EM noise)
		1, // REG_PCLK_POL (define active edge of PCLK)
		0, // REG_PCLK (temporary disable PCLK)
	)
	lcd.W(ft80.REG_HCYCLE).Write32(
		548, // REG_HCYCLE  (total number of clocks per line)
		43,  // REG_HOFFSET (tart of active line)
		480, // REG_HSIZE   (active width of LCD display)
		0,   // REG_HSYNC0  (start of horizontal sync pulse)
		41,  // REG_HSYNC1  (end of horizontal sync pulse)
		292, // REG_VCYCLE  (total number of lines per screen)
		12,  // REG_VOFFSET (start of active screen)
		272, // REG_VSIZE   (active height of LCD display)
		0,   // REG_VSYNC0  (start of vertical sync pulse)
		10,  // REG_VSYNC1  (end of vertical sync pulse)
	)

	check(lcd.Err(false))

	fmt.Print("Write initial display list and enable display:")

	dl := lcd.DL(ft80.RAM_DL)
	dl.ClearColorRGB(0, 0, 0)
	dl.Clear(eve.CST)
	dl.Display()

	// Alternative, method:
	//
	//  lcd.W(ft80.RAM_DL).Write32(
	//  	eve.CLEAR_COLOR_RGB,
	//  	eve.CLEAR|eve.CST,
	//  	eve.DISPLAY,
	//  )

	lcd.WriteByte(ft80.REG_DLSWAP, eve.DLSWAP_FRAME)

	gpio := lcd.ReadByte(ft80.REG_GPIO)
	lcd.WriteByte(ft80.REG_GPIO, gpio|0x80)
	lcd.WriteByte(ft80.REG_PCLK, pclkDiv) // Enable PCLK.
	check(lcd.Err(false))
	printFreq(lcd)

	delay.Millisec(20) // Wait for new main clock.

	printFreq(lcd)

	// FT800CB-HY50B display is unstable with fast SPI and VCC <= 3.3V. If you
	// have problems please comment the line bellow or better desolder U1 and U2
	// (74LCX125 buffers) and short the U1:2-3,5-6,11-2, U2:2-3,5-6 pins.
	dci.SPI().P.SetConf(dci.SPI().P.Conf()&^spi.BR256 | dci.SPI().P.BR(30e6))
	fmt.Printf("SPI set to %d Hz\n", dci.SPI().P.Baudrate(dci.SPI().P.Conf()))

	lcd.WriteByte(ft80.REG_PWM_DUTY, 64)

	fmt.Print("Draw two points:")

	dl = lcd.DL(ft80.RAM_DL)
	dl.Clear(eve.CST)
	dl.Begin(eve.POINTS)
	dl.ColorRGB(161, 244, 97)
	dl.PointSize(100 * 16)
	dl.Vertex2F(200*16, 100*16)
	dl.ColorRGB(255, 0, 255)
	dl.PointSize(50 * 16)
	dl.Vertex2F(300*16, 200*16)
	dl.Display()
	lcd.WriteByte(ft80.REG_DLSWAP, eve.DLSWAP_FRAME)
	check(lcd.Err(false))

	delay.Millisec(1000)

	fmt.Print("Load bitmap:")
	lcd.W(ft80.RAM_G).Write(LenaFace[:])
	check(lcd.Err(false))

	fmt.Print("Draw widgets on top of 1000 bitmaps:")

	var rnd rand.XorShift64
	rnd.Seed(1)

	addr := ft80.RAM_DL
	dl = lcd.DL(addr)
	dl.BitmapHandle(1)
	dl.BitmapSource(ft80.RAM_G)
	dl.BitmapLayout(eve.RGB565, 80, 40)
	dl.BitmapSize(0, 40, 40)
	dl.Clear(eve.CST)
	dl.Begin(eve.BITMAPS)
	dl.ColorA(255)
	dl.BitmapHandle(1)
	for i := 0; i < 1000; i++ {
		v := rnd.Uint64()
		vl := uint32(v)
		vh := uint32(v >> 32)
		dl.Vertex2F(int((vl%480-20)*16), int((vh%272-20)*16))
	}
	addr += dl.Close()

	lcd.WriteInt(ft80.REG_CMD_DL, addr)

	n := lcd.ReadInt(ft80.REG_CMD_WRITE)
	ge := lcd.GE(ft80.RAM_CMD + n)
	ge.Button(170, 110, 140, 40, 23, 0, "Push me!")
	ge.Clock(440, 40, 30, 0, 21, 22, 42, 00)
	ge.Gauge(440, 232, 30, 0, 5, 5, 33, 100)
	ge.Keys(30, 242, 120, 20, 18, 0, "ABCDE")
	ge.Progress(180, 248, 100, 10, 0, 75, 100)
	ge.Scrollbar(10, 10, 100, 10, 0, 50, 25, 100)
	ge.Slider(10, 30, 100, 10, 0, 25, 100)
	ge.Dial(40, 80, 30, 0, 3000)
	ge.Toggle(25, 130, 30, 18, 0, true, "yes")
	ge.Display()
	ge.Swap()
	n += ge.Close()
	lcd.WriteInt(ft80.REG_CMD_WRITE, n)

	check(lcd.Err(false))

	for {
		delay.Millisec(1000)
		fmt.Print("Swap DL:")
		lcd.WriteByte(ft80.REG_DLSWAP, eve.DLSWAP_FRAME)
		check(lcd.Err(false))
	}
}



func check(err error) {
	if err == nil {
		fmt.Printf(" OK\n")
		return
	}
	fmt.Printf(" %v\n", err)
	for {
	}
}

func lcdSPIISR() {
	dci.SPI().ISR()
}

func lcdRxDMAISR() {
	dci.SPI().DMAISR(dci.SPI().RxDMA)
}

func lcdTxDMAISR() {
	dci.SPI().DMAISR(dci.SPI().TxDMA)
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.SPI1:          lcdSPIISR,
	irq.DMA1_Channel2: lcdRxDMAISR,
	irq.DMA1_Channel3: lcdTxDMAISR,
}
