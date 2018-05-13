package main

import (
	"delay"
	"fmt"

	"sdcard"

	"stm32/hal/gpio"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"

	"stm32/hal/raw/rcc"
	"stm32/hal/raw/sdio"
)

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

	rcc.RCC.SDIOEN().Set()
	sd := sdio.SDIO
	sd.CLKCR.Store(sdio.CLKEN | (48e6/400e3-2)<<sdio.CLKDIVn) // CLK=400 kHz
	sd.POWER.Store(3)                                         // Power on.
}

type Host struct {
	status sdio.STA
}

func (h *Host) Cmd(cmd sdcard.Command, arg uint32) (resp sdcard.Response) {
	if h.status != 0 {
		return
	}
	sd := sdio.SDIO
	sd.ICR.Store(0xFFFFFFFF)
	sd.ARG.Store(sdio.ARG(arg))
	sd.CMD.Store(sdio.CPSMEN | sdio.CMD(cmd)&0xFF)
	for sd.CMDACT().Load() != 0 {
		// Wait for transfer end.
	}
	h.status = sd.STA.Load()
	const errFlags = sdio.CCRCFAIL | sdio.DCRCFAIL | sdio.CTIMEOUT |
		sdio.DTIMEOUT | sdio.TXUNDERR | sdio.RXOVERR
	if cmd&sdcard.RespLen == sdcard.NoResp {
		h.status &= errFlags
		return
	}
	// Wait for response
	for h.status&(sdio.CMDREND|errFlags) == 0 {
		h.status = sd.STA.Load()
	}
	h.status &= errFlags
	if h.status != 0 {
		if h.status&sdio.CCRCFAIL == 0 {
			return
		}
		if r := cmd & sdcard.R; r != sdcard.R3 && r != sdcard.R4 {
			return
		}
		// Ignore CRC error for responses R3 and R4
		h.status &^= sdio.CCRCFAIL
	}
	resp[0] = sd.RESP[0].U32.Load()
	if cmd&sdcard.RespLen == sdcard.LongResp {
		resp[1] = sd.RESP[1].U32.Load()
		resp[2] = sd.RESP[2].U32.Load()
		resp[3] = sd.RESP[3].U32.Load()
	}
	return
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

	h := new(Host)
	ocr := sdcard.V33 | sdcard.HCXC
	v2 := true

	fmt.Printf("\nInitializing SD card")

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

	for i := 0; ocr&sdcard.PWUP == 0 && i < 10; i++ {
		h.Cmd(sdcard.CMD55(0))
		ocr = h.Cmd(sdcard.ACMD41(ocr)).R3()
		checkErr("ACMD41", h.Err(true))
		fmt.Printf(".")
		delay.Millisec(100)
	}
	if ocr&sdcard.PWUP == 0 {
		fmt.Printf(" timeout\n")
		for {
		}
	}
	fmt.Printf(" OK\n\n")
	fmt.Printf("Physicaly layer version 2.00+: %t\n", v2)
	fmt.Printf("Operation Conditions Register: 0x%08X\n\n", ocr)

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

	csd := h.Cmd(sdcard.CMD9(rca)).R2CSD()
	checkErr("CMD9", h.Err(true))

	csdv := csd.Version()
	fmt.Printf("CSD version: %d\n", csdv)
	fmt.Printf("TAAC:        %d ns\n", csd.TAAC())
	fmt.Printf("NSAC:        %d clk\n", csd.NSAC())
	fmt.Printf("TRAN_SPEED:  %d kbit/s\n", csd.TRAN_SPEED())
	fmt.Printf("CCC:         %012b\n", csd.CCC())
}
