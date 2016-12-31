package main

import (
	"delay"
	"fmt"
	"rtos"

	"nrf24"

	"stm32/nrfdci"

	"stm32/hal/dma"
	"stm32/hal/exti"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/spi"
	"stm32/hal/system"
	"stm32/hal/system/timer/rtc"

	"stm32/hal/raw/rcc"
	"stm32/hal/raw/tim"
)

var (
	led gpio.Pin
	dci *nrfdci.DCI
)

func init() {
	system.Setup(8, 1, 72/8)
	rtc.Setup(32768)
	start := rtos.Nanosec()

	gpio.B.EnableClock(true)
	ctrport, ce, irqn := gpio.B, gpio.Pin0, gpio.Pin1
	csn := gpio.B.Pin(12)
	spiport, sck, miso, mosi := gpio.B, gpio.Pin13, gpio.Pin14, gpio.Pin15

	gpio.C.EnableClock(false)
	led = gpio.C.Pin(13)

	// LED

	led.Setup(&gpio.Config{Mode: gpio.Out, Driver: gpio.OpenDrain, Speed: gpio.Low})
	led.Set()

	// nRF24 SPI

	spiport.Setup(sck|mosi, &gpio.Config{Mode: gpio.Alt, Speed: gpio.High})
	spiport.Setup(miso, &gpio.Config{Mode: gpio.AltIn})
	csn.Setup(&gpio.Config{Mode: gpio.Out, Speed: gpio.High})
	d := dma.DMA1
	d.EnableClock(true)
	spid := spi.NewDriver(spi.SPI2, d.Channel(4, 0), d.Channel(5, 0))
	spid.P.EnableClock(true)
	rtos.IRQ(irq.SPI2).Enable()
	rtos.IRQ(irq.DMA1_Channel4).Enable()
	rtos.IRQ(irq.DMA1_Channel5).Enable()

	// nRF24 control lines.

	ctrport.Setup(ce, &gpio.Config{Mode: gpio.Alt, Speed: gpio.High})
	rcc.RCC.TIM3EN().Set()
	ctrport.Setup(irqn, &gpio.Config{Mode: gpio.In, Pull: gpio.PullUp})
	irqline := exti.Lines(irqn)
	irqline.Connect(ctrport)
	rtos.IRQ(irq.EXTI1).Enable()

	dci = nrfdci.NewDCI(spid, csn, system.APB1.Clock(), tim.TIM3, 3, irqline)
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
	config, _ := nrf.CONFIG()
	fmt.Printf("CONFIG:      %v\n", config)
	enaa, _ := nrf.EN_AA()
	fmt.Printf("EN_AA:       %v\n", enaa)
	enrxaddr, _ := nrf.EN_RXADDR()
	fmt.Printf("EN_RXADDR:   %v\n", enrxaddr)
	setupaw, _ := nrf.SETUP_AW()
	fmt.Printf("SETUP_AW:    %d\n", setupaw)
	arc, ard, _ := nrf.SETUP_RETR()
	fmt.Printf("SETUP_RETR:  %d, %dus\n", arc, ard)
	rfch, _ := nrf.RF_CH()
	fmt.Printf("RF_CH:       %d (%d MHz)\n", rfch, 2400+rfch)
	rfsetup, _ := nrf.RF_SETUP()
	fmt.Printf("RF_SETUP:    %v\n", rfsetup)
	plos, arc, _ := nrf.OBSERVE_TX()
	fmt.Printf("OBSERVE_TX:  %d lost, %d retr\n", plos, arc)
	rpd, _ := nrf.RPD()
	fmt.Printf("RPD:         %t\n", rpd)
	delay.Millisec(100) // To work with slow ST-LINK SWO receiver.
	var addr [5]byte
	for i := 0; i < 6; i++ {
		n := setupaw
		if i > 1 {
			n = 1
		}
		nrf.Read_RX_ADDR(i, addr[:n])
		fmt.Printf("RX_ADDR_P%d:  %x\n", i, addr[:n])
	}
	nrf.Read_TX_ADDR(addr[:setupaw])
	fmt.Printf("TX_ADDR:     %x\n", addr[:setupaw])
	for i := 0; i < 6; i++ {
		rxpw, _ := nrf.RX_PW(i)
		fmt.Printf("RX_PW_P%d:    %d\n", i, rxpw)
	}
	delay.Millisec(100) // To work with slow ST-LINK SWO receiver.
	fifostatus, _ := nrf.FIFO_STATUS()
	fmt.Printf("FIFO_STATUS: %v\n", fifostatus)
	dynpd, _ := nrf.DYNPD()
	fmt.Printf("DYNPD:       %v\n", dynpd)
	feurure, status := nrf.FEATURE()
	fmt.Printf("FEATURE:     %v\n", feurure)
	checkErr(nrf.Err())
	fmt.Printf("STATUS:      %v\n", status)
}

func main() {
	var buf [32]byte

	spibus := dci.SPI().P.Bus()
	fmt.Printf(
		"\nSPI on %s (%d MHz). SPI speed: %d Hz.\n\n",
		spibus, spibus.Clock()/1e6, dci.Baudrate(),
	)
	nrf := nrf24.NewRadio(dci)

	nrf.Set_RF_SETUP(nrf24.RF_DR_HIGH)
	nrf.Set_EN_AA(0)
	nrf.Set_EN_RXADDR(nrf24.P0)
	nrf.Set_SETUP_AW(3)
	config := nrf24.PWR_UP | nrf24.EN_CRC | nrf24.CRCO | nrf24.PRIM_RX&0
	if config&nrf24.PRIM_RX != 0 {
		nrf.Set_RX_PW(0, len(buf))
	} else {
		nrf.FLUSH_TX()
	}

	nrf.Set_CONFIG(config)
	pwrstart := rtos.Nanosec()

	printRegs(nrf)
	fmt.Println()

	// Wait for transition from Power Down to Standby I.
	rtos.SleepUntil(pwrstart + 4.5e6)

	n := 5000
	for i := 0; ; i++ {
		start := rtos.Nanosec()
		if config&nrf24.PRIM_RX != 0 {
			nrf.FLUSH_RX()
			dci.SetCE(1)
			for i := 0; i < n; i++ {
				// BUG: Must use FIFO_STATUS.
				nrf.ClearIRQ(nrf24.RX_DR)
				dci.IRQF().Reset(0)
				dci.IRQF().Wait(1, 0)
				nrf.R_RX_PAYLOAD(buf[:])
			}
			dci.SetCE(0)
		} else {
			nrf.W_TX_PAYLOAD(buf[:])
			for i := 1; i < n; i++ {
				nrf.ClearIRQ(nrf24.TX_DS)
				dci.IRQF().Reset(0)
				dci.SetCE(2)
				nrf.W_TX_PAYLOAD(buf[:])
				dci.IRQF().Wait(1, 0)
			}
			nrf.ClearIRQ(nrf24.TX_DS)
			dci.IRQF().Reset(0)
			dci.SetCE(2)
			dci.IRQF().Wait(1, 0)
		}
		dt := float32(rtos.Nanosec() - start)
		checkErr(nrf.Err())
		fmt.Printf(
			"%d: %d ms %.0f pkt/s (%.0f kb/s)\n",
			i, uint(dt/1e6),
			float32(n)*1e9/dt,
			float32(n*len(buf)*8)*1e6/dt,
		)
	}
}

func exti1ISR() {
	exti.L1.ClearPending()
	dci.ISR()
}

func nrfSPIISR() {
	dci.SPI().ISR()
}

func nrfRxDMAISR() {
	dci.SPI().DMAISR(dci.SPI().RxDMA)
}

func nrfTxDMAISR() {
	dci.SPI().DMAISR(dci.SPI().TxDMA)
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTCAlarm: rtc.ISR,

	irq.EXTI1:         exti1ISR,
	irq.SPI2:          nrfSPIISR,
	irq.DMA1_Channel4: nrfRxDMAISR,
	irq.DMA1_Channel5: nrfTxDMAISR,
}
