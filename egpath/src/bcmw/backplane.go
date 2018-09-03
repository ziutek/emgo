package bcmw

import (
	"sdcard"
)

// The methods functions handle read/write operations on Sonics Silicon
// Backplane Core Register Space. More info:
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
			cmd52(d.sd, backplane, sbsdioFunc1SBAddrLow+n, sdcard.Write, a)
		}
	}
}
