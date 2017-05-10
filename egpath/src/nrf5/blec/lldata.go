package blec

type llData []byte

func (d llData) AA() uint32 {
	return uint32(d[0]) | uint32(d[1])<<8 | uint32(d[2])<<16 | uint32(d[3])<<24
}

func (d llData) CRCInit() uint32 {
	return uint32(d[4]) | uint32(d[5])<<8 | uint32(d[6])<<16
}

func (d llData) WinSize() uint32 {
	return uint32(d[7])
}

func (d llData) WinOffset() uint32 {
	return uint32(d[8]) | uint32(d[9])<<8
}

func (d llData) Interval() uint32 {
	return uint32(d[10]) | uint32(d[11])<<8
}

func (d llData) Latency() uint32 {
	return uint32(d[12]) | uint32(d[13])<<8
}

func (d llData) Timeout() uint32 {
	return uint32(d[14]) | uint32(d[15])<<8
}

func (d llData) ChM() uint64 {
	cml := uint32(d[16]) | uint32(d[17])<<8 | uint32(d[18])<<16 |
		uint32(d[19])<<24
	return uint64(cml) | uint64(d[20])<<32
}

func (d llData) Hop() int {
	return int(d[21]) & 0x1F
}

var ssca = [8]byte{
	(500<<19+999999)/1000000 - 8,
	(250<<19+999999)/1000000 - 8,
	(150<<19+999999)/1000000 - 8,
	(100<<19+999999)/1000000 - 8,
	(75<<19+999999)/1000000 - 8,
	(50<<19+999999)/1000000 - 8,
	(30<<19+999999)/1000000 - 8,
	(20<<19+999999)/1000000 - 8,
}

// SSCA returns (maxSCAPPM<<19 + 1e6 - 1) / 1e6.
func (d llData) SSCA() uint32 {
	return uint32(ssca[d[21]>>5]) + 8
}
