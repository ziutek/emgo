package blec

type chmap struct {
	l    uint32
	h    byte
	used byte
	hop  byte
	uchi byte

	instant int32
	newL    uint32
	newH    byte

	connEventCnt uint16 // Zero for first received connEvent.
}

func (chm *chmap) ConnEventCnt() uint16 {
	return chm.connEventCnt
}

func (chm *chmap) set(l uint32, h byte) {
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
	chm.instant = -1
}

func (chm *chmap) Init(l uint32, h, hop byte) {
	chm.set(l, h)
	chm.hop = hop
	chm.connEventCnt = 0xFFFF
}

func (chm *chmap) Update(l uint32, h byte, instant uint16) {
	chm.instant = int32(instant)
	chm.newL = l
	chm.newH = h
}

func (chm *chmap) NextChi() int {
	chm.connEventCnt++
	if int32(chm.connEventCnt) == chm.instant {
		chm.set(chm.newL, chm.newH)
	}
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
