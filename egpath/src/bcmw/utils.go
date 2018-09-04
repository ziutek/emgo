package bcmw

import (
	"delay"
	"encoding/binary/be"

	"sdcard"
	"sdcard/sdio"
)

// The following code is heavily inspired and sometimes simply translated from
// SDIO and WLAN code in NuttX (http://nuttx.org/).

func cmd52(sd sdcard.Host, f, addr int, flags sdcard.IORWFlags, val byte) byte {
	val, _ = sd.SendCmd(sdcard.CMD52(f, addr, flags, val)).R5()
	return val
}

func cmd53read32(sd sdcard.Host, f, addr int) uint32 {
	var buf [1]uint64
	sd.SetupData(sdcard.Recv, buf[:], 4)
	sd.SendCmd(sdcard.CMD53(f, addr|access32bit, sdcard.Read, 4))
	return be.Decode32(sdcard.AsData(buf[:]).Bytes())
}

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
