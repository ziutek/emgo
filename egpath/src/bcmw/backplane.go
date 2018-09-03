package main

// These functions handle read/write operations on Sonics Silicon Backplane
// Core Register Space. More info: http://www.gc-linux.org/wiki/Wii:WLAN

func (d *Driver) setBackplaneWindow(addr uint32) {
	win := d.backplaneWindow
	addr = addr >> 8 & 0x7F
	d.backplaneWindow = addr
	for n := 0; n < 3; n++ {
		if a := addr & 0xFF; a != win&0xFF {
			cmd52(d.sd, backplane, sbsdioFunc1SBAddrLow+n, sdcard.Write, a)
		}
		addr >>= 8
		win >>= 8
	}
}
