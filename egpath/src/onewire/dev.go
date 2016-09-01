package onewire

import (
	"fmt"
)

type Dev [8]byte

func (d Dev) Format(f fmt.State, _ rune) {
	typ := d[0]
	crc := d[7]
	a1 := uint32(d[1])<<16 + uint32(d[2])<<8 + uint32(d[3])
	a2 := uint32(d[4])<<16 + uint32(d[5])<<8 + uint32(d[6])
	fmt.Fprintf(f, "%02x-%06x-%06x-%02x", typ, a1, a2, crc)
}

type Type byte

func (d Dev) Type() Type {
	return Type(d[0])
}