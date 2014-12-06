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
	recallE         = 0xb8
	readPowerSupply = 0xb4
)

// ConvertT (DS18x2x) initiates temperature conversion. In case of parasite
// power mode any other command can't be sent before the duration of tconv.
// If external supply is used, ReadBit method can be used (periodicaly): it
// returns 0 while the temperature conversion is in progress and a 1 when the
// conversion is done.
func (m *Master) ConvertT() error {
	return m.WriteByte(convertT)
}

func (m *Master) WriteScratchpad(th, tl, cfg byte) error {
	if err := m.WriteByte(writeScratchpad); err != nil {
		return err
	}
	_, err := m.Write([]byte{th, tl, cfg})
	return err
}

type Scratchpad [9]byte

func (s *Scratchpad) Temp16(typ Type) (int, error) {
	switch typ {
	case DS18B20, DS1822:
		return int(uint(s[1])<<8 + uint(s[0])), nil
	}
	return 0x1000, ErrDevType
}

func (s *Scratchpad) Temp(typ Type) (float32, error) {
	t16, err := s.Temp16(typ)
	return float32(t16) * 0.0625, err
}

func (s *Scratchpad) CRC() byte {
	return s[8]
}

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
