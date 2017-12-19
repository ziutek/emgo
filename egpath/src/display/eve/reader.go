package eve

// Reader allows to read data from the EVE.
type Reader struct {
	d    *Driver
	addr int
}

func (r *Reader) addrAdd(n int) {
	r.addr += n
	r.d.state.addr = r.addr
}

// Addr returns the current read address.
func (r *Reader) Addr() int {
	return r.addr
}

func (r *Reader) start() {
	d := r.d
	buf := [...]byte{byte(r.addr >> 16), byte(r.addr >> 8), byte(r.addr)}
	d.dci.Write(buf[:])
	d.dci.Read(buf[:1]) // Read dummy byte (input mode required by QSPI ).
}

// R starts a read transaction from the EVE memory at address addr. It returns
// Reader that provides set of reading methods.
func (d *Driver) R(addr int) Reader {
	d.end()
	checkAddr(addr)
	d.state.flags = stateOpen | stateRead
	r := Reader{d, addr}
	r.start()
	return r
}

func (r *Reader) restart(n int) {
	d := r.d
	if (d.state != state{r.addr, stateOpen | stateRead}) {
		if d.state.flags&stateOpen != 0 {
			d.end()
		}
		r.start()
		r.addr += n
		d.state = state{r.addr, stateOpen | stateRead}
	} else {
		r.addrAdd(n)
	}
}

func (r *Reader) ReadByte() byte {
	r.restart(1)
	var buf [1]byte
	r.d.dci.Read(buf[:])
	return buf[0]
}

func (r *Reader) ReadUint16() uint16 {
	r.restart(2)
	var buf [2]byte
	r.d.dci.Read(buf[:])
	return uint16(buf[0]) | uint16(buf[1])<<8
}

func (r *Reader) ReadUint32() uint32 {
	r.restart(4)
	var buf [4]byte
	r.d.dci.Read(buf[:])
	return uint32(buf[0]) | uint32(buf[1])<<8 | uint32(buf[2])<<16 |
		uint32(buf[3])<<24
}

func (r *Reader) ReadInt() int {
	return int(int32(r.ReadUint32()))
}

func (r *Reader) Read(s []byte) (int, error) {
	r.restart(len(s))
	r.d.dci.Read(s)
	if err := r.d.dci.Err(false); err != nil {
		return 0, err
	}
	return len(s), nil
}

// Sync closes the read transaction.
func (r *Reader) Sync() {
	r.d.end()
}
