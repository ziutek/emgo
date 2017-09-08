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

// Fixed19 is 32-bit binary fixed-point unsigned number with scaling factor
// 1/2^19 = 1/524288 = 0.0000019073486328125. It divides 32-bit word to 13-bit
// integer part and 19-bit fractional part. Fixed19 can store values from 0 to
// 8191.9999980926513671875.
type fixed19 uint32

// ppmToFixedUp converts ppm (Parts Per Million) to fixed19 (rounding mode: up).
// Allowed ppm values: 0 <= ppm <= 8191.
func ppmToFixedUp(ppm int) fixed19 {
	return (fixed19(ppm)<<19 + 999999) / 1000000
}

// MulUp multiples integer value v by x and returns integer value rounded up.
func (x fixed19) MulUp(v uint32) uint32 {
	return (v*uint32(x) + 1<<19 - 1) >> 19
}
