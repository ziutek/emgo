package uuid

import (
	"encoding/binary/le"
	"errors"
	"fmt"
)

// Long is full (128-bit) universally unique identifier. H and L represent
// respectively the most and the less significant bytes of UUID (in 8-4-4-4-12
// text format: HHHHHHHH-HHHH-HHHH-LLLL-LLLLLLLLLLLL).
type Long struct {
	H, L uint64
}

// DecodeLong decodes UUID from first 16 bytes of s.
func DecodeLong(s []byte) Long {
	return Long{le.Decode64(s[8:]), le.Decode64(s)}
}

var ErrBadUUID = errors.New("bad UUID")

// Parse parses text representation of 128-bit UUID. It requires 8-4-4-4-12
// format (len(s) must be 36).
func Parse(s []byte) (u Long, err error) {
	if len(s) != 36 {
		return u, ErrBadUUID
	}
	n := uint(128)
	for i := 0; i < 36; i++ {
		d := int(s[i])
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if d != '-' {
				return Long{}, ErrBadUUID
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
			return Long{}, ErrBadUUID
		}
		if n -= 4; n < 64 {
			u.L |= uint64(d) << n
		} else {
			u.H |= uint64(d) << (n - 64)
		}
	}
	return u, nil
}

func (u Long) Encode(s []byte) {
	le.Encode64(s, u.L)
	le.Encode64(s[8:], u.H)
}

// Format produces 8-4-4-4-12 text format of u.
func (u Long) Format(f fmt.State, c rune) {
	var buf [36]byte
	n := uint(128)
	for i := range buf {
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
	f.Write(buf[:])
}

// Base is the Base Bluetooth UUID, used for calculating 128-bit UUIDs from
// shortened (16-bit, 32-bit) UUIDs. Bluetooth UUIDs follow the template
// xxxxxxxx-0000-1000-8000-00805F9B34FB so BaseUUID == LongUUID(0).
//
//emgo:const
var Base = Long{0x1000, 0x800000805F9B34FB}

// Short is shortened (16-bit) bluetooth universally unique identifier. See Base
// for more information about shortened form of UUID.
type Short uint16

// DecodeShort decodes short UUID from first 2 bytes of s.
func DecodeShort(s []byte) Short {
	return Short(le.Decode16(s))
}

func (u Short) Long() Long {
	return Long{Base.H | uint64(u)<<32, Base.L}
}

func (u Short) Encode(s []byte) {
	le.Encode16(s, uint16(u))
}
