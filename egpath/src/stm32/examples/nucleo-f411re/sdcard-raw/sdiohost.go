package main

import (
	"rtos"
	"unsafe"

	"sdcard"

	"stm32/hal/dma"

	"stm32/hal/raw/sdio"
)

// SDIO Errata Sheet DocID027036 Rev 2 workarounds:
// 2.7.1 Don't use HW flow control (CLKCR.HWFC_EN).
// 2.7.2 Ignore STA.CCRCFAIL for R3 and R4.
// 2.7.3 Don't use clock dephasing (CLKCR.NEGEDGE).
// 2.7.5 Ensure 3*period(PCLK2)+3*period(SDIOCLK) < 32/BusWidth*period(SDIO_CK)
//       (always met for PCLK2 (APB2CLK) > 28.8 MHz).

type Host struct {
	dma    *dma.Channel
	status sdio.STA
}

func (h *Host) Enable() {
	sdio.SDIO.POWER.Store(3)
}

func (h *Host) Disable() {
	sdio.SDIO.POWER.Store(0)
}

// SetFreq sets SDIO divider to provide SDIO_CK frequency <= freqHz.
func (h *Host) SetFreq(freqHz int, pwrsave bool) {
	div := sdio.CLKCR((48e6 + freqHz - 1) / freqHz)
	// BUG: handle clkdiv == 1
	clkcr := sdio.CLKEN | (div-2)<<sdio.CLKDIVn
	if pwrsave {
		clkcr |= sdio.PWRSAV
	}
	sdio.SDIO.CLKCR.Store(clkcr)
}

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
		resp[3] = sd.RESP[0].U32.Load() // Most significant bits.
		resp[2] = sd.RESP[1].U32.Load()
		resp[1] = sd.RESP[2].U32.Load()
		resp[0] = sd.RESP[3].U32.Load() // Least significant bits.
	} else {
		resp[0] = sd.RESP[0].U32.Load()
	}
	return
}

func (h *Host) SetupDMA(dir dma.Mode, buf []uint32) {
	if h.status != 0 || len(buf) == 0 {
		return
	}
	sd := sdio.SDIO
	sd.DLEN.Store(sdio.DLEN(len(buf)) * 4)
	sd.DTIMER.Store(0xFFFFFFFF)
	ch := h.dma
	ch.Clear(dma.EvAll, dma.ErrAll)
	ch.Setup(dir | dma.PFC | dma.IncM | dma.FT4 | dma.PB4 | dma.MB4)
	ch.SetWordSize(4, 4)
	ch.SetAddrP(unsafe.Pointer(sd.FIFO.Addr()))
	ch.SetAddrM(unsafe.Pointer(&buf[0]))
	ch.Enable()
}

// Set DTEN, DBLOCKSIZE, Wait for SDIO DBCKEND.

const (
	WR = 0 << sdio.DTDIRn
	RD = 1 << sdio.DTDIRn
)

func (h *Host) StartBlockTransfer(dir sdio.DCTRL) {
	sdio.SDIO.DCTRL.Store(sdio.DTEN | dir | sdio.DMAEN | 9<<sdio.DBLOCKSIZEn)
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
		err = sdcard.ErrCmdTimeout
	} else {
		err = Error(h.status)
	}
	if clear {
		h.status = 0
	}
	return err
}
