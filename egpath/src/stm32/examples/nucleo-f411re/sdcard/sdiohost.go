package main

import (
	"unsafe"

	"sdcard"

	"stm32/hal/dma"
	"stm32/hal/sdmmc"
)

// SDIO Errata Sheet DocID027036 Rev 2 workarounds:
// 2.7.1 Don't use HW flow control (CLKCR.HWFC_EN).
// 2.7.2 Ignore STA.CCRCFAIL for R3 and R4.
// 2.7.3 Don't use clock dephasing (CLKCR.NEGEDGE).
// 2.7.5 Ensure 3*period(PCLK2)+3*period(SDIOCLK) < 32/BusWidth*period(SDIO_CK)
//       (always met for PCLK2 (APB2CLK) > 28.8 MHz).

type Host struct {
	p   *sdmmc.Periph
	dma *dma.Channel
	err sdmmc.Error
}

// SetFreq sets SDIO divider to provide SDIO_CK frequency <= freqHz.
func (h *Host) SetFreq(freqHz int, pwrsave bool) {
	clkdiv := (48e6 + freqHz - 1) / freqHz
	cfg := sdmmc.ClkEna
	if pwrsave {
		cfg |= sdmmc.PwrSave
	}
	// BUG: handle clkdiv == 1
	h.p.SetBusClock(cfg, clkdiv-2)
}

func (h *Host) Cmd(cmd sdcard.Command, arg uint32) (resp sdcard.Response) {
	if h.err != 0 {
		return
	}
	h.p.Clear(sdmmc.EvAll, sdmmc.ErrAll)
	h.p.SetArg(arg)
	h.p.SetCmd(sdmmc.CmdEna | sdmmc.Command(cmd)&255)

	var waitFor sdmmc.Event
	if cmd&sdcard.HasResp != 0 {
		waitFor = sdmmc.CmdRespOK
	} else {
		waitFor = sdmmc.CmdSent
	}
	for {
		ev, err := h.p.Status()
		if err != 0 || ev&waitFor != 0 {
			h.err = err
			break
		}
		//rtos.SchedYield()
	}
	if cmd&sdcard.HasResp == 0 {
		return
	}
	if h.err != 0 {
		if h.err&sdmmc.ErrCmdCRC == 0 {
			return
		}
		if r := cmd & sdcard.RespType; r != sdcard.R3 && r != sdcard.R4 {
			return
		}
		// Ignore CRC error for R3, R4 responses.
		h.err &^= sdmmc.ErrCmdCRC
	}
	if cmd&sdcard.LongResp != 0 {
		resp[3] = h.p.Resp(0) // Most significant bits.
		resp[2] = h.p.Resp(1)
		resp[1] = h.p.Resp(2)
		resp[0] = h.p.Resp(3) // Least significant bits.
	} else {
		resp[0] = h.p.Resp(0)
	}
	return
}

func (h *Host) SetupDMA(dir dma.Mode, buf []uint32) {
	if h.err != 0 || len(buf) == 0 {
		return
	}
	h.p.SetDataLen(len(buf) * 4)
	h.p.SetDataTimeout(0xFFFFFFFF)
	ch := h.dma
	ch.Clear(dma.EvAll, dma.ErrAll)
	ch.Setup(dir | dma.PFC | dma.IncM | dma.FT4 | dma.PB4 | dma.MB4)
	ch.SetWordSize(4, 4)
	ch.SetAddrP(unsafe.Pointer(h.p.FIFOAddr()))
	ch.SetAddrM(unsafe.Pointer(&buf[0]))
	ch.Enable()
}

func (h *Host) StartBlockTransfer(dir sdmmc.DataCtrl) {
	h.p.SetDataCtrl(sdmmc.DTEna|sdmmc.UseDMA|dir, 9)
}

func (h *Host) Err(clear bool) error {
	if h.err == 0 {
		return nil
	}
	var err error
	if h.err == sdmmc.ErrCmdTimeout {
		err = sdcard.ErrCmdTimeout
	} else {
		err = h.err
	}
	if clear {
		h.err = 0
	}
	return err
}
