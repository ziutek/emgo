package blec

import (
	"encoding/binary/le"
)

type llData []byte

func (d llData) AA() uint32 {
	return uint32(d[0]) | uint32(d[1])<<8 | uint32(d[2])<<16 | uint32(d[3])<<24
}

func (d llData) CRCInit() uint32 {
	return uint32(d[4]) | uint32(d[5])<<8 | uint32(d[6])<<16
}

// WinSize returns window size as number microseconds.
func (d llData) WinSize() uint32 {
	return uint32(d[7]) * 1250
}

// WinOffset returns window offset as number microseconds + 1250.
func (d llData) WinOffset() uint32 {
	return (uint32(d[8]) | uint32(d[9])<<8 + 1) * 1250
}

// Interval returns connection interval as number microseconds.
func (d llData) Interval() uint32 {
	return (uint32(d[10]) | uint32(d[11])<<8) * 1250
}

func (d llData) Latency() int {
	return int(d[12]) | int(d[13])<<8
}

// Timeout returns connection supervision timeout as number microseconds.
func (d llData) Timeout() uint32 {
	return (uint32(d[14]) | uint32(d[15])<<8) * 10000
}

var sca = [8]byte{
	(500<<19+999999)/1000000 - 8,
	(250<<19+999999)/1000000 - 8,
	(150<<19+999999)/1000000 - 8,
	(100<<19+999999)/1000000 - 8,
	(75<<19+999999)/1000000 - 8,
	(50<<19+999999)/1000000 - 8,
	(30<<19+999999)/1000000 - 8,
	(20<<19+999999)/1000000 - 8,
}

// SCA returns (maxSCAPPM<<19 + 1e6 - 1) / 1e6.
func (d llData) SCA() uint32 {
	return uint32(sca[d[21]>>5]) + 8
}

func (d llData) ChMapL() uint32 {
	return le.Decode32(d[16:])
}

func (d llData) ChMapH() byte {
	return d[20]
}

func (d llData) Hop() byte {
	return d[21] & 0x1F
}
