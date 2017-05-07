package blec

type LLData []byte

func (d LLData) AA() uint32 {
	return uint32(d[0]) | uint32(d[1])<<8 | uint32(d[2])<<16 | uint32(d[3])<<24
}

func (d LLData) CRCInit() uint32 {
	return uint32(d[4]) | uint32(d[5])<<8 | uint32(d[6])<<16
}

func (d LLData) WinSize() int {
	return int(d[7])
}

func (d LLData) WinOffset() int {
	return int(d[8]) | int(d[9])<<8
}

func (d LLData) Interval() int {
	return int(d[10]) | int(d[11])<<8
}

func (d LLData) Latency() int {
	return int(d[12]) | int(d[13])<<8
}

func (d LLData) Timeout() int {
	return int(d[14]) | int(d[15])<<8
}

func (d LLData) ChM() uint64 {
	cml := uint32(d[16]) | uint32(d[17])<<8 | uint32(d[18])<<16 |
		uint32(d[19])<<24
	return uint64(cml) | uint64(d[20])<<32
}

func (d LLData) Hop() int {
	return int(d[21]) & 0x1F
}

var scaPPM = [8]byte{500 / 2, 250 / 2, 150 / 2, 100 / 2, 76 / 2, 50 / 2, 30 / 2, 20 / 2}

func (d LLData) SCA() int {
	return int(scaPPM[d[21]>>5]) * 2
}
