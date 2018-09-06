package main

import (
	"bufio"
	"delay"
	"fmt"
	"rtos"
	"text/linewriter"

	"sdcard"
	"sdcard/sdio"

	"stm32/hal/dma"
	"stm32/hal/exti"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/sdmmc"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
	"stm32/hal/usart"

	"stm32/hal/raw/pwr"
	"stm32/hal/raw/rcc"
)

var (
	led     gpio.Pin
	sddrv   *sdmmc.DriverDMA
	tts     *usart.Driver
	bcmRSTn gpio.Pin
	bcmD1   gpio.Pin
)

func init() {
	system.Setup96(26)
	systick.Setup(2e6)

	// GPIO

	gpio.A.EnableClock(false)
	bcmIRQ := gpio.A.Pin(0)
	tx2 := gpio.A.Pin(2)
	rx2 := gpio.A.Pin(3)
	led = gpio.A.Pin(4)
	bcmCMD := gpio.A.Pin(6)
	//flashMOSI = gpio.A.Pin(7)
	bcmD1 = gpio.A.Pin(8) // Also LSE output (MCO1) to WLAN powersave clock.
	bcmD2 := gpio.A.Pin(9)
	//flashCSn := gpio.A.Pin(15)

	gpio.B.EnableClock(true)
	//flashSCK := gpio.B.Pin(3)
	//flashMISO := gpio.B.Pin(4)
	bcmD3 := gpio.B.Pin(5)
	bcmD0 := gpio.B.Pin(7)
	bcmRSTn = gpio.B.Pin(14)
	bcmCLK := gpio.B.Pin(15)

	// LED

	led.Set()
	led.Setup(&gpio.Config{
		Mode:   gpio.Out,
		Driver: gpio.OpenDrain,
		Speed:  gpio.Low,
	})

	// USART2

	tx2.Setup(&gpio.Config{Mode: gpio.Alt})
	rx2.Setup(&gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	tx2.SetAltFunc(gpio.USART2)
	rx2.SetAltFunc(gpio.USART2)
	d := dma.DMA1
	d.EnableClock(true) // DMA clock must remain enabled in sleep mode.
	tts = usart.NewDriver(
		usart.USART2, d.Channel(6, 4), d.Channel(5, 4), make([]byte, 88),
	)
	tts.Periph().EnableClock(true)
	tts.Periph().SetBaudRate(115200)
	tts.Periph().Enable()
	tts.EnableRx()
	tts.EnableTx()
	rtos.IRQ(irq.USART2).Enable()
	rtos.IRQ(irq.DMA1_Stream5).Enable()
	rtos.IRQ(irq.DMA1_Stream6).Enable()
	fmt.DefaultWriter = linewriter.New(
		bufio.NewWriterSize(tts, 88),
		linewriter.CRLF,
	)

	// WLAN (BCM43362: SDIO, reset, IRQ)

	bcmIRQ.Setup(&gpio.Config{Mode: gpio.In})
	bcmRSTn.Setup(&gpio.Config{Mode: gpio.Out, Speed: gpio.Low})

	cfg := &gpio.Config{Mode: gpio.Alt, Speed: gpio.VeryHigh, Pull: gpio.PullUp}
	for _, pin := range []gpio.Pin{bcmCLK, bcmCMD, bcmD0, bcmD1, bcmD2, bcmD3} {
		pin.Setup(cfg)
		pin.SetAltFunc(gpio.SDIO)
	}
	d = dma.DMA2
	d.EnableClock(true)
	sddrv = sdmmc.NewDriverDMA(sdmmc.SDIO, d.Channel(6, 4), bcmD0)
	sddrv.Periph().EnableClock(true)
	sddrv.Periph().Enable()
	rtos.IRQ(irq.SDIO).Enable()
	rtos.IRQ(irq.EXTI9_5).Enable()
}

func main() {
	sd := sdcard.Host(sddrv)

	// Initialize WLAN

	bcmRSTn.Store(0) // Set WLAN into reset state.

	sd.SetBusWidth(sdcard.Bus1)
	sd.SetClock(400e3, true)

	// Provide WLAN powersave clock on PA8 (SDIO_D1).
	RCC := rcc.RCC
	PWR := pwr.PWR
	RCC.PWREN().Set()
	PWR.DBP().Set()
	RCC.LSEON().Set()
	for RCC.LSERDY().Load() == 0 {
		led.Clear()
		delay.Millisec(50)
		led.Set()
		delay.Millisec(50)
	}
	RCC.MCO1().Store(1 << rcc.MCO1n) // LSE on MCO1.
	PWR.DBP().Clear()
	RCC.PWREN().Clear()
	bcmD1.SetAltFunc(gpio.MCO)

	delay.Millisec(2)
	bcmRSTn.Store(1)

	var (
		retry int
		rca   uint16
		cs    sdcard.CardStatus
	)

	print("\nEnumerate:")
	for retry = 250; retry > 0; retry-- {
		delay.Millisec(2)
		sd.SendCmd(sdcard.CMD0())
		checkErr("CMD0", sd.Err(true), 0)
		sd.SendCmd(sdcard.CMD5(0))
		sd.Err(true)
		rca, cs = sd.SendCmd(sdcard.CMD3()).R6()
		if sd.Err(true) == nil {
			break
		}
	}
	checkRetry(retry)
	fmt.Printf(" RCA=0x%X CardStatus=%s\n", rca, cs)

	print("Select card:")
	cs = sd.SendCmd(sdcard.CMD7(rca)).R1()
	checkErr("CMD7", sd.Err(true), 0)
	printOK()

	print("Enable FN1 (Backplane):")
	for retry = 250; retry > 0; retry-- {
		reg := sendCMD52(sd, CIA, sdio.CCCR_IOEN, sdcard.WriteRead, 1<<SSB)
		if reg&(1<<SSB) != 0 {
			break
		}
		delay.Millisec(2)
	}
	checkRetry(retry)
	printOK()

	print("Enable 4-bit data bus:")
	reg := sendCMD52(sd, CIA, sdio.CCCR_BUSICTRL, sdcard.Read, 0)
	reg = reg&^3 | byte(sdcard.Bus4)
	sendCMD52(sd, CIA, sdio.CCCR_BUSICTRL, sdcard.Write, reg)
	printOK()

	print("Set 64B block size for FN0:")
	for retry = 250; retry > 0; retry-- {
		if sendCMD52(sd, CIA, sdio.CCCR_BLKSIZE0, sdcard.WriteRead, 64) == 64 {
			break
		}
		delay.Millisec(2)
	}
	checkRetry(retry)
	printOK()

	print("Set 64B block size for FN1 (Backplane):")
	sendCMD52(sd, CIA, sdio.FBR1+sdio.FBR_BLKSIZE0, sdcard.Write, 64)
	printOK()

	print("Set 64B block size for FN2 (WLAN data):")
	sendCMD52(sd, CIA, sdio.FBR2+sdio.FBR_BLKSIZE0, sdcard.Write, 64)
	sendCMD52(sd, CIA, sdio.FBR2+sdio.FBR_BLKSIZE1, sdcard.Write, 0)
	printOK()

	/*
		print("Enable client interrupts:")
		sendCMD52(
			sd, CIA, sdio.CCCR_INTEN, sdcard.Write, sdio.IENM|sdio.FN1|sdio.FN2,
		)
		printOK()
	*/

	reg = sendCMD52(sd, CIA, sdio.CCCR_SPEEDSEL, sdcard.Read, 0)
	if reg&1 != 0 {
		print("Enable high speed mode (50 MHz):")
		sendCMD52(sd, CIA, sdio.CCCR_SPEEDSEL, sdcard.Write, reg|2)
		printOK()
		sd.SetClock(50e6, true)
	} else {
		sd.SetClock(25e6, true)
	}

	print("Wait for FN1 (Backplane) is ready:")
	for retry = 250; retry > 0; retry-- {
		if sendCMD52(sd, CIA, sdio.CCCR_IORDY, sdcard.Read, 0)&(1<<SSB) != 0 {
			break
		}
		delay.Millisec(2)
	}
	checkRetry(retry)
	printOK()

	print("Enable ALP clock:")
	sendCMD52(
		sd, BP, CHIP_CLOCK_CSR, sdcard.Write,
		SBSDIO_FORCE_HW_CLKREQ_OFF|SBSDIO_ALP_AVAIL_REQ|SBSDIO_FORCE_ALP,
	)
	for retry = 50; retry > 0; retry-- {
		reg := sendCMD52(sd, BP, CHIP_CLOCK_CSR, sdcard.Read, 0)
		if reg&SBSDIO_ALP_AVAIL != 0 {
			break
		}
		delay.Millisec(2)
	}
	checkRetry(retry)
	sendCMD52(sd, BP, CHIP_CLOCK_CSR, sdcard.Write, 0) // Clear enable request.
	printOK()

	print("Disable BCM43362 SDIO pull-ups:") // We use STM32 GPIO pull-ups.
	sendCMD52(sd, BP, PULL_UP, sdcard.Write, 0)
	printOK()

	print("Enable FN2 (WLAN data):")
	sendCMD52(sd, CIA, sdio.CCCR_IOEN, sdcard.Write, 1<<SSB|1<<WLAN)
	printOK()

	print("Enable out-of-band interrupts:")
	sendCMD52(
		sd, CIA, SEP_INT_CTL, sdcard.Write,
		SEP_INTR_CTL_MASK|SEP_INTR_CTL_EN|SEP_INTR_CTL_POL,
	)
	// EMW3165 uses default IRQ pin (Pin0). Redirection isn't needed.
	printOK()

	print("Enable FN2 (WLAN data) interrupts:")
	sendCMD52(sd, CIA, sdio.CCCR_INTEN, sdcard.Write, 1<<CIA|1<<WLAN)
	printOK()

	print("Ensure FN2 (WLAN data) is ready:")
	for retry = 250; retry > 0; retry-- {
		if sendCMD52(sd, CIA, sdio.CCCR_IORDY, sdcard.Read, 0)&(1<<WLAN) != 0 {
			break
		}
		delay.Millisec(2)
	}
	checkRetry(retry)
	printOK()

	print("Sending firmware:")

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

func sdioISR() {
	sddrv.ISR()
}

func exti9_5ISR() {
	pending := exti.Pending() & 0x3E0
	pending.ClearPending()
	if pending&sddrv.BusyLine() != 0 {
		sddrv.BusyISR()
	}
}

//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.USART2:       ttsISR,
	irq.DMA1_Stream5: ttsRxDMAISR,
	irq.DMA1_Stream6: ttsTxDMAISR,

	irq.SDIO:    sdioISR,
	irq.EXTI9_5: exti9_5ISR,
}
