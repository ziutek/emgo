package main

import (
	"delay"
	"fmt"
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
	//rtos.IRQ(irq.EXTI9_5).Enable()

	dci = evedci.NewSPI(spidrv, csn, pdn, irqline)
}

const (
	lcdWidth   = 480 // Active width of LCD display
	lcdHeight  = 272 // Active height of LCD display
	lcdHcycle  = 548 // Total number of clocks per line
	lcdHoffset = 43  // Start of active line
	lcdHsync0  = 0   // Start of horizontal sync pulse
	lcdHsync1  = 41  // End of horizontal sync pulse
	lcdVcycle  = 292 // Total number of lines per screen
	lcdVoffset = 12  // Start of active screen
	lcdVsync0  = 0   // Start of vertical sync pulse
	lcdVsync1  = 10  // End of vertical sync pulse
	lcdPclk    = 5   // Pixel Clock
	lcdPclkpol = 1   // Define active edge of PCLK
)

func main() {
	delay.Millisec(200)
	spibus := dci.SPI().P.Bus()
	baudrate := dci.SPI().P.Baudrate(dci.SPI().P.Conf())
	fmt.Printf(
		"\nSPI on %s (%d MHz).\nSPI speed: %d bps.\n",
		spibus, spibus.Clock()/1e6, baudrate,
	)

	// Wakeup from POWERDOWN to STANDBY (PDN must be low min. 20 ms).
	dci.PDN().Set()
	delay.Millisec(20) // Wait 20 ms for internal oscilator and PLL.

	lcd := eve.NewDriver(dci, 32)

	fmt.Print("Init... ")

	// Wakeup from STANDBY to ACTIVE.
	lcd.Cmd(ft80.ACTIVE, 0)

	// Select external 12 MHz oscilator as clock source.
	lcd.Cmd(ft80.CLKEXT, 0)

	if lcd.Reader(ft80.REG_ID).ReadByte() != 0x7c {
		fmt.Printf("Not EVE controller.\n")
		return
	}
	if lcd.Reader(ft80.ROM_CHIPID).ReadWord32() != 0x10008 {
		fmt.Printf("Not FT800 controller.\n")
		return
	}

	check(lcd.Err(false))

	fmt.Print("Configure WQVGA (480x272) display...")

	lcd.Writer(ft80.REG_PWM_DUTY).Write32(0)

	lcd.Writer(ft80.REG_PCLK_POL).Write32(
		lcdPclkpol, // REG_PCLK_POL
		0,          // REG_PCLK
	)

	lcd.Writer(ft80.REG_HCYCLE).Write32(
		lcdHcycle,  // REG_HCYCLE
		lcdHoffset, // REG_HOFFSET
		lcdWidth,   // REG_HSIZE
		lcdHsync0,  // REG_HSYNC0
		lcdHsync1,  // REG_HSYNC1
		lcdVcycle,  // REG_VCYCLE
		lcdVoffset, // REG_VOFFSET
		lcdHeight,  // REG_VSIZE
		lcdVsync0,  // REG_VSYNC0
		lcdVsync1,  // REG_VSYNC1
	)

	check(lcd.Err(false))

	fmt.Print("Write initial display list and enable display...")

	lcd.Writer(ft80.RAM_DL).Write32(
		ft80.DL_CLEAR_RGB,
		ft80.DL_CLEAR|ft80.CLR_COL|ft80.CLR_STN|ft80.CLR_TAG,
		ft80.DL_DISPLAY,
	)

	lcd.Writer(ft80.REG_DLSWAP).Write32(ft80.DLSWAP_FRAME)

	gpio := lcd.Reader(ft80.REG_GPIO).ReadWord32()
	lcd.Writer(ft80.REG_GPIO).Write32(gpio | 0x80)
	lcd.Writer(ft80.REG_PCLK).Write32(lcdPclk)

	check(lcd.Err(false))

	dci.SPI().P.SetConf(dci.SPI().P.Conf()&^spi.BR256 | dci.SPI().P.BR(30e6))
	fmt.Printf("SPI set to %d Hz\n", dci.SPI().P.Baudrate(dci.SPI().P.Conf()))

	lcd.Writer(ft80.REG_PWM_DUTY).Write32(100)

	fmt.Print("Testing DL...")

	lcd.Writer(ft80.RAM_DL).Write32(
		ft80.DL_CLEAR_RGB,
		ft80.DL_CLEAR|ft80.CLR_COL|ft80.CLR_STN|ft80.CLR_TAG,
		ft80.DL_BEGIN|ft80.POINTS,
		ft80.DL_COLOR_RGB|0xa161f4,
		ft80.DL_POINT_SIZE|100*16,
		ft80.DL_VERTEX2F|200*16<<15|100*16,
		ft80.DL_COLOR_RGB|0xffff00,
		ft80.DL_POINT_SIZE|50*16,
		ft80.DL_VERTEX2F|300*16<<15|200*16,
		ft80.DL_DISPLAY,
	)

	for {
		lcd.Writer(ft80.REG_DLSWAP).Write32(ft80.DLSWAP_FRAME)
		check(lcd.Err(false))
		delay.Millisec(500)
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
