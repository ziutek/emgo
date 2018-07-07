package main

import (
	"io"
)

type hexWriter struct {
	w     io.Writer
	toHex func(out, data []byte)
	ibuf  []byte
	obuf  []byte
	in    int
	ln    int
	lmax  int
}

func newHexWriter(w io.Writer, elsiz int, big bool) *hexWriter {
	hw := &hexWriter{
		w:    w,
		ibuf: make([]byte, elsiz),
		obuf: make([]byte, elsiz*2+5),
		lmax: 36 / (2 + elsiz),
	}
	if big {
		hw.toHex = bigHex
	} else {
		hw.toHex = littleHex
	}
	hw.obuf[0] = '\t'
	hw.obuf[1] = '0'
	hw.obuf[2] = 'x'
	hw.obuf[len(hw.obuf)-2] = ','
	hw.obuf[len(hw.obuf)-1] = '\n'
	return hw
}

func (hw *hexWriter) Write(data []byte) (n int, err error) {
	for n < len(data) {
		m := copy(hw.ibuf[hw.in:], data[n:])
		n += m
		hw.in += m
		if hw.in < len(hw.ibuf) {
			break
		}
		hw.in = 0
		hw.toHex(hw.obuf[3:len(hw.obuf)-2], hw.ibuf)
		on := len(hw.obuf)
		if hw.ln++; hw.ln < hw.lmax {
			on-- // No newline.
		}
		if _, err = hw.w.Write(hw.obuf[:on]); err != nil {
			break
		}
		switch hw.ln {
		case hw.lmax:
			hw.obuf[0] = '\t'
			hw.ln = 0
		case 1:
			hw.obuf[0] = ' '
		}
	}
	return n, err
}

var zeros [7]byte

func (hw *hexWriter) Flush() (err error) {
	if hw.in == 0 {
		_, err = hw.w.Write(hw.obuf[len(hw.obuf)-1:]) // Write newline.
	} else {
		copy(hw.ibuf[hw.in:], zeros[:])
		hw.in = 0
		hw.toHex(hw.obuf[3:len(hw.obuf)-2], hw.ibuf)
		_, err = hw.w.Write(hw.obuf) // Last number and newline.
	}
	return err
}

const digits = "0123456789ABCDEF"

func bigHex(out, data []byte) {
	for i := 0; i < len(data); i++ {
		b := data[i]
		out[i*2] = digits[b>>4]
		out[i*2+1] = digits[b&15]
	}
}

func littleHex(out, data []byte) {
	for i := 0; i < len(data); i++ {
		b := data[i]
		k := len(out) - i*2 - 2
		out[k] = digits[b>>4]
		out[k+1] = digits[b&15]
	}
}
