package eve

// Reader allows to read data from the EVE.
type Reader Driver

func (r *Reader) ReadByte() byte {
	var buf [1]byte
	r.dci.Read(buf[:])
	return buf[0]
}

func (r *Reader) ReadWord32() uint32 {
	var buf [4]byte
	r.dci.Read(buf[:])
	return uint32(buf[0]) | uint32(buf[1])<<8 | uint32(buf[2])<<16 |
		uint32(buf[3])<<24
}

func (r *Reader) Read(s []byte) (int, error) {
	r.dci.Read(s)
	if err := r.dci.Err(false); err != nil {
		return 0, err
	}
	return len(s), nil
}
