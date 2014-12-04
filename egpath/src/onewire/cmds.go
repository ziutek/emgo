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

func (m *Master) ReadROM() (d Dev, err error) {
	if err = m.Reset(); err != nil {
		return
	}
	if err = m.WriteByte(readROM); err != nil {
		return
	}
	for k := range d {
		d[k], err = m.ReadByte()
		if err != nil {
			return
		}
	}
	return
}
