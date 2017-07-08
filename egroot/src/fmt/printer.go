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

func (wpf *wpf) parse(flags string) {
	if len(flags) == 0 {
		wpf.width = -2
		wpf.prec = -2
		wpf.flags = ""
		return
	}
	i := 0
	for ; i < len(flags); i++ {
		c := flags[i]
		if c >= '1' && c <= '9' || c == '.' {
			break
		}
	}
	wpf.flags = flags[:i]
	if i == len(flags) {
		return
	}
	flags = flags[i:]
	for i = 0; i < len(flags); i++ {
		if flags[i] == '.' {
			break
		}
	}
	if i > 0 {
		width, _ := strconv.ParseStringUint32(flags[:i], 10)
		wpf.width = int16(width)
	}
	if i >= len(flags)-1 {
		return
	}
	prec, _ := strconv.ParseStringUint32(flags[i+1:], 10)
	wpf.prec = int16(prec)
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
// Value of type printer is small enough to be be assigned to interface type.
type printer struct {
	*writer
	wpf
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
	if verb != 'T' && verb != 'p' {
		if f, ok := i.(Formatter); ok {
			f.Format(p, rune(verb))
			return
		}
	}
	switch verb {
	case 'T':
		i = reflect.TypeOf(i).String()
	case 'v', 's', 'q': // Do not follow original rule in case of x and X.
		switch f := i.(type) {
		case error:
			i = f.Error()
		case Stringer:
			i = f.String()
		}
	}
	if i == nil {
		i = "<nil>"
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
	width, _ := p.Width()
	if p.Flag('-') {
		width = -width
	}
	zeros := p.Flag('0')
	k := v.Kind()
	switch {
	case k == reflect.Array || k == reflect.Slice:
		if verb == 's' && k == reflect.Slice {
			if b, ok := v.Interface().([]byte); ok {
				strconv.WriteBytes(p, b, width, zeros)
				break
			}
		}
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
		p.WriteString("<invalid>")
	case k == reflect.Bool:
		strconv.WriteBool(p, v.Bool(), 't', width)
	case k <= reflect.Uintptr:
		p.formatIntVal(v, width, verb, zeros)
	case k <= reflect.Float64:
		p.formatFloatVal(v, width, verb, zeros)
	case k <= reflect.Complex128:
		c := v.Complex()
		p.WriteByte('(')
		p.formatFloatVal(reflect.ValueOf(real(c)), width, verb, zeros)
		if imag(c) >= 0 {
			p.WriteByte('+')
		}
		p.formatFloatVal(reflect.ValueOf(imag(c)), width, verb, zeros)
		p.WriteString("i)")
	case k <= reflect.Func || k == reflect.Ptr || k == reflect.UnsafePointer:
		ptr := v.Pointer()
		if verb != 'v' {
			p.formatIntVal(reflect.ValueOf(ptr), width, verb, zeros)
			break
		}
		if ptr == 0 {
			p.WriteString("<nil>")
		} else {
			p.WriteString("0x")
			strconv.WriteUintptr(p, ptr, 16, 0)
		}
	case k == reflect.String:
		strconv.WriteString(p, v.String(), width, zeros)
	case k == reflect.Struct:
		p.WriteByte('{')
		p.WriteString("TODO")
		p.WriteByte('}')
	default:
		p.WriteString("<!not supported>")
	}
}

func (p *printer) formatIntVal(v reflect.Value, width int, verb byte, zeros bool) {
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
			p.WriteByte(byte(v.Int()))
		} else {
			p.WriteByte(byte(v.Uint()))
		}
		return
	default:
		p.badVerb(verb, v)
	}
	if zeros {
		base = -base
	}
	switch {
	case k == reflect.Int:
		strconv.WriteInt(p, int(v.Int()), base, width)
	case k <= reflect.Int32:
		strconv.WriteInt32(p, int32(v.Int()), base, width)
	case k == reflect.Int64:
		strconv.WriteInt64(p, v.Int(), base, width)
	case k == reflect.Uint:
		strconv.WriteUint(p, uint(v.Uint()), base, width)
	case k <= reflect.Uint32:
		strconv.WriteUint32(p, uint32(v.Uint()), base, width)
	case k == reflect.Uint64:
		strconv.WriteUint64(p, v.Uint(), base, width)
	default:
		strconv.WriteUintptr(p, uintptr(v.Uint()), base, width)
	}
}

func (p *printer) formatFloatVal(v reflect.Value, width int, verb byte, zeros bool) {
	bits := 32
	if v.Kind() == reflect.Float64 {
		bits = 64
	}
	prec, _ := p.Precision()
	fmt := int(verb)
	if fmt == 'v' {
		fmt = 'g'
	}
	if zeros {
		fmt = -fmt
	}
	strconv.WriteFloat(p, v.Float(), fmt, width, prec, bits)
}
