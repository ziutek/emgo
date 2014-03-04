package stlink

import "sync"

type size struct {
	buf byte
	tx  byte
	rx  byte
	_   byte
}

const bufLen = 80

// http://ncrmnt.org/wp/2013/05/06/stlink-as-a-serial-terminal/
type console struct {
	magic uint32
	siz   size
	txbuf [bufLen]byte
	rxbuf [bufLen]byte
	rd    int
}

func (c *console) BufLen() int {
	return int(c.siz.buf)
}

var Con console

func (c *console) Read(b []byte) (n int, _ error) {
	sync.Barrier()
	for c.siz.rx == 0 {
		sync.Barrier()
	}

	n = copy(b, c.rxbuf[c.rd:c.siz.rx])
	c.rd += n

	if c.rd == int(c.siz.rx) {
		sync.Memory()
		c.siz.rx = 0
		c.rd = 0
	}
	return
}

func (c *console) Write(b []byte) (n int, _ error) {
	for len(b) != 0 {
		sync.Barrier()
		for c.siz.tx != 0 {
			sync.Barrier()
		}

		m := copy(c.txbuf[:], b)
		sync.Memory()
		c.siz.tx = byte(m)

		n += m
		b = b[m:]
	}
	return
}

func init() {
	Con.magic = 0xDEADF00D - 1
	Con.siz.buf = bufLen
	Con.magic++
}
