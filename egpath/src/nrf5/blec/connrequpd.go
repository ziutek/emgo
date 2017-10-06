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

type connUpdate struct {
	valid     bool
	winSize   byte
	winOffset uint16
	interval  uint16
	latency   uint16
	timeout   uint16
	instant   uint16
}

func (cu *connUpdate) Init(llConnUpdateReq []byte) {
	cu.valid = true
	cu.winSize = llConnUpdateReq[0]
	cu.winOffset = le.Decode16(llConnUpdateReq[1:3])
	cu.interval = le.Decode16(llConnUpdateReq[3:5])
	cu.latency = le.Decode16(llConnUpdateReq[5:7])
	cu.timeout = le.Decode16(llConnUpdateReq[7:9])
	cu.instant = le.Decode16(llConnUpdateReq[9:11])
}

func (cu *connUpdate) CheckInstant(connEventCnt uint16) bool {
	return cu.valid && connEventCnt == cu.instant
}

func (cu *connUpdate) SetDone() {
	cu.valid = false
}

// WinSize returns window size (µs).
func (cu *connUpdate) WinSize() uint32 {
	return uint32(cu.winSize) * 1250
}

// WinOffset returns window offset (µs).
func (cu *connUpdate) WinOffset() uint32 {
	return uint32(cu.winOffset) * 1250
}

// Interval returns connection interval (µs).
func (cu *connUpdate) Interval() uint32 {
	return uint32(cu.interval) * 1250
}

func (cu *connUpdate) Latency() int {
	return int(cu.latency)
}

// Timeout returns connection supervision timeout (µs).
func (cu *connUpdate) Timeout() uint32 {
	return uint32(cu.timeout) * 10000
}

func (cu *connUpdate) Instant() uint16 {
	return cu.instant
}
