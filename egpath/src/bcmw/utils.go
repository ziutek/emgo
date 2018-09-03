package bcmw

import (
	"delay"

	"sdcard"
	"sdcard/sdio"
)

func cmd52(sd sdcard.Host, f, addr int, flags sdcard.IORWFlags, val byte) byte {
	val, _ = sd.SendCmd(sdcard.CMD52(f, addr, flags, val)).R5()
	return val
}

func cmd53read16(sd sdcard.Host, f, addr int) uint16 {
	var buf [1]uint64
	sd.SetupData(sdcard.Recv, buf[:])
	sd.SendCmd(sdcard.CMD53(f, addr, sdcard.Read, 2))
}

// Following code is heavily inspired and sometimes simply translated from WLAN
// code in NuttX (http://nuttx.org/).

func enableFunction(sd sdcard.Host, f int) (timeout bool) {
	m := byte(1 << uint(f))

	r := cmd52(sd, cia, sdio.CCCR_IOEN, sdcard.Read, 0)
	cmd52(sd, cia, sdio.CCCR_IOEN, sdcard.Write, r|m)

	for retry := 250; retry > 0; retry-- {
		delay.Millisec(2)
		r = cmd52(sd, cia, sdio.CCCR_IORDY, sdcard.Read, 0)
		if sd.Err(false) == nil || r&m != 0 {
			return false
		}
	}
	return true
}
