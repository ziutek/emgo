package onewire

type Search struct {
	// Parameters:
	typ Type
	cmd byte

	// State:
	err   error
	dev   Dev
	lastd int8
}

func (s *Search) Reset() {
	s.err = nil
	s.dev = Dev{}
	s.lastd = -1
}

func (s *Search) Init(typ Type, alarm bool) {
	s.Reset()
	s.typ = typ
	if alarm {
		s.cmd = alarmSearch
	} else {
		s.cmd = searchROM
	}
}

func MakeSearch(typ Type, alarm bool) Search {
	var s Search
	s.Init(typ, alarm)
	return s
}

func (s *Search) Dev() Dev {
	return Dev(s.dev)
}

func (s *Search) Err() error {
	return s.err
}

func (m *Master) nextBit() (bit, neg byte, err error) {
	bit, err = m.ReadBit()
	if err != nil {
		return
	}
	neg, err = m.ReadBit()
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
	k := 0
	if s.typ != 0 {
		for i := 0; i < 8; i++ {
			nb := byte(s.typ>>uint(i)) & 1
			bit, cmp, err := m.nextBit()
			if err != nil {
				s.err = err
				return false
			}
			if bit != nb && bit != cmp {
				// There is no device with nb bit.
				return false
			}
			if s.err = m.WriteBit(nb); s.err != nil {
				return false
			}
		}
		s.dev[k] = byte(s.typ)
		k++
	}
	last0 := -1
	for k < len(s.dev) {
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
				// Discrepancy.
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
			default: // 1
				*dp |= b
			}
			if s.err = m.WriteBit(bit); s.err != nil {
				return false
			}
		}
		k++
	}
	if CRC8(0, s.dev[:7]) != s.dev[7] {
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
