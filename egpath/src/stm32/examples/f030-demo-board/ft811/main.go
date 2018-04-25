package main

import (
	"delay"
	"rtos"

	"display/eve"
	"display/eve/ft81"

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
	dci.Setup(11e6)
}

func lcdHostCmd(cmd eve.HostCmd, param byte) {
	dci.Write([]byte{byte(cmd), param, 0})
	dci.End()
}

func lcdWriteUint32(addr int, val uint32) int {
	dci.Write([]byte{
		1<<7 | byte(addr>>16), byte(addr >> 8), byte(addr),
		byte(val), byte(val >> 8), byte(val >> 16), byte(val >> 24),
	})
	dci.End()
	return addr + 4
}

func lcdReadUint32(addr int) uint32 {
	buf := [5]byte{byte(addr >> 16), byte(addr >> 8), byte(addr)}
	dci.Write(buf[:3])
	dci.Read(buf[:5])
	dci.End()
	return uint32(buf[1]) | uint32(buf[2])<<8 | uint32(buf[3])<<16 |
		uint32(buf[4])<<24
}

func lcdWait(flags uint32) {
	for lcdReadUint32(ft81.REG_INT_FLAGS)&flags == 0 {
		<-dci.IRQ()
	}
}

func f(n int16) uint32 { return uint32(n) * 16 & 0x7FFF }

func main() {
	dci.SetPDN(0)
	delay.Millisec(20)
	dci.SetPDN(1)
	delay.Millisec(20)

	lcdHostCmd(ft81.CLKEXT, 0)
	lcdHostCmd(ft81.ACTIVE, 0)
	delay.Millisec(300)

	cfg := &eve.Default800x480
	lcdWriteUint32(ft81.REG_HCYCLE, uint32(cfg.Hcycle))
	lcdWriteUint32(ft81.REG_HOFFSET, uint32(cfg.Hoffset))
	lcdWriteUint32(ft81.REG_HSIZE, uint32(cfg.Hsize))
	lcdWriteUint32(ft81.REG_HSYNC0, uint32(cfg.Hsync0))
	lcdWriteUint32(ft81.REG_HSYNC1, uint32(cfg.Hsync1))
	lcdWriteUint32(ft81.REG_VCYCLE, uint32(cfg.Vcycle))
	lcdWriteUint32(ft81.REG_VOFFSET, uint32(cfg.Voffset))
	lcdWriteUint32(ft81.REG_VSIZE, uint32(cfg.Vsize))
	lcdWriteUint32(ft81.REG_VSYNC0, uint32(cfg.Vsync0))
	lcdWriteUint32(ft81.REG_VSYNC1, uint32(cfg.Vsync1))
	lcdWriteUint32(ft81.REG_PCLK_POL, uint32(cfg.ClkPol))

	lcdWriteUint32(ft81.REG_GPIO, 0x80)
	lcdWriteUint32(ft81.REG_PCLK, 5)
	lcdWriteUint32(ft81.REG_INT_EN, 1)

	dci.Setup(30e6)

	var x, y int16
	for {
		a := lcdWriteUint32(ft81.RAM_DL, eve.CLEAR|eve.CST)
		a = lcdWriteUint32(a, eve.BEGIN|eve.POINTS)
		a = lcdWriteUint32(a, eve.POINT_SIZE|f(150))
		a = lcdWriteUint32(a, eve.VERTEX2F|f(400)<<15|f(240))
		a = lcdWriteUint32(a, eve.POINT_SIZE|f(100))
		a = lcdWriteUint32(a, eve.COLOR_RGB|0x9600C8)
		a = lcdWriteUint32(a, eve.COLOR_A|128)
		a = lcdWriteUint32(a, eve.VERTEX2F|f(x)<<15|f(y))
		a = lcdWriteUint32(a, eve.DISPLAY)
		lcdWriteUint32(ft81.REG_DLSWAP, eve.DLSWAP_FRAME)
		lcdWait(eve.INT_SWAP)
		for {
			xy := lcdReadUint32(ft81.REG_TOUCH_SCREEN_XY)
			if xy != 0x80008000 {
				x, y = int16(xy>>16), int16(xy)
				break
			}
			lcdWait(eve.INT_TOUCH)
		}
	}
}

func lcdSPIISR() {
	dci.SPI().ISR()
}

func lcdDMAISR() {
	p := dci.SPI()
	p.DMAISR(p.RxDMA())
	p.DMAISR(p.TxDMA())
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
