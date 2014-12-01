package onewire

func (m *Master) nextBit() (bit, neg byte, err error) {
	bit, err = m.RecvBit()
	if err != nil {
		return
	}
	neg, err = m.RecvBit()
	if err != nil {
		return
	}
	if bit == 1 && neg == 1 {
		err = ErrNoResponse
	}
	return
}

// Search perform a search for 1-Wire slave devices connected to bus controled by m.
// If alarm is true it performs a search for devices in alarm state.
// Search calls found function for any found device.
func (m *Master) Search(found func(Dev), alarm bool) error {
	searchCmd := byte(searchROM)
	if alarm {
		searchCmd = alarmSearch
	}
	var dev uint64
	lastd := -1
	for {
		if err := m.Reset(); err != nil {
			return err
		}
		if err := m.WriteByte(searchCmd); err != nil {
			return err
		}
		last0 := -1
		for i := 0; i < 64; i++ {
			bit, cmp, err := m.nextBit()
			if err != nil {
				return err
			}
			b := uint64(1) << uint(i)
			switch bit {
			case cmp:
				// discrepancy
				switch {
				case i > lastd:
					// Take 0-dir.
					bit = 0
					dev &^= b
				case i == lastd:
					// Take 1-dir.
					bit = 1
					dev |= b
				default: // i < lastd
					// Take previous dir.
					bit = byte((dev & b) >> uint(i))
				}
				if bit == 0 {
					last0 = i
				}
			case 0:
				dev &^= b
			case 1:
				dev |= b
			}
			if err := m.SendBit(bit); err != nil {
				return err
			}
		}

		found(Dev(dev))

		if last0 == -1 {
			break
		}
		lastd = last0
	}
	return nil
}

type Search struct {
	dev   uint64
	lastd int
	cmd   byte
	err   error
}

func (s *Search) Init(alarm bool) {
	s.dev = 0
	s.lastd = -1
	if alarm {
		s.cmd = alarmSearch
	} else {
		s.cmd = searchRom
	}
}

func (s *Search) Dev() Dev {
	return s.dev
}

// Search initializes search and returns
func (m *Master) Next(s *Search) bool {
	if s.err := m.Reset(); s.err != nil {
		return false
	}
	if s.err := m.WriteByte(s.cmd); s.err != nil {
		return false
	}
	last0 := -1
	for i := 0; i < 64; i++ {
		bit, cmp, err := m.nextBit()
		if err != nil {
			return err
		}
		b := uint64(1) << uint(i)
		switch bit {
		case cmp:
			// discrepancy
			switch {
			case i > s.lastd:
				// Take 0-dir.
				bit = 0
				dev &^= b
			case i == s.lastd:
				// Take 1-dir.
				bit = 1
				dev |= b
			default: // i < s.lastd
				// Take previous dir.
				bit = byte((dev & b) >> uint(i))
			}
			if bit == 0 {
				last0 = i
			}
		case 0:
			dev &^= b
		case 1:
			dev |= b
		}
		if err := m.SendBit(bit); err != nil {
			return err
		}
	}
	if last0 == -1 {
		return false
	}
	s.lastd = last0
	return true
}
