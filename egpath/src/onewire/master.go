package onewire

type Master struct {
	Driver
}

// WriteByte writes byte b to 1-wire bus.
func (m *Master) WriteByte(b byte) error {
	for i := 0; i < 8; i++ {
		if err := m.WriteBit(b & 1); err != nil {
			return err
		}
		b >>= 1
	}
	return nil
}

// ReadByte reads byte from 1-wire bus.
func (m *Master) ReadByte() (byte, error) {
	var b int
	for i := uint(0); i < 8; i++ {
		bit, err := m.ReadBit()
		if err != nil {
			return 0, err
		}
		b += int(bit) << i
	}
	return byte(b), nil
}

func (m *Master) Write(data []byte) (int, error) {
	for n, b := range data {
		if err := m.WriteByte(b); err != nil {
			return n, err
		}
	}
	return len(data), nil
}

func (m *Master) ReadFull(data []byte) (int, error) {
	for n := range data {
		b, err := m.ReadByte()
		if err != nil {
			return n, err
		}
		data[n] = b
	}
	return len(data), nil

}
