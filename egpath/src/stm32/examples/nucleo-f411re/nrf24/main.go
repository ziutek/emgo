package main

import (
	"bufio"
	"delay"
	"fmt"
	"mmio"
	"rtos"
	"text/linewriter"

	"arch/cortexm/bitband"
	"arch/cortexm/debug/itm"

	"nrf24"

	"stm32/hal/dma"
	"stm32/hal/exti"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/spi"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
	"stm32/hal/usart"

	"stm32/hal/raw/rcc"
	"stm32/hal/raw/tim"
)

const dbg = itm.Port(0)

type nrfDCI struct {
	spi  *spi.Driver
	csn  bitband.Bit
	cet  *tim.TIM_Periph
	ocmn uint
	irq  exti.Lines
}

func (dci *nrfDCI) WriteRead(oi ...[]byte) (n int, err error) {
	dci.csn.Clear()
	dci.spi.WriteReadMany(oi...)
	dci.csn.Set()
	return n, dci.spi.Err()
}

func (dci *nrfDCI) SetCE(v int) error {
	switch v {
	case 0:
		dci.cet.CCMR2.Store(4 << dci.ocmn)
	case 1:
		dci.cet.CCMR2.Store(5 << dci.ocmn)
	case 2:
		dci.cet.CCMR2.Store(7 << dci.ocmn)
		dci.cet.CR1.Store(tim.OPM | tim.CEN)
	}
	return nil
}

func (dci *nrfDCI) Setup(spidrv *spi.Driver, csn bitband.Bit, pclk uint, cet *tim.TIM_Periph, ch int, irqline exti.Lines) {
	dci.spi = spidrv
	spidrv.P.SetConf(
		spi.Master | spi.MSBF | spi.CPOL0 | spi.CPHA0 |
			spidrv.P.BR(10e6) | // 10 MHz max.
			spi.SoftSS | spi.ISSHigh,
	)
	spidrv.P.Enable()

	dci.csn = csn
	csn.Set()

	dci.cet = cet
	if pclk != system.AHB.Clock() {
		pclk *= 2
	}
	cet.PSC.U16.Store(48e3 - 1) // 2e3 Hz
	cet.ARR.U32.Store(1e3 - 1)  // 1 Hz
	cet.UG().Set()
	var (
		ccr *mmio.U32
		cce tim.CCER_Bits
	)
	switch ch {
	case 1:
		ccr = &cet.CCR1.U32
		cce = tim.CC1E
		dci.ocmn = tim.OC1Mn
	case 2:
		ccr = &cet.CCR2.U32
		cce = tim.CC2E
		dci.ocmn = tim.OC2Mn
	case 3:
		ccr = &cet.CCR3.U32
		cce = tim.CC3E
		dci.ocmn = tim.OC3Mn
	case 4:
		ccr = &cet.CCR4.U32
		cce = tim.CC4E
		dci.ocmn = tim.OC4Mn
	}
	ccr.Store(200)
	cet.CCER.Store(cce)

	dci.irq = irqline
	irqline.EnableFallTrig()
	irqline.EnableInt()
}

var (
	tts      *usart.Driver
	dmarxbuf [88]byte
	nrfdci   nrfDCI
)

func init() {
	system.Setup96(8)
	systick.Setup()
	start := rtos.Nanosec()

	// GPIO

	gpio.A.EnableClock(true)
	uartport, tx, rx := gpio.A, gpio.Pin2, gpio.Pin3
	spiport, sck, miso, mosi := gpio.A, gpio.Pin5, gpio.Pin6, gpio.Pin7

	gpio.B.EnableClock(true)
	ctrport, csn, irqn, ce := gpio.B, gpio.Pin6, gpio.Pin8, gpio.Pin9
	csnbb := ctrport.OutPins().Bit(6)

	// UART

	uartport.Setup(tx, gpio.Config{Mode: gpio.Alt})
	uartport.Setup(rx, gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	uartport.SetAltFunc(tx|rx, gpio.USART2)
	d := dma.DMA1
	d.EnableClock(true) // DMA clock must remain enabled in sleep mode.
	tts = usart.NewDriver(
		usart.USART2, d.Channel(5, 4), d.Channel(6, 4), dmarxbuf[:],
	)
	tts.P.EnableClock(true)
	tts.P.SetBaudRate(115200)
	tts.P.Enable()
	tts.EnableRx()
	tts.EnableTx()
	rtos.IRQ(irq.USART2).Enable()
	rtos.IRQ(irq.DMA1_Stream5).Enable()
	rtos.IRQ(irq.DMA1_Stream6).Enable()
	fmt.DefaultWriter = linewriter.New(
		bufio.NewWriterSize(tts, 88),
		linewriter.CRLF,
	)

	// nRF24 SPI

	spiport.Setup(sck|mosi, gpio.Config{Mode: gpio.Alt, Speed: gpio.High})
	spiport.Setup(miso, gpio.Config{Mode: gpio.AltIn})
	spiport.SetAltFunc(sck|miso|mosi, gpio.SPI1)
	d = dma.DMA2
	d.EnableClock(true)
	spid := spi.NewDriver(spi.SPI1, d.Channel(2, 3), d.Channel(3, 3))
	spid.P.EnableClock(true)
	rtos.IRQ(irq.SPI1).Enable()
	rtos.IRQ(irq.DMA2_Stream2).Enable()
	rtos.IRQ(irq.DMA2_Stream3).Enable()

	// nRF24 control lines

	ctrport.Setup(csn, gpio.Config{Mode: gpio.Out, Speed: gpio.High})
	ctrport.Setup(ce, gpio.Config{Mode: gpio.Alt, Speed: gpio.High})
	ctrport.SetAltFunc(ce, gpio.TIM4)
	rcc.RCC.TIM4EN().Set()
	ctrport.Setup(irqn, gpio.Config{Mode: gpio.In, Pull: gpio.PullUp})
	irqline := exti.Lines(irqn)
	irqline.Connect(ctrport)
	rtos.IRQ(irq.EXTI9_5).Enable()

	nrfdci.Setup(spid, csnbb, system.APB1.Clock(), tim.TIM4, 4, irqline)

	// nRF24 requires wait at least 100 ms from start before use it.
	rtos.SleepUntil(start + 100e6)
}

func checkErr(err error) {
	if err == nil {
		return
	}
	fmt.Printf("Error: %v\n", err)
	for {
	}
}

func printRegs(nrf *nrf24.Radio) {
	fmt.Printf("CONFIG:      %v\n", nrf.Config())
	fmt.Printf("EN_AA:       %v\n", nrf.AA())
	fmt.Printf("EN_RXADDR:   %v\n", nrf.RxAEn())
	fmt.Printf("SETUP_AW:    %d\n", nrf.AW())
	arc, ard := nrf.Retr()
	fmt.Printf("SETUP_RETR:  %d, %dus\n", arc, ard)
	ch := nrf.Ch()
	fmt.Printf("RF_CH:       %d (%d MHz)\n", ch, 2400+ch)
	fmt.Printf("RF_SETUP:    %v\n", nrf.RF())
	plos, arc := nrf.ObserveTx()
	fmt.Printf("OBSERVE_TX:  %d lost, %d retr\n", plos, arc)
	fmt.Printf("RPD:         %t\n", nrf.RPD())
	var addr [5]byte
	for i := 0; i < 6; i++ {
		n := 5
		if i > 1 {
			n = 1
		}
		nrf.RxAddr(i, addr[:n])
		fmt.Printf("RX_ADDR_P%d:  %x\n", i, addr[:n])
	}
	nrf.TxAddr(addr[:])
	fmt.Printf("Tx_ADDR:     %x\n", addr[:])
	for i := 0; i < 6; i++ {
		fmt.Printf("RX_PW_P%d:    %d\n", i, nrf.RxPW(i))
	}
	fmt.Printf("FIFO_STATUS: %v\n", nrf.FIFO())
	fmt.Printf("DYNPD:       %v\n", nrf.DPL())
	fmt.Printf("FEATURE:     %v\n", nrf.Feature())

	checkErr(nrf.Err)
	fmt.Printf("STATUS:      %v\n", nrf.Status)
}

func main() {
	fmt.Printf(
		"\nPCLK: %.1f MHz, SPI speed: %.1f MHz\n\n",
		float32(nrfdci.spi.P.Bus().Clock())*1e-6,
		float32(nrfdci.spi.P.Baudrate(nrfdci.spi.P.Conf()))*1e-6,
	)
	nrf := nrf24.NewRadio(&nrfdci)
	nrf.SetAA(0)

	printRegs(nrf)

	nrf.SetCE(1)
	delay.Millisec(2000)
	nrf.SetCE(0)
	delay.Millisec(2000)
	nrf.SetCE(2)
}

func exti9_5ISR() {
	lines := exti.Pending() & (exti.L9 | exti.L8 | exti.L7 | exti.L6 | exti.L5)
	lines.ClearPending()
	if lines&nrfdci.irq != 0 {
		dbg.WriteString("nRF24 ISR\n")
	}
}

func ttsISR() {
	tts.ISR()
}

func ttsRxDMAISR() {
	tts.RxDMAISR()
}

func ttsTxDMAISR() {
	tts.TxDMAISR()
}

func nrfSPIISR() {
	nrfdci.spi.ISR()
}

func nrfRxDMAISR() {
	nrfdci.spi.DMAISR(nrfdci.spi.RxDMA)
}

func nrfTxDMAISR() {
	nrfdci.spi.DMAISR(nrfdci.spi.TxDMA)
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.USART2:       ttsISR,
	irq.DMA1_Stream5: ttsRxDMAISR,
	irq.DMA1_Stream6: ttsTxDMAISR,

	irq.EXTI9_5:      exti9_5ISR,
	irq.SPI1:         nrfSPIISR,
	irq.DMA2_Stream2: nrfRxDMAISR,
	irq.DMA2_Stream3: nrfTxDMAISR,
}
