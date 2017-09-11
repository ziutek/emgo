package blec

import (
	"encoding/binary/le"
)

type connReqLLData []byte

func (d connReqLLData) AA() uint32 {
	return le.Decode32(d[0:4])
}

func (d connReqLLData) CRCInit() uint32 {
	return uint32(d[4]) | uint32(d[5])<<8 | uint32(d[6])<<16
}

// WinSize returns window size (µs).
func (d connReqLLData) WinSize() uint32 {
	return uint32(d[7]) * 1250
}

// WinOffset returns final window offset (µs) (WinOffset field + 1250 µs).
func (d connReqLLData) WinOffset() uint32 {
	return (uint32(le.Decode16(d[8:10])) + 1) * 1250
}

// Interval returns connection interval (µs).
func (d connReqLLData) Interval() uint32 {
	return uint32(le.Decode16(d[10:12])) * 1250
}

func (d connReqLLData) Latency() int {
	return int(le.Decode16(d[12:14]))
}

// Timeout returns connection supervision timeout (µs).
func (d connReqLLData) Timeout() uint32 {
	return uint32(le.Decode16(d[14:16])) * 10000
}

func (d connReqLLData) ChMapL() uint32 {
	return le.Decode32(d[16:])
}

func (d connReqLLData) ChMapH() byte {
	return d[20]
}

func (d connReqLLData) Hop() byte {
	return d[21] & 0x1F
}

//emgo:const
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

// SCA returns master's Sleep Clock Accuracy as fixed19 number.
func (d connReqLLData) SCA() fixed19 {
	return fixed19(sca[d[21]>>5]) + 8
}

type llConnUpdate []byte

// WinSize returns window size (µs).
func (d llConnUpdate) WinSize() uint32 {
	return uint32(d[0]) * 1250
}

// WinOffset returns window offset (µs).
func (d llConnUpdate) WinOffset() uint32 {
	return (uint32(le.Decode16(d[1:3]))) * 1250
}

// Interval returns connection interval (µs).
func (d llConnUpdate) Interval() uint32 {
	return uint32(le.Decode16(d[3:5])) * 1250
}

func (d llConnUpdate) Latency() int {
	return int(le.Decode16(d[5:7]))
}

// Timeout returns connection supervision timeout (µs).
func (d llConnUpdate) Timeout() uint32 {
	return uint32(le.Decode16(d[7:9])) * 10000
}

func (d llConnUpdate) Instant() uint16 {
	return le.Decode16(d[9:11])
}
