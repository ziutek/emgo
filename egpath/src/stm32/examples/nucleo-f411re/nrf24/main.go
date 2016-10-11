package main

import (
	"bufio"
	"fmt"
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
)

const dbg = itm.Port(0)

type nrfDCI struct {
	spi *spi.Driver
	csn bitband.Bit
	irq exti.Lines
}

func (dc *nrfDCI) WriteRead(oi ...[]byte) (n int, err error) {
	dc.csn.Clear()
	dc.spi.WriteReadMany(oi...)
	dc.csn.Set()
	return n, dc.spi.Err()
}

func (dc *nrfDCI) SetCE(v int) error {
	return nil
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
	nrfdci.csn = ctrport.OutPins().Bit(6)

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
	nrfdci.spi = spi.NewDriver(spi.SPI1, d.Channel(2, 3), d.Channel(3, 3))
	nrfdci.spi.P.EnableClock(true)
	nrfdci.spi.P.SetConf(
		spi.Master | spi.MSBF | spi.CPOL0 | spi.CPHA0 |
			nrfdci.spi.P.BR(10e6) | // 10 MHz max.
			spi.SoftSS | spi.ISSHigh,
	)
	nrfdci.spi.P.Enable()
	rtos.IRQ(irq.SPI1).Enable()
	rtos.IRQ(irq.DMA2_Stream2).Enable()
	rtos.IRQ(irq.DMA2_Stream3).Enable()

	// nRF24 control lines

	ctrport.Setup(csn|ce, gpio.Config{Mode: gpio.Out, Speed: gpio.High})
	ctrport.SetPins(csn)
	ctrport.Setup(irqn, gpio.Config{Mode: gpio.In, Pull: gpio.PullUp})
	nrfdci.irq = exti.Lines(irqn)
	nrfdci.irq.Connect(ctrport)
	nrfdci.irq.EnableFallTrig()
	nrfdci.irq.EnableInt()
	rtos.IRQ(irq.EXTI9_5).Enable()

	// nRF24 requires wait at least 100 ms from start before use it.
	rtos.SleepUntil(start + 100e6)
}

func main() {
	fmt.Printf(
		"\nPCLK: %.1f MHz, SPI speed: %.1f MHz\n\n",
		float32(nrfdci.spi.P.Bus().Clock())*1e-6,
		float32(nrfdci.spi.P.Baudrate(nrfdci.spi.P.Conf()))*1e-6,
	)

	nrf := nrf24.Device{DCI: &nrfdci}

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
	if nrf.Err != nil {
		fmt.Printf("Error: %v\n", nrf.Err)
		return
	}
	fmt.Printf("STATUS:      %v\n", nrf.Status)
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
