package blec

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

type chmap struct {
	l    uint32
	h    byte
	used byte
	hop  byte
	uchi byte
}

func (d llData) ChM() chmap {
	var chm chmap
	chm.l = uint32(d[16]) | uint32(d[17])<<8 | uint32(d[18])<<16 |
		uint32(d[19])<<24
	chm.h = d[20]
	chm.hop = d[21] & 0x1F
	used := 0
	for v := chm.l; v != 0; v >>= 1 {
		used++
	}
	for v := chm.h; v != 0; v >>= 1 {
		used++
	}
	chm.used = byte(used)
	return chm
}

func (chm *chmap) NextChi() int {
	uchi := uint(chm.uchi) + uint(chm.hop)
	if uchi >= 37 {
		uchi -= 37
	}
	chm.uchi = byte(uchi)
	if uchi < 32 {
		if chm.l&(1<<uchi) != 0 {
			return int(uchi)
		}
	} else {
		if chm.h&(1<<(uchi-32)) != 0 {
			return int(uchi)
		}
	}
	remapIdx := uchi % uint(chm.used)
	uchi = 0
	for uchi < 32 {
		for chm.l&(1<<uchi) == 0 {
			uchi++
		}
		if remapIdx == 0 {
			return int(uchi)
		}
		remapIdx--
		uchi++
	}
	uchi = 0
	for {
		for chm.h&(1<<uchi) == 0 {
			uchi++
		}
		if remapIdx == 0 {
			return int(uchi + 32)
		}
		remapIdx--
		uchi++
	}
}
