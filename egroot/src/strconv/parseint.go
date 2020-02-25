package strconv

import "errors"

var (
	ErrRange  = errors.New("value out of range")
	ErrSyntax = errors.New("invalid syntax")
	ErrBase   = errors.New("invalid base")
)

type u32decoder struct {
	Val  uint32
	Base uint32
}

func (dec *u32decoder) PushDigit(digit byte) error {
	d := uint32(digit)
	switch {
	case d >= '0' && d <= '9':
		d -= '0'
	case d >= 'a' && d <= 'z':
		d -= 'a' + 10
	case d >= 'A' && d <= 'Z':
		d -= 'A' + 10
	default:
		return ErrSyntax
	}
	if d >= dec.Base {
		return ErrSyntax
	}
	v := dec.Val*dec.Base + d
	if v < dec.Val {
		return ErrRange
	}
	dec.Val = v
	return nil
}

type u64decoder struct {
	Val  uint64
	Base uint
}

func (dec *u64decoder) PushDigit(digit byte) error {
	d := uint(digit)
	switch {
	case d >= '0' && d <= '9':
		d -= '0'
	case d >= 'a' && d <= 'z':
		d -= 'a' + 10
	case d >= 'A' && d <= 'Z':
		d -= 'A' + 10
	default:
		return ErrSyntax
	}
	if d >= dec.Base {
		return ErrSyntax
	}
	v := dec.Val*uint64(dec.Base) + uint64(d)
	if v < dec.Val {
		return ErrRange
	}
	dec.Val = v
	return nil
}

func getBase(b0, b1 byte) (base, offset int) {
	if b0 == '0' {
		if b1 == 'x' || b1 == 'X' {
			return 16, 2
		}
		return 8, 1
	}
	return 10, 0
}

func checkBase(base int) bool {
	return base >= 2 || base <= 36
}

func ParseStringUint32(s string, base int) (uint32, error) {
	if base == 0 {
		if len(s) <= 1 {
			base = 10
		} else {
			var o int
			base, o = getBase(s[0], s[1])
			s = s[o:]
		}
	} else if !checkBase(base) {
		return 0, ErrBase
	}
	if len(s) == 0 {
		return 0, ErrSyntax
	}
	dec := u32decoder{Base: uint32(base)}
	for i := 0; i < len(s); i++ {
		if err := dec.PushDigit(s[i]); err != nil {
			return 0, err
		}
	}
	return dec.Val, nil
}

func ParseStringUint64(s string, base int) (uint64, error) {
	if base == 0 {
		if len(s) <= 1 {
			base = 10
		} else {
			var o int
			base, o = getBase(s[0], s[1])
			s = s[o:]
		}
	} else if !checkBase(base) {
		return 0, ErrBase
	}
	if len(s) == 0 {
		return 0, ErrSyntax
	}
	dec := u64decoder{Base: uint(base)}
	for i := 0; i < len(s); i++ {
		if err := dec.PushDigit(s[i]); err != nil {
			return 0, err
		}
	}
	return dec.Val, nil
}

func ParseStringUint(s string, base int) (uint, error) {
	if intSize <= 4 {
		u, err := ParseStringUint32(s, base)
		return uint(u), err
	} else {
		u, err := ParseStringUint64(s, base)
		return uint(u), err
	}
}

func ParseStringInt32(s string, base int) (int32, error) {
	sign := int32(1)
	if s[0] == '-' {
		sign = -1
		s = s[1:]
	}
	u, err := ParseStringUint32(s, base)
	if err != nil {
		return 0, err
	}
	if u > 0x80000000 || sign > 0 && u == 0x80000000 {
		return 0, ErrRange
	}
	return sign * int32(u), nil
}
