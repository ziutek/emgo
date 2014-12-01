package onewire

// ROM commands.
const (
	searchROM   = 0xf0
	readROM     = 0x33
	matchROM    = 0x55
	skipROM     = 0xCC
	alarmSearch = 0xec
)

// DS18B20 function commands.
const (
	convertT = 0x44
)

func (m *Master) ReadROM() (Dev, error) {
	if err := m.Reset(); err != nil {
		return 0, err
	}
	if err := m.WriteByte(readROM); err != nil {
		return 0, err
	}
	var rom uint64
	for i := uint(0); i < 64; i += 8 {
		b, err := m.ReadByte()
		if err != nil {
			return 0, err
		}
		rom |= uint64(b) << i
	}
	return Dev(rom), nil
}
