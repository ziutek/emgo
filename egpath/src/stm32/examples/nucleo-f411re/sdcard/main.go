package main

import (
	"delay"
	"fmt"

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

const (
	shortResp = 1 << sdio.WAITRESPn // Total 48 bits (start,content,stop)
	longResp  = 3 << sdio.WAITRESPn // Total 136 bits (start,content,CRC,stop)
)

const sdioErrFlags = sdio.CCRCFAIL | sdio.DCRCFAIL | sdio.CTIMEOUT | sdio.DTIMEOUT |
	sdio.TXUNDERR | sdio.RXOVERR

func sdioCMD(cmd sdio.CMD, arg sdio.ARG) (status sdio.STA) {
	sd := sdio.SDIO
	for i := 0; i < 10; i++ {
		sd.ICR.Store(0xFFFFFFFF)
		sd.ARG.Store(arg)
		sd.CMD.Store(sdio.CPSMEN | cmd)
		for sd.CMDACT().Load() != 0 {
			// Wait for transfer end.
		}
		if cmd&sdio.WAITRESP == 0 {
			break
		}
		// Wait for response
		for {
			status = sd.STA.Load()
			if status&(sdio.CMDREND|sdioErrFlags) != 0 {
				break
			}
		}
		if status&(sdio.CMDREND|sdio.CTIMEOUT) != 0 {
			break // Response received or timeout
		}
		if cid := cmd & sdio.CMDINDEX; status&sdio.CCRCFAIL != 0 &&
			(cid == 5 || cid == 41) {
			// SDIO periph always checks CRC.
			// Ignore CRC error for commands 5 and 41.
			status &^= sdio.CCRCFAIL
			break
		}
		// Try again.
	}
	return status & sdioErrFlags // Return error flags if any.
}

const (
	resp7 = shortResp
)

const (
	GO_IDLE_STATE = 0
	SEND_IF_COND  = 8 | resp7
)

func main() {
	delay.Millisec(200) // For SWO output

	sd := sdio.SDIO

	status := sdioCMD(GO_IDLE_STATE, 0)
	checkStatus("GO_IDLE_STATE", status)
	status = sdioCMD(SEND_IF_COND, 0x1AA)
	checkStatus("SEND_IF_COND", status)
	fmt.Printf("SEND_IF_COND resp: %x\n", sd.RESP[0].Load())
}

func checkStatus(cmd string, status sdio.STA) {
	if status == 0 {
		return
	}
	fmt.Printf("%s error: %x\n", cmd, status)
	for {
	}
}
