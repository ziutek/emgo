package main

import (
	"delay"
	"fmt"
	"rtos"

	"sdcard"
	"sdcard/sdmc"

	"stm32/hal/dma"
	"stm32/hal/exti"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/sdmmc"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var sddrv *sdmmc.Driver

func init() {
	system.Setup96(8) // Setups USB/SDIO/RNG clock to 48 MHz
	systick.Setup(2e6)

	// GPIO

	gpio.A.EnableClock(true)
	// irq := gpio.A.Pin(0)
	cmd := gpio.A.Pin(6)
	d1 := gpio.A.Pin(8)
	d2 := gpio.A.Pin(9)

	gpio.B.EnableClock(true)
	d3 := gpio.B.Pin(5)
	d0 := gpio.B.Pin(7)
	clk := gpio.B.Pin(15)

	cfg := gpio.Config{Mode: gpio.Alt, Speed: gpio.VeryHigh, Pull: gpio.PullUp}
	for _, pin := range []gpio.Pin{clk, cmd, d0, d1, d2, d3} {
		pin.Setup(&cfg)
		pin.SetAltFunc(gpio.SDIO)
	}

	d := dma.DMA2
	d.EnableClock(true)
	//sddrv = sdmmc.NewDriverDMA(sdmmc.SDIO, d.Channel(6, 4))
	sddrv = sdmmc.NewDriver(sdmmc.SDIO, d0)
	sddrv.Periph().EnableClock(true)
	sddrv.Periph().Enable()

	rtos.IRQ(irq.SDIO).Enable()
	rtos.IRQ(irq.EXTI9_5).Enable()
}

func main() {
	delay.Millisec(200) // For SWO output

	card := sdmc.NewCard(sddrv)

	fmt.Printf("\nInitializing SD Memory Card...\n")
	cid, err := card.Init(
		25e6, sdcard.Bus4,
		sdcard.HCXC|sdcard.V30|sdcard.V31|sdcard.V32|sdcard.V33,
	)
	// Init can return valid CID and capacity even if it returned error.
	printCID(cid)
	fmt.Printf(
		"Card capacity:         %d blk ≈ %d MiB ≈ %d MB\n",
		card.Cap(), card.Cap()>>11, card.Cap()*512/1e6,
	)
	checkErr(err, card.LastStatus())

	buf := sdcard.MakeDataBlocks(32)
	start := card.Cap() / 1000
	stop := start + start

	fmt.Printf("Write pattern to card: ")

	for addr := start; addr < stop; addr += int64(buf.NumBlocks()) {
		for n := range buf.Words() {
			buf.Words()[n] = uint64(addr + int64(n))
		}
		err = card.WriteBlocks(addr, buf)
		checkErr(err, card.LastStatus())
		fmt.Printf(".")
	}

	fmt.Printf("\nRead and check pattern: ")

	for addr := start; addr < stop; addr += int64(buf.NumBlocks()) {
		for n := range buf.Words() {
			buf.Words()[n] = 0
		}
		err = card.ReadBlocks(addr, buf)
		checkErr(err, card.LastStatus())
		for n := range buf.Words() {
			if buf.Words()[n] != uint64(addr+int64(n)) {
				fmt.Printf("Data don't march!")
				return
			}
		}
		fmt.Printf(".")
	}
	fmt.Printf("\nEnd\n")
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
	irq.SDIO:    sdioISR,
	irq.EXTI9_5: exti9_5ISR,
}
