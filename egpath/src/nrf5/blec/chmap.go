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
		used += int(l & 1)
		l >>= 1
	}
	for h != 0 {
		used += int(h & 1)
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
	chi := uint(chm.uchi) + uint(chm.hop)
	if chi >= 37 {
		chi -= 37
	}
	chm.uchi = byte(chi)
	if chi < 32 {
		if chm.l&(1<<chi) != 0 {
			return int(chi)
		}
	} else {
		if chm.h&(1<<(chi-32)) != 0 {
			return int(chi)
		}
	}
	remapIdx := chi % uint(chm.used)
	chi = 0
	for chi < 32 {
		for chm.l&(1<<chi) == 0 {
			chi++
		}
		if remapIdx == 0 {
			return int(chi)
		}
		remapIdx--
		chi++
	}
	chi = 0
	for {
		for chm.h&(1<<chi) == 0 {
			chi++
		}
		if remapIdx == 0 {
			return int(chi + 32)
		}
		remapIdx--
		chi++
	}
}
