package blec

type chmap struct {
	l    uint32
	h    byte
	used byte
	hop  byte
	uchi byte
}

func (chm *chmap) SetLH(l uint32, h byte) {
	chm.l = l
	chm.h = h
	used := 0
	for l != 0 {
		used++
		l >>= 1
	}
	for h != 0 {
		used++
		h >>= 1
	}
	chm.used = byte(used)
}

func (chm *chmap) SetHop(hop byte) {
	chm.hop = hop
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
