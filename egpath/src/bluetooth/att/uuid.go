package att

import (
	"errors"
	"fmt"
)

// UUID is full (128-bit) universally unique identifier. H and L represent
// respectively the most and the less significant bytes of UUID (in 8-4-4-4-12
// text format: HHHHHHHH-HHHH-HHHH-LLLL-LLLLLLLLLLLL).
type UUID struct {
	H, L uint64
}

// DecodeUUID decodes UUID from first 16 bytes of s.
func DecodeUUID(s []byte) (u UUID) {
	u.L = Decode64(s)
	u.H = Decode64(s[8:])
	return
}

var ErrBadUUID = errors.New("bad UUID")

// ParseUUID parses text representation of 128-bit UUID. It requires 8-4-4-4-12
// format (len(s) must be 36).
func ParseUUID(s []byte) (u UUID, err error) {
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

func (u UUID) Encode(s []byte) int {
	Encode64(s, u.L)
	Encode64(s[8:], u.H)
	return 16
}

// Format produces 8-4-4-4-12 text format of u.
func (u UUID) Format(f fmt.State, c rune) {
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

// BaseUUID is the Base Bluetooth UUID, used for calculating 128-bit UUIDs from
// shortened (16-bit, 32-bit) UUIDs. Bluetooth UUIDs follow the template
// xxxxxxxx-0000-1000-8000-00805F9B34FB so BaseUUID == FullUUID(0).
//
//emgo:const
var BaseUUID = UUID{0x1000, 0x800000805F9B34FB}

// UUID16 is shortened (16-bit) bluetooth universally unique identifier.
// See BaseUUID for more information about shortened form of UUID.
type UUID16 uint16

// DecodeUUID16 decodes UUID16 from first 2 bytes of s.
func DecodeUUID16(s []byte) UUID16 {
	return UUID16(Decode16(s))
}

func (u UUID16) Full() UUID {
	return UUID{BaseUUID.H | uint64(u)<<32, BaseUUID.L}
}

func (u UUID16) Encode(s []byte) int {
	Encode16(s, uint16(u))
	return 2
}

// UUID32 is shortened (32-bit) bluetooth universally unique identifier.
// See BaseUUID for more information about shortened form of UUID.
type UUID32 uint32

// DecodeUUID32 decodes UUID16 from first 4 bytes of s.
func DecodeUUID32(s []byte) UUID32 {
	return UUID32(Decode32(s))
}

func (u UUID32) Full() UUID {
	return UUID{BaseUUID.H | uint64(u)<<32, BaseUUID.L}
}

func (u UUID32) Encode(s []byte) int {
	Encode32(s, uint32(u))
	return 4
}
