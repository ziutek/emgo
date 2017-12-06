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
	system.Setup96(8)
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
	d := dma.DMA2
	d.EnableClock(true)
	spidrv := spi.NewDriver(spi.SPI1, d.Channel(2, 3), d.Channel(3, 3))
	spidrv.P.EnableClock(true)
	rtos.IRQ(irq.SPI1).Enable()
	rtos.IRQ(irq.DMA2_Stream2).Enable()
	rtos.IRQ(irq.DMA2_Stream3).Enable()

	// EVE control lines

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.High}
	pdn.Setup(&cfg)
	csn.Setup(&cfg)
	irqn.Setup(&gpio.Config{Mode: gpio.In})
	irqline := exti.Lines(irqn.Mask())
	irqline.Connect(irqn.Port())
	//rtos.IRQ(irq.EXTI9_5).Enable()

	dci = evedci.NewSPI(spidrv, csn, pdn, irqline)
}

func curFreq(lcd *eve.Driver) uint32 {
	clk1 := lcd.StartR(ft80.REG_CLOCK).ReadWord32()
	t1 := rtos.Nanosec()
	delay.Millisec(8)
	clk2 := lcd.StartR(ft80.REG_CLOCK).ReadWord32()
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
		// Triming if internal oscilator is used.
		for trim := uint32(0); trim <= 31; trim++ {
			lcd.StartW(ft80.REG_TRIM).Write32(trim)
			if f := curFreq(lcd); f > 47040000 {
				lcd.StartW(ft80.REG_FREQUENCY).Write32(f)
				break
			}
		}
	*/

	// Select external 12 MHz oscilator as clock source.
	lcd.Cmd(ft80.CLKEXT, 0)

	if lcd.StartR(ft80.REG_ID).ReadByte() != 0x7c {
		fmt.Printf("Not EVE controller.\n")
		return
	}
	if lcd.StartR(ft80.ROM_CHIPID).ReadWord32() != 0x10008 {
		fmt.Printf("Not FT800 controller.\n")
		return
	}
	check(lcd.Err(false))

	printFreq(lcd)

	fmt.Print("Configure WQVGA (480x272) display:")

	lcd.StartW(ft80.REG_PWM_DUTY).Write32(0)

	const pclkDiv = 5 // Pixel Clock divider: pclk = mainClk / pclkDiv.

	lcd.StartW(ft80.REG_PCLK_POL).Write32(
		1, // REG_PCLK_POL (define active edge of PCLK)
		0, // REG_PCLK (temporary disable PCLK)
	)
	lcd.StartW(ft80.REG_HCYCLE).Write32(
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

	// Refresh rate: pclk/(hcycle*vcycle) = 48 MHz/5/(548*292) = 59.99 Hz.

	check(lcd.Err(false))

	fmt.Print("Write initial display list and enable display:")

	dl := lcd.StartDL(ft80.RAM_DL)
	dl.ClearColorRGB(0, 0, 0)
	dl.Clear(eve.CST)
	dl.Display()

	// Alternative, method:
	//
	//  lcd.StartW(ft80.RAM_DL).Write32(
	//  	eve.CLEAR_COLOR_RGB,
	//  	eve.CLEAR|eve.CST,
	//  	eve.DISPLAY,
	//  )

	lcd.StartW(ft80.REG_DLSWAP).Write32(eve.DLSWAP_FRAME)

	gpio := lcd.StartR(ft80.REG_GPIO).ReadWord32()
	lcd.StartW(ft80.REG_GPIO).Write32(gpio | 0x80)
	lcd.StartW(ft80.REG_PCLK).Write32(pclkDiv) // Enable PCLK.
	check(lcd.Err(false))
	printFreq(lcd)

	delay.Millisec(20) // Wait for new main clock.

	printFreq(lcd)

	fmt.Println(dci.SPI().P.BR(12e6) >> 3)
	dci.SPI().P.SetConf(dci.SPI().P.Conf()&^spi.BR256 | dci.SPI().P.BR(12e6))
	fmt.Printf("SPI set to %d Hz\n", dci.SPI().P.Baudrate(dci.SPI().P.Conf()))

	printFreq(lcd)

	lcd.StartW(ft80.REG_PWM_DUTY).Write32(100)

	fmt.Print("Points:")

	dl = lcd.StartDL(ft80.RAM_DL)
	dl.Clear(eve.CST)
	dl.Begin(eve.POINTS)
	dl.ColorRGB(161, 244, 97)
	dl.PointSize(100 * 16)
	dl.Vertex2F(200*16, 100*16)
	dl.ColorRGB(255, 0, 255)
	dl.PointSize(50 * 16)
	dl.Vertex2F(300*16, 200*16)
	dl.Display()

	lcd.StartW(ft80.REG_DLSWAP).Write32(eve.DLSWAP_FRAME)
	check(lcd.Err(false))

	delay.Millisec(1000)

	fmt.Print("Load bitmap:")

	w := lcd.StartW(ft80.RAM_G)
	w.Write(LenaFace[:])

	check(lcd.Err(false))

	fmt.Print("Draw bitmap:")

	var rnd rand.XorShift64
	rnd.Seed(1)

	dla := ft80.RAM_DL

	dl = lcd.StartDL(dla)
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

	dla += dl.Close()
	check(lcd.Err(false))

	lcd.StartW(ft80.REG_CMD_DL).WriteInt(dla)

	n := lcd.StartR(ft80.REG_CMD_WRITE).ReadInt()
	ge := lcd.StartGE(ft80.RAM_CMD + n)
	ge.Button(170, 110, 140, 40, 23, 0, "Push me!")
	ge.Display()
	ge.Swap()
	n += ge.Close()
	lcd.StartW(ft80.REG_CMD_WRITE).WriteInt(n)

	check(lcd.Err(false))

	for {
		delay.Millisec(1000)
		lcd.StartW(ft80.REG_DLSWAP).Write32(eve.DLSWAP_FRAME)
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
	irq.SPI1:         lcdSPIISR,
	irq.DMA2_Stream2: lcdRxDMAISR,
	irq.DMA2_Stream3: lcdTxDMAISR,
}
