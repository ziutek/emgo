package main

import (
	"delay"
	"encoding/binary/be"
	"fmt"
	"rtos"

	"sdcard"

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/sdmmc"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var sd *sdmmc.DriverDMA

//var sd *sdmmc.Driver

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
	sd = sdmmc.NewDriverDMA(sdmmc.SDIO, d.Channel(6, 4))
	//sd = sdmmc.NewDriver(sdmmc.SDIO)
	sd.Periph().EnableClock(true)
	sd.Periph().Enable()

	rtos.IRQ(irq.SDIO).Enable()
}

func main() {
	delay.Millisec(200) // For SWO output

	ocr := sdcard.V31 | sdcard.V32 | sdcard.V33 | sdcard.HCXC
	v2 := true

	// Set SDIO_CK to no more than 400 kHz (max. open-drain freq). Clock must be
	// continuously enabled (pwrsave = false) to allow correct initialisation.
	sd.SetClock(400e3, false)

	// SD card power-up takes maximum of 1 ms or 74 SDIO_CK cycles.
	delay.Millisec(1)

	// Reset.
	sd.SendCmd(sdcard.CMD0())
	checkErr("CMD0", sd.Err(true), 0)

	// CMD0 may require up to 8 SDIO_CK cycles to reset the card.
	delay.Millisec(1)

	// Verify card interface operating condition.
	vhs, pattern := sd.SendCmd(sdcard.CMD8(sdcard.V27_36, 0xAC)).R7()
	if err := sd.Err(true); err != nil {
		if err == sdcard.ErrCmdTimeout {
			ocr &^= sdcard.HCXC
			v2 = false
		} else {
			checkErr("CMD8", err, 0)
		}
	} else if vhs != sdcard.V27_36 || pattern != 0xAC {
		fmt.Printf("CMD8 bad response: %x, %x\n", vhs, pattern)
		for {
		}
	}
	fmt.Printf("\nPhysical layer version 2.00+: %t\n", v2)

	fmt.Printf("Initializing SD card ")
	var oca sdcard.OCR
	for i := 0; oca&sdcard.PWUP == 0 && i < 20; i++ {
		sd.SendCmd(sdcard.CMD55(0))
		oca = sd.SendCmd(sdcard.ACMD41(ocr)).R3()
		checkErr("ACMD41", sd.Err(true), 0)
		fmt.Printf(".")
		delay.Millisec(50)
	}
	if oca&sdcard.PWUP == 0 {
		fmt.Printf(" timeout\n")
		for {
		}
	}
	fmt.Printf(" OK\n\n")
	fmt.Printf("Operation Conditions Register: 0x%08X\n\n", oca)

	// Read Card Identification Register.
	cid := sd.SendCmd(sdcard.CMD2()).R2CID()
	checkErr("CMD2", sd.Err(true), 0)

	printCID(cid)

	// Generate new Relative Card Address.
	rca, _ := sd.SendCmd(sdcard.CMD3()).R6()
	checkErr("CMD3", sd.Err(true), 0)

	fmt.Printf("Relative Card Address: 0x%04X\n\n", rca)

	// After CMD3 card is in Data Transfer Mode (Standby State) and SDIO_CK can
	// be set to no more than 25 MHz (max. push-pull freq). Clock power save
	// mode can be enabled.
	sd.SetClock(25e6, true)

	// Read Card Specific Data register.
	csd := sd.SendCmd(sdcard.CMD9(rca)).R2CSD()
	checkErr("CMD9", sd.Err(true), 0)

	printCSD(csd)

	// Select card (put into Transfer State).
	st := sd.SendCmd(sdcard.CMD7(rca)).R1()
	checkErr("CMD7", sd.Err(true), st)

	// Now the card is int Transfer State.

	buf := make(sdcard.Data, 8*512/8)

	// Read SD Configuration Register.
	sd.SendCmd(sdcard.CMD55(rca))
	sd.SetupData(sdcard.Recv|sdcard.Block8, buf[:1])
	st = sd.SendCmd(sdcard.ACMD51()).R1()
	checkErr("ACMD51", sd.Err(true), st)

	scr := sdcard.SCR(be.Decode64(buf.Bytes()))
	printSCR(scr)

	// Disable 50k pull-up resistor on D3/CD.
	sd.SendCmd(sdcard.CMD55(rca))
	st = sd.SendCmd(sdcard.ACMD42(false)).R1()
	checkErr("ACMD42", sd.Err(true), st)

	if scr.SD_BUS_WIDTHS()&sdcard.SDBus4 != 0 {
		fmt.Printf("Enable 4-bit data bus... ")
		sd.SendCmd(sdcard.CMD55(rca))
		st = sd.SendCmd(sdcard.ACMD6(sdcard.Bus4)).R1()
		checkErr("ACMD6", sd.Err(true), st)
		sd.SetBusWidth(sdcard.Bus4)
		fmt.Printf("OK\n")
	}

	// Enable High Speed.
	if scr.SD_SPEC() > 0 {
		sd.SetupData(sdcard.Recv|sdcard.Block64, buf[:64/8])
		st = sd.SendCmd(sdcard.CMD6(sdcard.ModeSwitch | sdcard.HighSpeed)).R1()
		checkErr("CMD6", sd.Err(true), st)

		sel := printCMD6Status(buf.Bytes()[:64])

		if sel&sdcard.AccessMode == sdcard.HighSpeed && false {
			fmt.Printf("Card supports High Speed: set clock to 50 MHz.\n\n")
			delay.Millisec(1) // Function switch takes max. 8 SDIO_CK cycles.
			sd.SetClock(50e6, true)
		}
	}

	delay.Millisec(500)

	// Set block size to 512 B (required for protocol version < 2 or SDSC card).
	if ocr&sdcard.HCXC == 0 {
		st = sd.SendCmd(sdcard.CMD16(512)).R1()
		checkErr("CMD16", sd.Err(true), st)
	}

	block := buf[:512/8]
	for i := range block.Bytes() {
		block.Bytes()[i] = byte(i)
	}

	fmt.Printf("Write block of data...\n")
	sd.SetupData(sdcard.Send|sdcard.Block512, block)
	st = sd.SendCmd(sdcard.CMD24(512)).R1()
	checkErr("CMD24", sd.Err(true), st)

	delay.Millisec(500)

	for i := range block {
		block[i] = 0
	}
	fmt.Printf("Read block of data...\n")
	sd.SetupData(sdcard.Recv|sdcard.Block512, block)
	st = sd.SendCmd(sdcard.CMD17(512)).R1()
	checkErr("CMD17", sd.Err(true), st)

	for i, d := 0, block.Bytes(); i < len(d); i += 16 {
		fmt.Printf("%02x\n", d[i:i+16])
	}

	bufSize := len(buf.Bytes())
	testLen := uint(1e4)

	fmt.Printf(
		"Reading %d blocks (%d KiB) using %d B buffer ",
		testLen, testLen/2, bufSize,
	)

	t := rtos.Nanosec()
	for n, step := uint(0), uint(bufSize)/512; n < testLen; n += step {
		fmt.Printf(".")
		addr := n
		if oca&sdcard.HCXC == 0 {
			addr *= 512
		}
		sd.SetupData(sdcard.Recv|sdcard.Block512, buf)
		if len(buf.Bytes()) > 512 {
			st = sd.SendCmd(sdcard.CMD18(addr)).R1()
			checkErr("CMD18", sd.Err(true), st)
			st = sd.SendCmd(sdcard.CMD12()).R1()
			checkErr("CMD12", sd.Err(true), st)
		} else {
			st = sd.SendCmd(sdcard.CMD17(addr)).R1()
			checkErr("CMD17", sd.Err(true), st)
		}
	}
	dt := rtos.Nanosec() - t
	fmt.Printf("%d KiB/s\n", 1e9/2*int64(testLen)/dt)
}

func sdioISR() {
	sd.ISR()
}

//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.SDIO: sdioISR,
}
