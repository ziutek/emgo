package onewire

// Types of DS18x20 family.
const (
	DS18S20 Type = 0x10
	DS1822  Type = 0x22
	DS18B20 Type = 0x28
)

// DS18x2x specific commands.
const (
	convertT        = 0x44
	readScratchpad  = 0xbe
	writeScratchpad = 0x4e
	copyScratchpad  = 0x48
	recallEE        = 0xb8
	readPowerSupply = 0xb4
)

// ConvertT (DS18x2x) initiates temperature conversion.
// Temperature conversion needs extra power. In case of parasite power mode
// the strong pullup is required during the entire conversion so no any other
// command can't be sent before the duration of tconv. If external supply is
// used, ReadBit method can be used to pooling conversion state: it returns 0
// while the temperature conversion is in progress and a 1 when the conversion
// is done.
func (m *Master) ConvertT() error {
	return m.WriteByte(convertT)
}

const (
	T9bit  = 0x1f
	T10bit = 0x3f
	T11bit = 0x5f
	T12bit = 0x7f
)

// WriteScratchpad writes data to scratchpad. DS18S20 ignores cfg byte.
func (m *Master) WriteScratchpad(th, tl int8, cfg byte) error {
	if err := m.WriteByte(writeScratchpad); err != nil {
		return err
	}
	_, err := m.Write([]byte{byte(th), byte(tl), cfg})
	return err
}

// Scratchpad represents 9 bytes of data that can be read from DS18x2x device.
type Scratchpad [9]byte

// Temp16 returns temperature data from s as uint16 translated to uint16 value
// equal to T[C] * 16.
func (s *Scratchpad) Temp16(typ Type) (int, error) {
	switch typ {
	case DS18B20, DS1822:
		return int(uint(s[1])<<8 + uint(s[0])), nil
	}
	return 0x1000, ErrDevType
}

// Temp returns temperature data from s as float32 value T[C].
func (s *Scratchpad) Temp(typ Type) (float32, error) {
	t16, err := s.Temp16(typ)
	return float32(t16) * 0.0625, err
}

// CRC returns value of CRC field. It should be equal to CRC8(0, s[:8]).
func (s *Scratchpad) CRC() byte {
	return s[8]
}

// ReadScratchpad (DS18x2x) reads content of scratchpad and checks its CRC.
// It retuns readed data and error (ErrCRC in case of bad CRC).
func (m *Master) ReadScratchpad() (s Scratchpad, err error) {
	if err = m.WriteByte(readScratchpad); err != nil {
		return
	}
	if _, err = m.ReadFull(s[:]); err != nil {
		return
	}
	if CRC8(0, s[:8]) != s.CRC() {
		err = ErrCRC
	}
	return
}

// CopyScratchpad copies the contents of the scratchpad Th, Tl and configuration
// registers to EEPROM. This command needs extra power.
func (m *Master) CopyScratchpad() error {
	return m.WriteByte(copyScratchpad)
}

// RecallEE
func (m *Master) RecallEE() error {
	return m.WriteByte(recallEE)
}

// ReadPowerSupply
func (m *Master) ReadPowerSupply() error {
	return m.WriteByte(readPowerSupply)
}
