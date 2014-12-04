package onewire

type Search struct {
	err   error
	dev   Dev
	lastd int8
	cmd   byte
}

func (s *Search) Reset() {
	s.dev = Dev{}
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
	for k := range s.dev {
		for i := 0; i < 8; i++ {
			bit, cmp, err := m.nextBit()
			if err != nil {
				s.err = err
				return false
			}
			dp := &s.dev[k]
			b := byte(1) << uint(i)
			switch bit {
			case cmp:
				// discrepancy
				switch {
				case i > int(s.lastd):
					// Take 0-dir.
					bit = 0
					*dp &^= b
				case i == int(s.lastd):
					// Take 1-dir.
					bit = 1
					*dp |= b
				default: // i < s.lastd
					// Take previous dir.
					bit = byte((*dp & b) >> uint(i))
				}
				if bit == 0 {
					last0 = i
				}
			case 0:
				*dp &^= b
			case 1:
				*dp |= b
			}
			if s.err = m.SendBit(bit); s.err != nil {
				return false
			}
		}
	}
	if CRC8(0, s.dev[:7]...) != s.dev[7] {
		s.err = ErrCRC
		return false
	}
	if last0 == -1 {
		s.lastd = -2 // No more devices.
	} else {
		s.lastd = int8(last0)
	}
	return true
}
