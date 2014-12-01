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
	err   error
	lastd int8
	cmd   byte
}

func (s *Search) Reset() {
	s.dev = 0
	s.err = nil
	s.lastd = -1
}

func (s *Search) Init(alarm bool) {
	s.Reset()
	if alarm {
		s.cmd = alarmSearch
	} else {
		s.cmd = searchROM
	}
}

func MakeSearch(alarm bool) Search {
	var s Search
	s.Init(alarm)
	return s
}

func (s *Search) Dev() Dev {
	return Dev(s.dev)
}

func (s *Search) Err() error {
	return s.err
}

// SearchNext can be used to perform a search for 1-Wire slave devices
// connected to the bus controled by m.
// It saves current state in s and returns true if next device was found or
// false if no more devices or error occurred. Use s.Dev() to get the device
// found, s.Err() to check for error.  
func (m *Master) SearchNext(s *Search) bool {
	if s.lastd == -2 {
		return false
	}
	if s.err = m.Reset(); s.err != nil {
		return false
	}
	if s.err = m.WriteByte(s.cmd); s.err != nil {
		return false
	}
	last0 := -1
	for i := 0; i < 64; i++ {
		bit, cmp, err := m.nextBit()
		if err != nil {
			s.err = err
			return false
		}
		b := uint64(1) << uint(i)
		switch bit {
		case cmp:
			// discrepancy
			switch {
			case i > int(s.lastd):
				// Take 0-dir.
				bit = 0
				s.dev &^= b
			case i == int(s.lastd):
				// Take 1-dir.
				bit = 1
				s.dev |= b
			default: // i < s.lastd
				// Take previous dir.
				bit = byte((s.dev & b) >> uint(i))
			}
			if bit == 0 {
				last0 = i
			}
		case 0:
			s.dev &^= b
		case 1:
			s.dev |= b
		}
		if s.err = m.SendBit(bit); s.err != nil {
			return false
		}
	}
	if last0 == -1 {
		s.lastd = -2 // No more devices.
	} else {
		s.lastd = int8(last0)
	}
	return true
}
