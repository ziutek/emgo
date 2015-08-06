package fmt

import (
	"io"
	"reflect"
	"strconv"
)

// prineter implements State interface
type printer struct {
	W   io.Writer
	Err error
	N   int

	width int
	prec  int
	flags string
	buf   [65]byte
}

func (p *printer) parse(verb string) {
	if len(verb) == 0 {
		p.width = -2
		p.prec = -2
		p.flags = ""
		return
	}
	return
}

func (p *printer) Write(b []byte) (int, error) {
	if p.Err != nil {
		return 0, p.Err
	}
	var n int
	n, p.Err = p.W.Write(b)
	p.N += n
	return n, p.Err
}

func (p *printer) WriteString(s string) (int, error) {
	if p.Err != nil {
		return 0, p.Err
	}
	var n int
	// io.WriteString uses WriteString method if p.W imlements it.
	n, p.Err = io.WriteString(p.W, s)
	p.N += n
	return n, p.Err
}

func (p *printer) WriteByte(b byte) error {
	p.Write([]byte{b})
	return p.Err
}

func (p *printer) Width() (int, bool) {
	set := (p.width != -2)
	if set {
		return p.width, set
	}
	return 0, set
}

func (p *printer) Precision() (int, bool) {
	set := (p.prec != -2)
	if set {
		return p.prec, set
	}
	return 6, set
}

func (p *printer) Flag(c int) bool {
	for i := 0; i < len(p.flags); i++ {
		if int(p.flags[i]) == c {
			return true
		}
	}
	return false
}

func (p *printer) padSpaces(length int) {
	width, wok := p.Width()
	if !wok {
		return
	}
	width -= length
	if width <= 0 {
		return
	}
	spaces := "        "
	for {
		if width <= len(spaces) {
			p.WriteString(spaces[:width])
			return
		}
		p.WriteString(spaces)
		if p.Err != nil {
			return
		}
		width -= len(spaces)
	}
}

func (p *printer) format(verb rune, i interface{}) {
	switch f := i.(type) {
	case nil:
		i = "<nil>"
	case Formatter:
		f.Format(p, verb)
		return
	case error:
		i = f.Error()
	case Stringer:
		i = f.String()
	}
	p.formatValue(verb, reflect.ValueOf(i))
}

func (p *printer) formatTryInterface(verb rune, v reflect.Value) {
	if i := v.Interface(); i != nil || !v.IsValid() {
		p.format(verb, i)
	} else {
		p.formatValue(verb, v)
	}
}

func (p *printer) formatValue(verb rune, v reflect.Value) {
	var (
		length int
		str    string
	)
	k := v.Kind()
	switch {
	case k == reflect.Array || k == reflect.Slice:
		p.WriteByte('[')
		n := v.Len()
		for i := 0; i < n; i++ {
			if i > 0 {
				p.WriteByte(' ')
			}
			p.formatTryInterface(verb, v.Index(i))
		}
		p.WriteByte(']')
	case k == reflect.Invalid:
		str = "<invalid>"
		length = len(str)
	case k == reflect.Bool:
		length = strconv.FormatBool(p.buf[:], v.Bool(), 2)
	case k == reflect.Int:
		length = strconv.FormatInt(p.buf[:], int(v.Int()), 10)
	case k <= reflect.Int32:
		length = strconv.FormatInt32(p.buf[:], int32(v.Int()), 10)
	case k == reflect.Int64:
		length = strconv.FormatInt64(p.buf[:], v.Int(), 10)
	case k == reflect.Uint:
		length = strconv.FormatUint(p.buf[:], uint(v.Uint()), 10)
	case k <= reflect.Uint32:
		length = strconv.FormatUint32(p.buf[:], uint32(v.Uint()), 10)
	case k == reflect.Uint64:
		length = strconv.FormatUint64(p.buf[:], v.Uint(), 10)
	case k == reflect.Uintptr:
		length = strconv.FormatUintptr(p.buf[:], uintptr(v.Uint()), 10)
	case k <= reflect.Float64:
		prec, _ := p.Precision()
		length = strconv.FormatFloat(p.buf[:], v.Float(), 'e', prec)
	case k <= reflect.Complex128:
		c := v.Complex()
		p.WriteByte('(')
		p.format(verb, real(c))
		if imag(c) >= 0 {
			p.WriteByte('+')
		}
		p.format(verb, imag(c))
		p.WriteString("i)")
	case k == reflect.Ptr:
		if v.IsNil() {
			str = "<nil>"
			length = len(str)
			break
		}
		p.WriteByte('&')
		p.formatValue(verb, v.Elem())
	case k <= reflect.Func || k == reflect.UnsafePointer:
		ptr := v.Pointer()
		if ptr == 0 {
			str = "<nil>"
			length = len(str)
			break
		}
		p.buf[0] = '0'
		p.buf[1] = 'x'
		length = 2 + strconv.FormatUintptr(p.buf[2:], ptr, 16)
	case k == reflect.String:
		str = v.String()
		length = len(str)
	case k == reflect.Struct:
		p.WriteByte('{')

		p.WriteByte('}')
	default:
		str = "<!not supported>"
		length = len(str)
	}
	left := p.Flag('-')
	if !left {
		p.padSpaces(length)
	}
	if length != 0 {
		if len(str) != 0 {
			p.WriteString(str)
		} else {
			p.Write(p.buf[:length])
		}
	}
	if left {
		p.padSpaces(length)
	}
	return
}
