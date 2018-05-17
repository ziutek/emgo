package main

import (
	"delay"
	"fmt"
	"rtos"
	"unsafe"

	"sdcard"

	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"

	"stm32/hal/raw/rcc"
	"stm32/hal/raw/sdio"
)

type Host struct {
	dma    *dma.Channel
	status sdio.STA
}

var h Host

func init() {
	system.Setup96(8) // Setups USB/SDIO/RNG clock to 48 MHz
	systick.Setup(2e6)

	// GPIO

	gpio.A.EnableClock(true)
	cmd := gpio.A.Pin(6)
	//d1 := gpio.A.Pin(8)
	//d2 := gpio.A.Pin(9)

	gpio.B.EnableClock(true)
	//d3 := gpio.B.Pin(5)
	d0 := gpio.B.Pin(7)
	clk := gpio.B.Pin(15)

	cfg := gpio.Config{Mode: gpio.Alt, Speed: gpio.High}
	clk.Setup(&cfg)
	clk.SetAltFunc(gpio.SDIO)
	cfg.Pull = gpio.PullUp
	for _, pin := range []gpio.Pin{cmd, d0 /*d1, d2, d3*/} {
		pin.Setup(&cfg)
		pin.SetAltFunc(gpio.SDIO)
	}

	d := dma.DMA2
	d.EnableClock(true)
	h.dma = d.Channel(6, 4)

	rcc.RCC.SDIOEN().Set()
	sd := sdio.SDIO
	sd.CLKCR.Store(sdio.CLKEN | (48e6/400e3-2)<<sdio.CLKDIVn) // CLK=400 kHz
	sd.POWER.Store(3)                                         // Power on.
}

// SDIO Errata Sheet DocID027036 Rev 2 workarounds:
// 2.7.1 Don't use HW flow control (CLKCR.HWFC_EN).
// 2.7.2 Ignore STA.CCRCFAIL for R3 and R4.
// 2.7.3 Don't use clock dephasing (CLKCR.NEGEDGE).
// 2.7.5 Ensure 3*period(PCLK2)+3*period(SDIOCLK) < 32/BusWidth*period(SDIO_CK)
//       (always met for PCLK2 (APB2CLK) > 28.8 MHz).

func (h *Host) Cmd(cmd sdcard.Command, arg uint32) (resp sdcard.Response) {
	if h.status != 0 {
		return
	}
	sd := sdio.SDIO
	sd.ICR.Store(0xFFFFFFFF)
	sd.ARG.Store(sdio.ARG(arg))
	sd.CMD.Store(sdio.CPSMEN | sdio.CMD(cmd)&0xFF)
	errFlags := sdio.CCRCFAIL | sdio.DCRCFAIL | sdio.CTIMEOUT | sdio.DTIMEOUT |
		sdio.TXUNDERR | sdio.RXOVERR
	waitFlags := errFlags
	if cmd&sdcard.HasResp == 0 {
		waitFlags |= sdio.CMDSENT
	} else {
		waitFlags |= sdio.CMDREND
	}
	for {
		h.status = sd.STA.Load()
		if h.status&waitFlags != 0 {
			break
		}
		rtos.SchedYield()
	}
	h.status &= errFlags
	if cmd&sdcard.HasResp == 0 {
		return
	}
	if h.status != 0 {
		if h.status&sdio.CCRCFAIL == 0 {
			return
		}
		if r := cmd & sdcard.RespType; r != sdcard.R3 && r != sdcard.R4 {
			return
		}
		// Ignore CRC error for R3, R4 responses.
		h.status &^= sdio.CCRCFAIL
	}
	if cmd&sdcard.LongResp != 0 {
		resp[3] = sd.RESP[0].U32.Load() // most significant bits
		resp[2] = sd.RESP[1].U32.Load()
		resp[1] = sd.RESP[2].U32.Load()
		resp[0] = sd.RESP[3].U32.Load() // least significant bits
	} else {
		resp[0] = sd.RESP[0].U32.Load()
	}
	return
}

func (h *Host) Recv(buf []uint32) {
	if h.status != 0 || len(buf) == 0 {
		return
	}
	sd := sdio.SDIO
	ch := h.dma
	ch.Setup(dma.PTM | dma.PFC | dma.IncM | dma.FIFO_1_4)
	ch.SetWordSize(4, 4)
	ch.SetBurst(4, 1)
	ch.SetAddrP(unsafe.Pointer(sd.FIFO.Addr()))
	ch.SetAddrM(unsafe.Pointer(&buf[0]))
	ch.Enable()
	// ...
}

type Error sdio.STA

func (err Error) Error() string {
	return "SDIO error"
}

func (h *Host) Err(clear bool) error {
	if h.status == 0 {
		return nil
	}
	var err error
	if h.status == sdio.CTIMEOUT {
		err = sdcard.ErrTimeout
	} else {
		err = Error(h.status)
	}
	if clear {
		h.status = 0
	}
	return err
}

func checkErr(what string, err error) {
	if err == nil {
		return
	}
	fmt.Printf("%s: %v (0x%X)\n", what, err, err)
	for {
	}
}

func main() {
	delay.Millisec(250) // For SWO output

	ocr := sdcard.V33 | sdcard.HCXC
	v2 := true

	fmt.Printf("\nInitializing SD card ")

	h.Cmd(sdcard.CMD0())
	checkErr("CMD0", h.Err(true))

	vhs, pattern := h.Cmd(sdcard.CMD8(sdcard.V27_36, 0xAC)).R7()
	if err := h.Err(true); err != nil {
		if err == sdcard.ErrTimeout {
			ocr &^= sdcard.HCXC
			v2 = false
		} else {
			checkErr("CMD8", err)
		}
	}
	if vhs != sdcard.V27_36 || pattern != 0xAC {
		fmt.Printf("CMD8 bad response: %x, %x\n", vhs, pattern)
		for {
		}
	}
	var oca sdcard.OCR
	for i := 0; oca&sdcard.PWUP == 0 && i < 20; i++ {
		h.Cmd(sdcard.CMD55(0))
		oca = h.Cmd(sdcard.ACMD41(ocr)).R3()
		checkErr("ACMD41", h.Err(true))
		fmt.Printf(".")
		delay.Millisec(50)
	}
	if oca&sdcard.PWUP == 0 {
		fmt.Printf(" timeout\n")
		for {
		}
	}
	fmt.Printf(" OK\n\n")
	fmt.Printf("Physicaly layer version 2.00+: %t\n", v2)
	fmt.Printf("Operation Conditions Register: 0x%08X\n\n", oca)

	cid := h.Cmd(sdcard.CMD2()).R2CID()
	checkErr("CMD2", h.Err(true))

	y, m := cid.MDT()
	pnm := cid.PNM()
	oid := cid.OID()
	prv := cid.PRV()
	fmt.Printf("Manufacturer ID:       %d\n", cid.MID())
	fmt.Printf("OEM/Application ID:    %s\n", oid[:])
	fmt.Printf("Product name:          %s\n", pnm[:])
	fmt.Printf("Product revision:      %d.%d\n", prv>>4&15, prv&15)
	fmt.Printf("Product serial number: %d\n", cid.PSN())
	fmt.Printf("Manufacturing date:    %04d-%02d\n\n", y, m)

	rca, _ := h.Cmd(sdcard.CMD3()).R6()
	checkErr("CMD3", h.Err(true))

	fmt.Printf("Relative Card Address: 0x%04X\n\n", rca)

	// After CMD3 card is in Data Transfer Mode and CLK can be up to 25 MHz.

	csd := h.Cmd(sdcard.CMD9(rca)).R2CSD()
	checkErr("CMD9", h.Err(true))

	csdv := csd.Version()
	fmt.Printf("CSD version:        %d\n", csdv)
	fmt.Printf("TAAC:               %d ns\n", csd.TAAC())
	fmt.Printf("NSAC:               %d clk\n", csd.NSAC())
	fmt.Printf("TRAN_SPEED:         %d kbit/s\n", csd.TRAN_SPEED())
	fmt.Printf("CCC:                0b%012b\n", csd.CCC())
	fmt.Printf("READ_BL_LEN:        %d B\n", csd.READ_BL_LEN())
	fmt.Printf("READ_BL_PARTIAL:    %t\n", csd.READ_BL_PARTIAL())
	fmt.Printf("WRITE_BLK_MISALIGN: %t\n", csd.WRITE_BLK_MISALIGN())
	fmt.Printf("READ_BLK_MISALIGN:  %t\n", csd.READ_BLK_MISALIGN())
	fmt.Printf("DSR_IMP:            %t\n", csd.DSR_IMP())
	csize := csd.C_SIZE()
	fmt.Printf("C_SIZE:             %d KiB (%d kB)\n", csize>>1, csize<<9/1000)
	fmt.Printf("ERASE_BLK_EN:       %t\n", csd.ERASE_BLK_EN())
	fmt.Printf("SECTOR_SIZE:        %d * WRITE_BL_LEN\n", csd.SECTOR_SIZE())
	fmt.Printf("WP_GRP_SIZE:        %d * SECTOR_SIZE\n", csd.WP_GRP_SIZE())
	fmt.Printf("WP_GRP_ENABLE:      %t\n", csd.WP_GRP_ENABLE())
	fmt.Printf("R2W_FACTOR:         %d\n", csd.R2W_FACTOR())
	fmt.Printf("WRITE_BL_LEN:       %d B\n", csd.WRITE_BL_LEN())
	fmt.Printf("WRITE_BL_PARTIAL:   %t\n", csd.WRITE_BL_PARTIAL())
	fmt.Printf("FILE_FORMAT:        %d\n", csd.FILE_FORMAT())
	fmt.Printf("COPY:               %t\n", csd.COPY())
	fmt.Printf("PERM_WRITE_PROTECT: %t\n", csd.PERM_WRITE_PROTECT())
	fmt.Printf("TMP_WRITE_PROTECT:  %t\n", csd.TMP_WRITE_PROTECT())
}
