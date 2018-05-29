package be

func Decode16(s []byte) uint16 {
	return uint16(s[0])<<8 | uint16(s[1])
}

func Decode32(s []byte) uint32 {
	return uint32(s[0])<<24 | uint32(s[1])<<16 | uint32(s[2])<<8 | uint32(s[3])
}

func Decode64(s []byte) uint64 {
	h := uint(s[0])<<24 | uint(s[1])<<16 | uint(s[2])<<8 | uint(s[3])
	l := uint(s[4])<<24 | uint(s[5])<<16 | uint(s[6])<<8 | uint(s[7])
	return uint64(l) | uint64(h)<<32
}

func Encode16(s []byte, u uint16) {
	s[0] = byte(u >> 8)
	s[1] = byte(u)
}

func Encode32(s []byte, u uint32) {
	s[0] = byte(u >> 24)
	s[1] = byte(u >> 16)
	s[2] = byte(u >> 8)
	s[3] = byte(u)
}

func Encode64(s []byte, u uint64) {
	h := uint(u >> 32)
	s[0] = byte(h >> 24)
	s[1] = byte(h >> 16)
	s[2] = byte(h >> 8)
	s[3] = byte(h)
	l := uint(u)
	s[4] = byte(l >> 24)
	s[5] = byte(l >> 16)
	s[6] = byte(l >> 8)
	s[7] = byte(l)
}
