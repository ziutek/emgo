// Package stlink provides STM32 ST-LINK debuging functions.
// This package is demonstration only (there are bugs in this package or in
// st-term that sometimes cause locking terminal) so better use one of your
// MCU's UART and serial package.
package stlink

import "sync/fence"

type size struct {
	buf byte
	tx  byte
	rx  byte
	_   byte
}

const TermBufLen = 64

// See http://ncrmnt.org/wp/2013/05/06/stlink-as-a-serial-terminal/
type terminal struct {
	magic uint32
	siz   size
	txbuf [TermBufLen]byte
	rxbuf [TermBufLen]byte
	rd    int
}

var (
	magick = uint32(0xDEADF00D - 1)
	Term   = &terminal{magic: magick + 1, siz: size{buf: TermBufLen}}
)

func (t *terminal) BufLen() int {
	return int(t.siz.buf)
}

func (t *terminal) waitrx() {
	fence.Compiler()
	for t.siz.rx == 0 {
		fence.Compiler()
	}
}

func (t *terminal) Read(b []byte) (n int, _ error) {
	t.waitrx()

	n = copy(b, t.rxbuf[t.rd:t.siz.rx])
	t.rd += n

	if t.rd == int(t.siz.rx) {
		fence.Memory()
		t.siz.rx = 0
		t.rd = 0
	}
	return
}

func (t *terminal) waittx() {
	fence.Compiler()
	for t.siz.tx != 0 {
		fence.Compiler()
	}
}

func (t *terminal) Write(b []byte) (n int, _ error) {
	for len(b) != 0 {
		t.waittx()
		m := copy(t.txbuf[:], b)
		fence.Memory()
		t.siz.tx = byte(m)

		n += m
		b = b[m:]
	}
	return
}

func (t *terminal) WriteString(s string) (n int, _ error) {
	//b := (*[]byte)(unsafe.Pointer(&s))
	//return t.Write((*b)[:len(s):len(s)])

	for len(s) != 0 {
		t.waittx()
		m := copy(t.txbuf[:], s)
		fence.Memory()
		t.siz.tx = byte(m)

		n += m
		s = s[m:]
	}
	return
}

func (t *terminal) WriteByte(b byte) (_ error) {
	t.Write([]byte{b})
	return
}
