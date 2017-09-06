package blec

import (
	"bits"

	"nrf5/hal/ficr"
)

func getDevAddr() int64 {
	FICR := ficr.FICR
	l := FICR.DEVICEADDR[0].Load()
	h := FICR.DEVICEADDR[1].Load() & 0xFFFF
	if FICR.DEVICEADDRTYPE.Load()&1 != 0 {
		h |= 0x8000C000
	}
	return int64(h)<<32 | int64(l)
}

func decodeDevAddr(b []byte, random bool) int64 {
	if len(b) < 6 {
		return 0
	}
	l := uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
	h := uint32(b[4]) | uint32(b[5])<<8 | uint32(bits.One(random))<<31
	return int64(h)<<32 | int64(l)
}

type encodedPPM uint32

func encodePPM(ppm int) encodedPPM {
	return (encodedPPM(ppm)<<19 + 999999) / 1000000
}

func (e encodedPPM) Mul(v uint32) uint32 {
	return (v*uint32(e) + 1<<19 - 1) >> 19
}
