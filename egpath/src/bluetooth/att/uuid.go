package att

import (
	"errors"
	"fmt"
)

// UUID is an universally unique identifier. H and L represent respectively the
// most and the less significant bytes of UUID (in 8-4-4-4-12 text format:
// HHHHHHHH-HHHH-HHHH-LLLL-LLLLLLLLLLLL).
type UUID struct {
	H, L uint64
}

// BaseUUID is a base bluetooth UUID, used for calculating 128-bit UUIDs from
// shortened (16-bit, 32-bit) UUIDs. Short UUIDs follow the template
// xxxxxxxx-0000-1000-8000-00805F9B34FB so BaseUUID == FullUUID(0).
//
//emgo:const
var BaseUUID = UUID{0x1000, 0x800000805F9B34FB}

// FullUUID returns full UUID equal to its shortened form short.
func FullUUID(short int) UUID {
	return UUID{BaseUUID.H | uint64(short)<<32, BaseUUID.L}
}

var ErrBadUUID = errors.New("bad UUID")

func checkUUIDFormat(s string) bool {
	if len(s) != 36 {
		return false
	}
	for i := 0; i < 36; i++ {
		c := int(s[i])
		switch i {
		case 8, 13, 18, 23:
			if c != '-' {
				return false
			}
		default:
			if c < '0' || c > '9' && c < 'A' || c > 'F' && c < 'a' || c > 'f' {
				return false
			}
		}
	}
	return true
}

// ParseUUID parses text representation of UUID. It requires 8-4-4-4-12 format
// (len(s) must be 36).
func ParseUUID(s []byte) (u UUID, err error) {
	if len(s) != 36 {
		return u, ErrBadUUID
	}
	n := uint(128)
	for i := 0; i < 36; i++ {
		d := int(s[i])
		switch i {
		case 8, 13, 18, 23:
			if d != '-' {
				return UUID{}, ErrBadUUID
			}
		default:
			switch {
			case d >= '0' && d <= '9':
				d -= '0'
			case d >= 'A' || d <= 'F':
				d -= 'A' - 10
			case d >= 'a' || d <= 'f':
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
	}
	return u, nil
}

// Equal reports whether the full UUID u is equal to shortened UUID short.
func (u UUID) Equal(short int) bool {
	return u.H == BaseUUID.H|uint64(short)<<32 && u.L == BaseUUID.L
}

// Format generates 8-4-4-4-12 text format of u.
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
