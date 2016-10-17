package main

import (
	"delay"
	"fmt"
	"rtos"

	"arch/cortexm/bitband"

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
	led bitband.Bit
	dci *nrfdci.DCI
)

func init() {
	system.Setup(8, 72/8, false)
	rtc.Setup(32768)
	start := rtos.Nanosec()

	gpio.A.EnableClock(true)
	spiport, csn, sck, miso, mosi := gpio.A, gpio.Pin4, gpio.Pin5, gpio.Pin6, gpio.Pin7

	gpio.B.EnableClock(true)
	ctrport, ce, irqn := gpio.B, gpio.Pin0, gpio.Pin1

	gpio.C.EnableClock(false)
	ledport, ledpin := gpio.C, 13

	// LED

	ledport.SetupPin(
		ledpin,
		gpio.Config{Mode: gpio.Out, Driver: gpio.OpenDrain, Speed: gpio.Low},
	)
	led = ledport.OutPins().Bit(ledpin)
	led.Set()

	// nRF24 SPI

	spiport.Setup(sck|mosi, gpio.Config{Mode: gpio.Alt, Speed: gpio.High})
	spiport.Setup(miso, gpio.Config{Mode: gpio.AltIn})
	spiport.Setup(csn, gpio.Config{Mode: gpio.Out, Speed: gpio.High})
	d := dma.DMA1
	d.EnableClock(true)
	spid := spi.NewDriver(spi.SPI1, d.Channel(2, 0), d.Channel(3, 0))
	spid.P.EnableClock(true)
	rtos.IRQ(irq.SPI1).Enable()
	rtos.IRQ(irq.DMA1_Channel2).Enable()
	rtos.IRQ(irq.DMA1_Channel3).Enable()

	// nRF24 control lines.

	ctrport.Setup(ce, gpio.Config{Mode: gpio.Alt, Speed: gpio.High})
	rcc.RCC.TIM3EN().Set()
	ctrport.Setup(irqn, gpio.Config{Mode: gpio.In, Pull: gpio.PullUp})
	irqline := exti.Lines(irqn)
	irqline.Connect(ctrport)
	rtos.IRQ(irq.EXTI1).Enable()

	dci = nrfdci.NewDCI(
		spid, spiport, csn, system.APB1.Clock(), tim.TIM3, 3, irqline,
	)
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
	cfg, _ := nrf.Config()
	fmt.Printf("CONFIG:      %v\n", cfg)
	aa, _ := nrf.AA()
	fmt.Printf("EN_AA:       %v\n", aa)
	rxae, _ := nrf.RxAEn()
	fmt.Printf("EN_RXADDR:   %v\n", rxae)
	aw, _ := nrf.AW()
	fmt.Printf("SETUP_AW:    %d\n", aw)
	arc, ard, _ := nrf.Retr()
	fmt.Printf("SETUP_RETR:  %d, %dus\n", arc, ard)
	ch, _ := nrf.Ch()
	fmt.Printf("RF_CH:       %d (%d MHz)\n", ch, 2400+ch)
	rf, _ := nrf.RF()
	fmt.Printf("RF_SETUP:    %v\n", rf)
	plos, arc, _ := nrf.ObserveTx()
	fmt.Printf("OBSERVE_TX:  %d lost, %d retr\n", plos, arc)
	rpd, _ := nrf.RPD()
	fmt.Printf("RPD:         %t\n", rpd)
	var addr [5]byte
	for i := 0; i < 6; i++ {
		fmt.Printf("RX_ADDR_P%d:  ", i)
		if i < 2 {
			nrf.RxAddr(i, addr[:])
			fmt.Printf("%x\n", addr[:])
		} else {
			lsb, _ := nrf.RxAddrLSB(i)
			fmt.Printf("%x\n", lsb)
		}
	}
	delay.Millisec(100) // To work with slow ST-LINK SWO receiver.
	nrf.TxAddr(addr[:])
	fmt.Printf("Tx_ADDR:     %x\n", addr[:])
	for i := 0; i < 6; i++ {
		rxpw, _ := nrf.RxPW(i)
		fmt.Printf("RX_PW_P%d:    %d\n", i, rxpw)
	}
	fifo, _ := nrf.FIFO()
	fmt.Printf("FIFO_STATUS: %v\n", fifo)
	dpl, _ := nrf.DPL()
	fmt.Printf("DYNPD:       %v\n", dpl)
	feurure, status := nrf.Feature()
	fmt.Printf("FEATURE:     %v\n", feurure)
	checkErr(nrf.Err())
	fmt.Printf("STATUS:      %v\n", status)
}

func main() {
	var buf [32]byte

	fmt.Printf("\nSPI speed: %d Hz\n\n", dci.Baudrate())
	nrf := nrf24.NewRadio(dci)
	//nrf.SetRF(nrf24.DRHigh)
	nrf.SetRF(0)
	//nrf.SetRF(nrf24.DRLow)
	nrf.SetAA(0)
	mode := nrf24.PrimRx
	if mode == nrf24.PrimRx {
		nrf.SetRxPW(0, len(buf))
	}
	nrf.SetConfig(nrf24.PwrUp | mode)
	pwrstart := rtos.Nanosec()

	printRegs(nrf)
	fmt.Println()

	rtos.SleepUntil(pwrstart + 4.5e6) // Wait for PowerDown -> StandbyI.

	n := 5000
	for i := 0; ; i++ {
		start := rtos.Nanosec()
		if mode == nrf24.PrimRx {
			nrf.FlushRx()
			dci.SetCE(1)
			for i := 0; i < n; i++ {
				nrf.Clear(nrf24.RxDR)
				dci.Wait(0)
				led.Clear()
				nrf.ReadRx(buf[:])
				led.Set()
			}
			dci.SetCE(0)
		} else {
			for i := 0; i < n; i++ {
				nrf.WriteTx(buf[:])
				nrf.Clear(nrf24.TxDS)
				dci.SetCE(2)
				dci.Wait(0)
			}
		}
		dt := float32(rtos.Nanosec() - start)
		checkErr(nrf.Err())
		fmt.Printf(
			"%d: %.0f pkt/s (%.0f kb/s)\n",
			i,
			float32(n)*1e9/dt,
			float32(n*len(buf)*8)*1e6/dt,
		)
	}
	/*
		for {
			delay.Millisec(100)
			led.Clear()
			delay.Millisec(1000)
			led.Set()
			fmt.Printf("Hello!\n")
		}
	*/
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
	irq.SPI1:          nrfSPIISR,
	irq.DMA1_Channel2: nrfRxDMAISR,
	irq.DMA1_Channel3: nrfTxDMAISR,
}
