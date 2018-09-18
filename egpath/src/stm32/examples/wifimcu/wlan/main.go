package main

import (
	"bufio"
	"fmt"
	"rtos"
	"text/linewriter"

	"bcmw"

	"stm32/hal/dma"
	"stm32/hal/exti"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/sdmmc"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
	"stm32/hal/usart"
	//"stm32/hal/raw/pwr"
	//"stm32/hal/raw/rcc"
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

	gpio.A.EnableClock(true)
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
	wlan := bcmw.NewDriver(sddrv)

	print("Initialize WLAN:")

	/*
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
	*/

	wlan.Init(bcmRSTn.Store, 0)

	checkErr(wlan.Err(true))
	printOK()

	/*
		print("Uploading firmware:")

		wlan.UploadFirmware(nil, firmware[:])

		checkErr(wlan.Err(true))
		printOK()

		print("Uploading NVRAM:")

		wlan.UploadNVRAM(nil, nvram)

		checkErr(wlan.Err(true))
		printOK()
	*/
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
