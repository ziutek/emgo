// Warning! This example destroys all data on your SD card.
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

//var sddrv *sdmmc.Driver
var sddrv *sdmmc.DriverDMA

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
	sddrv = sdmmc.NewDriverDMA(sdmmc.SDIO, d.Channel(6, 4), d0)
	//sddrv = sdmmc.NewDriver(sdmmc.SDIO, d0)
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

	buf := sdcard.MakeDataBlocks(16)
	start := card.Cap() / 8000
	stop := start + start

	fmt.Printf("Write pattern to card: ")

	// Use constant pattern at beggining of buffer to allow to test the transfer
	// with long/weak cables (breadboard). Better (more random) pattern can be
	// used with well-designed connections. The problems with weak connections
	// can be observed when reading multiple blocks (len(buf) > 1) with clock
	// > 2MHz. In case of long (> 5 cm) wires twisting GND and CMD helps much.
	//
	// To stop the ongoing transfer started by CMD18 (READ_MULTIPLE_BLOCK) the
	// CMD12 command must be sent to the card. This is the only case when the
	// signals on DATAx lines can interfere with the command on CMD line. To
	// avoid this the first 3/4 words of the buffer is set to 0.

	pattern := uint64(0)
	for n := range buf {
		if n >= len(buf)*3/4 {
			pattern = uint64(n)
		}
		buf[n] = pattern
	}

	for addr := start; addr < stop; addr += int64(buf.NumBlocks()) {
		err = card.WriteBlocks(addr, buf)
		checkErr(err, card.LastStatus())
		fmt.Printf(".")
	}

	fmt.Printf("\nRead and check pattern: ")

	for addr := start; addr < stop; addr += int64(buf.NumBlocks()) {
		for n := range buf {
			buf[n] = uint64(n)
		}
		err = card.ReadBlocks(addr, buf)
		checkErr(err, card.LastStatus())
		pattern = 0
		for n := range buf {
			if n >= len(buf)*3/4 {
				pattern = uint64(n)
			}
			if buf[n] != pattern {
				fmt.Printf("\nData don't march!\n")
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
