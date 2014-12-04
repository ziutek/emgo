package onewire

// CRC8 with poly: x^8+x^5+x^4+1 -> 0x31==inv(0x8c)
func CRC8(crc byte, data ...byte) byte {
	for _, b := range data {
		for i := 0; i < 8; i++ {
			lsb := (crc ^ b) & 0x01
			crc >>= 1
			b >>= 1
			if lsb != 0 {
				crc ^= 0x8c
			}
		}
	}
	return crc
}

// CRC16 with poly: x^16+x^15+x^2+1 -> 0x4003==inv(0xc002<<1)
func CRC16(crc uint16, data ...uint16) uint16 {
	for _, w := range data {
		for i := 0; i < 16; i++ {
			lsb := (crc ^ w) & 0x01
			crc >>= 1
			w >>= 1
			if lsb != 0 {
				crc ^= 0xc002
			}
		}
	}
	return crc
}
