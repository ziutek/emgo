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
