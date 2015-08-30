package fmt

import (
	"io"
	"reflect"
	"strconv"
	"strings"
)

type wpf struct {
	width int16
	prec  int16
	flags string
}

func (wpf *wpf) parse(format string) {
	if len(format) == 0 {
		wpf.width = -2
		wpf.prec = -2
		wpf.flags = ""
		return
	}
	return
}

func (wpf *wpf) Width() (int, bool) {
	set := (wpf.width != -2)
	if set {
		return int(wpf.width), set
	}
	return 0, set
}

func (wpf *wpf) Precision() (int, bool) {
	set := (wpf.prec != -2)
	if set {
		return int(wpf.prec), set
	}
	return 6, set
}

func (wpf *wpf) Flag(c int) bool {
	return strings.IndexByte(wpf.flags, byte(c)) != -1
}

type writer struct {
	w   io.Writer
	err error
	n   int

	buf [65]byte
}

func (w *writer) Write(b []byte) (int, error) {
	if w.err != nil {
		return 0, w.err
	}
	var n int
	n, w.err = w.w.Write(b)
	w.n += n
	return n, w.err
}

func (w *writer) WriteString(s string) (int, error) {
	if w.err != nil {
		return 0, w.err
	}
	var n int
	// io.WriteString uses WriteString method if w.w implements it.
	n, w.err = io.WriteString(w.w, s)
	w.n += n
	return n, w.err
}

func (w *writer) WriteByte(b byte) error {
	w.Write([]byte{b})
	return w.err
}

// printer implements State interface.
// Value of type printer can be assigned to interface type.
type printer struct {
	*writer
	wpf
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
		if p.err != nil {
			return
		}
		width -= len(spaces)
	}
}

func (p *printer) Ferr(verb byte, info string, a interface{}) {
	if a == nil {
		a = ""
	}
	p.WriteString("%!")
	if verb != 0 {
		p.WriteByte(verb)
	}
	p.WriteByte('(')
	p.WriteString(info)
	p.parse("")
	p.format('v', a)
	p.WriteByte(')')
}

func (p *printer) badVerb(verb byte, v reflect.Value) {
	p.WriteString("%!")
	p.WriteByte(verb)
	p.WriteByte('(')
	p.WriteString(v.Type().String())
	p.WriteByte('=')
	p.formatValue('v', v)
	p.WriteByte(')')
}

func (p *printer) format(verb byte, i interface{}) {
	switch verb {
	case 'T':
		i = reflect.TypeOf(i).String()
	default:
		switch f := i.(type) {
		case nil:
			i = "<nil>"
		case Formatter:
			f.Format(p, rune(verb))
			return
		case error:
			i = f.Error()
		case Stringer:
			i = f.String()
		}
	}
	p.formatValue(verb, reflect.ValueOf(i))
}

func (p *printer) tryFormatAsInterface(verb byte, v reflect.Value) {
	if i := v.Interface(); i != nil || !v.IsValid() {
		p.format(verb, i)
	} else {
		p.formatValue(verb, v)
	}
}

func (p *printer) formatValue(verb byte, v reflect.Value) {
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
			p.tryFormatAsInterface(verb, v.Index(i))
		}
		p.WriteByte(']')
	case k == reflect.Invalid:
		str = "<invalid>"
		length = len(str)
	case k == reflect.Bool:
		length = strconv.FormatBool(p.buf[:], v.Bool(), 2)
	case k <= reflect.Uintptr:
		length = p.formatIntVal(verb, v)
	case k <= reflect.Float64:
		length = p.formatFloat(verb, v)
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

func (p *printer) formatIntVal(verb byte, v reflect.Value) int {
	k := v.Kind()
	base := 10
	switch verb {
	case 'v', 'd':
		base = 10
	case 'x', 'X':
		base = 16
	case 'b':
		base = 2
	case 'o':
		base = 8
	case 'c':
		if k <= reflect.Int64 {
			p.buf[0] = byte(v.Int())
		} else {
			p.buf[0] = byte(v.Uint())
		}
		return 1
	default:
		p.badVerb(verb, v)
		return 0
	}
	switch {
	case k == reflect.Int:
		return strconv.FormatInt(p.buf[:], int(v.Int()), base)
	case k <= reflect.Int32:
		return strconv.FormatInt32(p.buf[:], int32(v.Int()), base)
	case k == reflect.Int64:
		return strconv.FormatInt64(p.buf[:], v.Int(), base)
	case k == reflect.Uint:
		return strconv.FormatUint(p.buf[:], uint(v.Uint()), base)
	case k <= reflect.Uint32:
		return strconv.FormatUint32(p.buf[:], uint32(v.Uint()), base)
	case k == reflect.Uint64:
		return strconv.FormatUint64(p.buf[:], v.Uint(), base)
	default:
		return strconv.FormatUintptr(p.buf[:], uintptr(v.Uint()), base)
	}
}

func (p *printer) formatFloat(verb byte, v reflect.Value) int {
	bitsize := 32
	if v.Kind() == reflect.Float64 {
		bitsize = 64
	}
	switch verb {
	case 'v':
		verb = 'e'
	case 'e', 'E', 'b':
	default:
		p.badVerb(verb, v)
		return 0
	}
	prec, _ := p.Precision()
	return strconv.FormatFloat(p.buf[:], v.Float(), int(verb), prec, bitsize)
}
