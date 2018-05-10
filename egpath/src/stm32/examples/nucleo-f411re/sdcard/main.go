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
	if cmd&sdcard.RespLen == sdcard.NoResp {
		goto end
	}
	const errFlags = sdio.CCRCFAIL | sdio.DCRCFAIL | sdio.CTIMEOUT |
		sdio.DTIMEOUT | sdio.TXUNDERR | sdio.RXOVERR
	// Wait for response
	for h.status&(sdio.CMDREND|errFlags) == 0 {
		h.status = sd.STA.Load()
	}
	if h.status&errFlags != 0 {
		r := cmd >> 10
		if h.status&sdio.CCRCFAIL == 0 {
			goto end
		}
		if r != sdcard.R3 && r != sdcard.R4 {
			goto end
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
end:
	h.status &= errFlags // Return error flags if any.
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
	err := Error(h.status)
	if clear {
		h.status = 0
	}
	return err
}

func checkErr(h *Host, cmd string) {
	err := h.Err(false)
	if err == nil {
		return
	}
	fmt.Printf("%v: %x\n", err, err)
	for {
	}
}

func main() {
	delay.Millisec(250) // For SWO output

	h := new(Host)

	h.Cmd(sdcard.CMD0())
	checkErr(h, "CMD0")

	vhs, pattern := h.Cmd(sdcard.CMD8(sdcard.V27_36, 0xAC)).R7()
	checkErr(h, "CMD8")

	if vhs != sdcard.V27_36 || pattern != 0xAC {
		fmt.Printf("CMD8 bad resp: %x, %x\n", vhs, pattern)
		return
	}

	status := h.Cmd(sdcard.CMD55(0)).R1()
	checkErr(h, "CMD55")

	fmt.Printf("%b\n", status)
}
