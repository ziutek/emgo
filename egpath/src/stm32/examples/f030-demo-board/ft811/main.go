package main

import (
	"rtos"

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
	system.SetupPLL(8, 1, 48/8)
	systick.Setup(2e6)

	// GPIO

	gpio.A.EnableClock(true)
	spiport, sck, miso, mosi := gpio.A, gpio.Pin5, gpio.Pin6, gpio.Pin7
	irqn := gpio.A.Pin(9)
	pdn := gpio.A.Pin(10)

	gpio.B.EnableClock(true)
	csn := gpio.B.Pin(1)

	// EVE control lines

	cfg := gpio.Config{Mode: gpio.Out, Speed: gpio.High}
	pdn.Setup(&cfg)
	csn.Setup(&cfg)
	irqn.Setup(&gpio.Config{Mode: gpio.In})
	irqline := exti.Lines(irqn.Mask())
	irqline.Connect(irqn.Port())
	irqline.EnableFallTrig()
	irqline.EnableIRQ()
	rtos.IRQ(irq.EXTI4_15).Enable()

	// EVE SPI

	spiport.Setup(sck|mosi, &gpio.Config{Mode: gpio.Alt, Speed: gpio.High})
	spiport.Setup(miso, &gpio.Config{Mode: gpio.AltIn})
	spiport.SetAltFunc(sck|miso|mosi, gpio.SPI1)
	d := dma.DMA1
	d.EnableClock(true)
	spidrv := spi.NewDriver(spi.SPI1, d.Channel(3, 0), d.Channel(2, 0))
	rtos.IRQ(irq.SPI1).Enable()
	rtos.IRQ(irq.DMA1_Channel2_3).Enable()

	dci = evedci.NewSPI(spidrv, csn, pdn)
}

func main() {
}

func lcdSPIISR() {
	dci.SPI().ISR()
}

func lcdDMAISR() {
	dci.SPI().DMAISR(dci.SPI().RxDMA())
	dci.SPI().DMAISR(dci.SPI().TxDMA())
}

func exti4_15ISR() {
	pending := exti.Pending()
	pending &= exti.L4 | exti.L5 | exti.L6 | exti.L7 | exti.L8 | exti.L9 |
		exti.L10 | exti.L11 | exti.L12 | exti.L13 | exti.L14 | exti.L15
	pending.ClearPending()
	if pending&exti.L9 != 0 {
		dci.ISR()
	}
}

//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.SPI1:            lcdSPIISR,
	irq.DMA1_Channel2_3: lcdDMAISR,
	irq.EXTI4_15:        exti4_15ISR,
}
