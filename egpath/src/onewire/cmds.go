package onewire

// Generic ROM commands.
const (
	searchROM   = 0xf0
	readROM     = 0x33
	matchROM    = 0x55
	skipROM     = 0xcc
	alarmSearch = 0xec
)

// ReadROM allows the bus master to read the slave’s 64-bit ROM code without
// using the SearcROM method. It can only be used when there is only one slave
// device on the bus.
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

// MatchROM allows the bus master to address a specific slave device.
func (m *Master) MatchROM(d Dev) error {
	if err := m.Reset(); err != nil {
		return err
	}
	if err := m.WriteByte(matchROM); err != nil {
		return err
	}
	for _, b := range d {
		if err := m.WriteByte(b); err != nil {
			return err
		}
	}
	return nil
}

// SkipROM can be used to address all devices on the bus simultaneously.
func (m *Master) SkipROM(d Dev) error {
	if err := m.Reset(); err != nil {
		return err
	}
	if err := m.WriteByte(skipROM); err != nil {
		return err
	}
	return nil
}
