package bcmw

import (
	"delay"
	"encoding/binary/le"
	"rtos"

	"sdcard"
	"sdcard/sdio"
)

func (d *Driver) error() bool {
	return d.sd.Err(false) != nil ||
		d.ioStatus&^sdcard.IO_CURRENT_STATE != 0 ||
		d.timeout
}

func (d *Driver) cmd52(f, addr int, flags sdcard.IORWFlags, val byte) byte {
	if d.error() {
		return 0
	}
	val, d.ioStatus = d.sd.SendCmd(sdcard.CMD52(f, addr, flags, val)).R5()
	return val
}

func (d *Driver) setBackplaneWindow(addr uint32) {
	if d.error() {
		return
	}
	addr &^= 0x7FFF
	win := d.backplaneWindow
	if win == addr {
		return
	}
	d.backplaneWindow = addr
	for n := 0; n < 3; n++ {
		addr >>= 8
		win >>= 8
		if a := byte(addr); a != byte(win) {
			d.cmd52(backplane, sbsdioFunc1SBAddrLow+n, sdcard.Write, a)
		}
	}
}

// The setBackplaneWindow, rbr*, wbr*, wbb* are methods that allow to access
// core registers in the way specific to Sonics Silicon Backplane. More info:
// http://www.gc-linux.org/wiki/Wii:WLAN

func (d *Driver) rbr8(addr uint32) byte {
	d.setBackplaneWindow(addr)
	return d.cmd52(backplane, int(addr&0x7FFF), sdcard.Read, 0)
}

func (d *Driver) wbr8(addr uint32, val byte) {
	d.setBackplaneWindow(addr)
	d.cmd52(backplane, int(addr&0x7FFF), sdcard.Write, val)
}

const busyTimeout = 1e9 // 1 s

func waitDataReady(sd sdcard.Host) bool {
	return sd.Wait(rtos.Nanosec() + busyTimeout)
}

func (d *Driver) rbr32(addr uint32) uint32 {
	d.setBackplaneWindow(addr)
	if d.error() {
		return 0
	}
	sd := d.sd
	if d.timeout = !waitDataReady(sd); d.timeout {
		return 0
	}
	var buf [1]uint64
	sd.SetupData(sdcard.Recv|sdcard.IO|sdcard.Stream, buf[:], 4)
	_, d.ioStatus = sd.SendCmd(sdcard.CMD53(
		backplane, int(addr&0x7FFF|access32bit), sdcard.Read, 4,
	)).R5()
	return le.Decode32(sdcard.AsData(buf[:]).Bytes())
}

func (d *Driver) wbr32(addr uint32, val uint32) {
	d.setBackplaneWindow(addr)
	if d.error() {
		return
	}
	sd := d.sd
	if d.timeout = !waitDataReady(sd); d.timeout {
		return
	}
	var buf [1]uint64
	le.Encode32(sdcard.AsData(buf[:]).Bytes(), val)
	sd.SetupData(sdcard.Send|sdcard.IO|sdcard.Stream, buf[:], 4)
	_, d.ioStatus = sd.SendCmd(sdcard.CMD53(
		backplane, int(addr&0x7FFF|access32bit), sdcard.Write, 4,
	)).R5()
}

func (d *Driver) wbb(addr uint32, buf []uint64) {
	sd := d.sd
	for len(buf) >= 8 {
		d.setBackplaneWindow(addr)
		if d.timeout = !waitDataReady(sd); d.timeout {
			return
		}
		nbl := len(buf) >> 3
		if nbl >= 0x1FF {
			nbl = 0x1FF
		}
		sd.SetupData(sdcard.Send|sdcard.IO|sdcard.Block64, buf, nbl*64)
		_, d.ioStatus = sd.SendCmd(sdcard.CMD53(
			backplane, int(addr&0x7FFF), sdcard.BlockWrite|sdcard.IncAddr, nbl,
		)).R5()
		if d.error() {
			return
		}
		buf = buf[nbl*8:]
		addr += uint32(nbl) * 64
	}
	if len(buf) == 0 {
		return
	}
	d.setBackplaneWindow(addr)
	if d.timeout = !waitDataReady(sd); d.timeout {
		return
	}
	n := len(buf) * 8
	sd.SetupData(sdcard.Send|sdcard.IO|sdcard.Stream, buf, n)
	_, d.ioStatus = sd.SendCmd(sdcard.CMD53(
		backplane, int(addr&0x7FFF|access32bit), sdcard.Write|sdcard.IncAddr, n,
	)).R5()
}

func (d *Driver) enableFunction(f int) {
	if d.error() {
		return
	}
	m := d.cmd52(cia, sdio.CCCR_IOEN, sdcard.Read, 0) | byte(1<<uint(f))
	for retry := 250; retry > 0; retry-- {
		r := d.cmd52(cia, sdio.CCCR_IOEN, sdcard.WriteRead, m)
		if d.error() || r == m {
			return
		}
		delay.Millisec(2)
	}
	d.timeout = true
}

func (d *Driver) disableCore(core int) {
	if d.error() {
		return
	}
	base := d.chip.baseAddr[core]
	d.rbr8(base + ssbResetCtl)
	if d.rbr8(base+ssbResetCtl)&1 != 0 {
		return // Already in reset state.
	}
	d.wbr8(base+ssbIOCtl, 0)
	d.rbr8(base + ssbIOCtl)
	delay.Millisec(1)
	d.wbr8(base+ssbResetCtl, 1)
	delay.Millisec(1)
}

func (d *Driver) resetCore(core int) {
	if d.error() {
		return
	}
	d.disableCore(core)

	// Initialization sequence.

	base := d.chip.baseAddr[core]
	d.wbr8(base+ssbIOCtl, ioCtlClk|ioCtlFGC)
	d.rbr8(base + ssbIOCtl)
	d.wbr8(base+ssbResetCtl, 0)
	delay.Millisec(1)
	d.wbr8(base+ssbIOCtl, ioCtlClk)
	d.rbr8(base + ssbIOCtl)
	delay.Millisec(1)
}
