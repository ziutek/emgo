package uuid

import (
	"encoding/binary/le"
	"errors"
	"fmt"
)

// UUID is full (128-bit) universally unique identifier. H and L represent
// respectively the most and the less significant bytes of UUID (in 8-4-4-4-12
// text format: HHHHHHHH-HHHH-HHHH-LLLL-LLLLLLLLLLLL).
type UUID struct {
	H, L uint64
}

// Decode decodes UUID from first 16 bytes of s.
func Decode(s []byte) UUID {
	return UUID{le.Decode64(s[8:]), le.Decode64(s)}
}

var ErrBadUUID = errors.New("bad UUID")

// Parse parses text representation of 128-bit UUID. It requires 8-4-4-4-12
// format (len(s) must be 36).
func Parse(s []byte) (u UUID, err error) {
	if len(s) != 36 {
		return u, ErrBadUUID
	}
	n := uint(128)
	for i := 0; i < 36; i++ {
		d := int(s[i])
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if d != '-' {
				return UUID{}, ErrBadUUID
			}
			continue
		}
		switch {
		case d >= '0' && d <= '9':
			d -= '0'
		case d >= 'A' && d <= 'F':
			d -= 'A' - 10
		case d >= 'a' && d <= 'f':
			d -= 'a' - 10
		default:
			return UUID{}, ErrBadUUID
		}
		if n -= 4; n < 64 {
			u.L |= uint64(d) << n
		} else {
			u.H |= uint64(d) << (n - 64)
		}
	}
	return u, nil
}

func (u UUID) Encode(s []byte) {
	le.Encode64(s, u.L)
	le.Encode64(s[8:], u.H)
}

func (u UUID) CanShorten(bits int) bool {
	return u.L == Base.L && u.H&0xFFFFFFFF == Base.H &&
		uint(u.H>>32)>>uint(bits) == 0
}

func (u UUID) Short16() UUID16 {
	if u.CanShorten(16) {
		return UUID16(u.H >> 32)
	}
	panic("uuid: can not shorten")
}

func (u UUID) Short32() UUID32 {
	if u.CanShorten(32) {
		return UUID32(u.H >> 32)
	}
	panic("uuid: can not shorten")
}

// Format produces 8-4-4-4-12 text format of u.
func (u UUID) Format(f fmt.State, c rune) {
	var buf [36]byte
	blen := len(buf)
	if c == 'v' && u.CanShorten(32) {
		blen = 8
	}
	for i, n := 0, uint(128); i < blen; i++ {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			buf[i] = '-'
			continue
		}
		var d int
		if n -= 4; n < 64 {
			d = int(u.L>>n) & 0xF
		} else {
			d = int(u.H>>(n-64)) & 0xF
		}
		switch {
		case d < 10:
			d += '0'
		case c == 'X':
			d += 'A' - 10
		default:
			d += 'a' - 10
		}
		buf[i] = byte(d)
	}
	if blen < len(buf) {
		copy(buf[blen:], "-bluetooth")
	}
	f.Write(buf[:])
}

// Base is the Base Bluetooth UUID, used for calculating 128-bit UUIDs from
// shortened (16-bit, 32-bit) UUIDs. Bluetooth UUIDs follow the template
// xxxxxxxx-0000-1000-8000-00805F9B34FB so Short(0).UUID() == Base.
//
//emgo:const
var Base = UUID{0x1000, 0x800000805F9B34FB}

// UUID16 is shortened (16-bit) bluetooth universally unique identifier. See
// Base for more information about shortened form of UUID.
type UUID16 uint16

// DecodeShort decodes short UUID from first 2 bytes of s.
func DecodeShort(s []byte) UUID16 {
	return UUID16(le.Decode16(s))
}

func (u UUID16) Full() UUID {
	return UUID{Base.H | uint64(u)<<32, Base.L}
}

func (u UUID16) Encode(s []byte) {
	le.Encode16(s, uint16(u))
}

// UUID32 is shortened (32-bit) bluetooth universally unique identifier. See
// Base for more information about shortened form of UUID.
type UUID32 uint32

func (u UUID32) Full() UUID {
	return UUID{Base.H | uint64(u)<<32, Base.L}
}
