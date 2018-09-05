package bcmw

import (
	"delay"
	"encoding/binary/be"
	"rtos"

	"sdcard"
	"sdcard/sdio"
)

func (d *Driver) error() bool {
	return d.sd.Err(false) != nil || d.timeout
}

func (d *Driver) cmd52(f, addr int, flags sdcard.IORWFlags, val byte) byte {
	if d.error() {
		return 0
	}
	val, _ = d.sd.SendCmd(sdcard.CMD52(f, addr, flags, val)).R5()
	return val
}

// The setBackplaneWindow, read32, write32 are methods that allow to access
// core registers in the way specific to Sonics Silicon Backplane. More info:
// http://www.gc-linux.org/wiki/Wii:WLAN

func (d *Driver) setBackplaneWindow(addr uint32) {
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

const busyTimeout = 1e9 // 1 s

func waitDataReady(sd sdcard.Host) bool {
	return sd.Wait(rtos.Nanosec() + busyTimeout)
}

func (d *Driver) read32(f, addr int) uint32 {
	if d.error() {
		return 0
	}
	sd := d.sd
	if d.timeout = !waitDataReady(sd); d.timeout {
		return 0
	}
	var buf [1]uint64
	sd.SetupData(sdcard.Recv, buf[:], 4)
	sd.SendCmd(sdcard.CMD53(f, addr|access32bit, sdcard.Read, 4))
	return be.Decode32(sdcard.AsData(buf[:]).Bytes())
}

func (d *Driver) enableFunction(f int) {
	if d.error() {
		return
	}
	m := byte(1 << uint(f))
	r := d.cmd52(cia, sdio.CCCR_IOEN, sdcard.Read, 0)
	d.cmd52(cia, sdio.CCCR_IOEN, sdcard.Write, r|m)

	for retry := 250; retry > 0; retry-- {
		delay.Millisec(2)
		r = d.cmd52(cia, sdio.CCCR_IORDY, sdcard.Read, 0)
		if d.error() || r&m != 0 {
			return
		}
	}
	d.timeout = true
}

func (d *Driver) disableCore(core int) {
	if d.error() {
		return
	}
	d.setBackplaneWindow(d.chip.baseAddr[core])
	if d.cmd52(backplane, ssbResetCtl, sdcard.Read, 0)&1 != 0 {
		return // Already in reset state.
	}
	delay.Millisec(10)
	d.cmd52(backplane, ssbResetCtl, sdcard.Write, 1)
	delay.Millisec(1)
	d.cmd52(backplane, ssbIOCtl, sdcard.Write, 0)
	d.cmd52(backplane, ssbIOCtl, sdcard.Read, 0)
	delay.Millisec(1)
}

func (d *Driver) resetCore(core int) {
	if d.error() {
		return
	}
	d.disableCore(core)

	// Initialization sequence.

	d.cmd52(backplane, ssbIOCtl, sdcard.Write, ioCtlClk|ioCtlFGC)
	d.read32(backplane, ssbIOCtl)

}
