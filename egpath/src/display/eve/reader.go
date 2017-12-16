package eve

// Reader allows to read data from the EVE.
type Reader struct {
	d *Driver
}

// R starts a read transaction from the EVE memory at address addr. It returns
// Reader that provides set of reading methods.
func (d *Driver) R(addr int) Reader {
	d.end()
	checkAddr(addr)
	d.dci.Write([]byte{byte(addr >> 16), byte(addr >> 8), byte(addr)})
	d.dci.Read([]byte{0}) // Read dummy byte (input mode required by QSPI ).
	return Reader{d}
}

func (r Reader) ReadByte() byte {
	var buf [1]byte
	r.d.dci.Read(buf[:])
	return buf[0]
}

func (r Reader) ReadUint16() uint16 {
	var buf [2]byte
	r.d.dci.Read(buf[:])
	return uint16(buf[0]) | uint16(buf[1])<<8
}

func (r Reader) ReadUint32() uint32 {
	var buf [4]byte
	r.d.dci.Read(buf[:])
	return uint32(buf[0]) | uint32(buf[1])<<8 | uint32(buf[2])<<16 |
		uint32(buf[3])<<24
}

func (r Reader) ReadInt() int {
	return int(int32(r.ReadUint32()))
}

func (r Reader) Read(s []byte) (int, error) {
	r.d.dci.Read(s)
	if err := r.d.dci.Err(false); err != nil {
		return 0, err
	}
	return len(s), nil
}

// Close closes the read transaction.
func (r Reader) Close() {
	r.d.end()
}
