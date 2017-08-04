package att

func panicShort() {
	panic("att: to short")
}

func Decode16(s []byte) uint16 {
	if len(s) < 2 {
		panicShort()
	}
	return uint16(s[0]) | uint16(s[1])<<8
}

func Decode32(s []byte) uint32 {
	if len(s) < 4 {
		panicShort()
	}
	return uint32(s[0]) | uint32(s[1])<<8 | uint32(s[2])<<16 | uint32(s[3])<<24
}

func Decode64(s []byte) uint64 {
	if len(s) < 8 {
		panicShort()
	}
	l := uint(s[0]) | uint(s[1])<<8 | uint(s[2])<<16 | uint(s[3])<<24
	h := uint(s[4]) | uint(s[5])<<8 | uint(s[6])<<16 | uint(s[7])<<24
	return uint64(l) | uint64(h)<<32
}

func Encode16(s []byte, u uint16) {
	if len(s) < 2 {
		panicShort()
	}
	s[0] = byte(u)
	s[1] = byte(u >> 8)
}

func Encode32(s []byte, u uint32) {
	if len(s) < 4 {
		panicShort()
	}
	s[0] = byte(u)
	s[1] = byte(u >> 8)
	s[2] = byte(u >> 16)
	s[3] = byte(u >> 24)
}

func Encode64(s []byte, u uint64) {
	if len(s) < 8 {
		panicShort()
	}
	l := uint(u)
	s[0] = byte(l)
	s[1] = byte(l >> 8)
	s[2] = byte(l >> 16)
	s[3] = byte(l >> 24)
	h := uint(u >> 32)
	s[4] = byte(h)
	s[5] = byte(h >> 8)
	s[6] = byte(h >> 16)
	s[7] = byte(h >> 24)
}
