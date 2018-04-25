package main

import (
	"rtos"

	"display/eve"

	"stm32/evedci"
	"stm32/hal/dma"
	"stm32/hal/exti"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/spi"
	"stm32/hal/system"
	"stm32/hal/system/timer/rtcst"
)

var dci *evedci.SPI

func init() {
	system.SetupPLL(8, 1, 72/8)
	rtcst.Setup(32768)

	// GPIO

	gpio.A.EnableClock(true)
	spiport, sck, miso, mosi := gpio.A, gpio.Pin5, gpio.Pin6, gpio.Pin7

	gpio.B.EnableClock(true)
	csn := gpio.B.Pin(0)
	irqn := gpio.B.Pin(1)
	pdn := gpio.B.Pin(10)

	// EVE control lines

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.High}
	pdn.Setup(&cfg)
	csn.Setup(&cfg)
	irqn.Setup(&gpio.Config{Mode: gpio.In})
	irqline := exti.Lines(irqn.Mask())
	irqline.Connect(irqn.Port())
	irqline.EnableFallTrig()
	irqline.EnableIRQ()
	rtos.IRQ(irq.EXTI1).Enable()

	// EVE SPI

	spiport.Setup(sck|mosi, &gpio.Config{Mode: gpio.Alt, Speed: gpio.High})
	spiport.Setup(miso, &gpio.Config{Mode: gpio.AltIn})
	d := dma.DMA1
	d.EnableClock(true)
	spidrv := spi.NewDriver(spi.SPI1, d.Channel(3, 0), d.Channel(2, 0))
	rtos.IRQ(irq.SPI1).Enable()
	rtos.IRQ(irq.DMA1_Channel2).Enable()
	rtos.IRQ(irq.DMA1_Channel3).Enable()

	dci = evedci.NewSPI(spidrv, csn, pdn)
}

func f(n int) int { return n * 16 }

func main() {
	lcd := eve.NewDriver(dci, 128)
	lcd.Setup(11e6)
	lcd.Init(&eve.Default800x480, nil)
	lcd.Setup(30e6)
	lcd.SetBacklight(64)

	var x, y int
	for {
		dl := lcd.DL(-1)
		dl.Clear(eve.CST)
		dl.Begin(eve.POINTS)
		dl.PointSize(f(150))
		dl.Vertex2f(f(400), f(240))
		dl.PointSize(f(100))
		dl.ColorRGB(150, 0, 200)
		dl.ColorA(128)
		dl.Vertex2f(f(x), f(y))
		dl.Display()
		lcd.SwapDL()
		for {
			x, y = lcd.TouchScreenXY()
			if x|y != -32768 {
				break
			}
			lcd.Wait(eve.INT_TOUCH)
			lcd.ClearIntFlags(eve.INT_TOUCH)
		}
	}
}

func lcdSPIISR() {
	dci.SPI().ISR()
}

func lcdRxDMAISR() {
	dci.SPI().DMAISR(dci.SPI().RxDMA())
}

func lcdTxDMAISR() {
	dci.SPI().DMAISR(dci.SPI().TxDMA())
}

func exti1ISR() {
	exti.L1.ClearPending()
	dci.ISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTCAlarm: rtcst.ISR,

	irq.SPI1:          lcdSPIISR,
	irq.DMA1_Channel2: lcdRxDMAISR,
	irq.DMA1_Channel3: lcdTxDMAISR,
	irq.EXTI1:         exti1ISR,
}
