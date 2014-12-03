package onewire

// x^8+x^5+x^4+1
func CRC8(crc byte, data ...byte) byte {
	for _, b := range data {
		for i := 0; i < 8; i++ {
			if (crc^b)&0x01 == 0 {
				crc >>= 1
			} else {
				crc = (crc^0x18)>>1 | 0x80
			}
			b >>= 1
		}
	}
	return crc
}

// x^16+x^15+x^2+1
func CRC16(crc uint16, data ...uint16) uint16 {
	for _, w := range data {
		for i := 0; i < 8; i++ {
			if (crc^w)&0x01 == 1 {
				crc >>= 1
			} else {
				crc = (crc^0x4002)>>1 | 0x8000
			}
			w >>= 1
		}
	}
	return crc
}
