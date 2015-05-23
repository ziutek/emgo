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
	spaces := [8]byte{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '}
	for {
		if width <= len(spaces) {
			p.Write(spaces[:width])
			return
		}
		p.Write(spaces[:])
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
	v := reflect.ValueOf(i)
	var (
		length int
		str    string
	)
	switch v.Kind() {
	case reflect.Bool:
		length = strconv.FormatBool(p.buf[:], v.Bool(), 2)
	case reflect.Int:
		length = strconv.FormatInt(p.buf[:], int(v.Int()), 10)
	case reflect.Int8, reflect.Int16, reflect.Int32:
		length = strconv.FormatInt32(p.buf[:], int32(v.Int()), 10)
	case reflect.Int64:
		length = strconv.FormatInt64(p.buf[:], v.Int(), 10)
	case reflect.Uint:
		length = strconv.FormatUint(p.buf[:], uint(v.Uint()), 10)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32:
		length = strconv.FormatUint32(p.buf[:], uint32(v.Uint()), 10)
	case reflect.Uint64:
		length = strconv.FormatUint64(p.buf[:], v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		prec, _ := p.Precision()
		length = strconv.FormatFloat(p.buf[:], v.Float(), 'e', prec)
	case reflect.Complex64, reflect.Complex128:
		c := v.Complex()
		p.Write([]byte{'('})
		p.format(verb, real(c))
		if imag(c) >= 0 {
			p.Write([]byte{'+'})
		}
		p.format(verb, imag(c))
		p.Write([]byte{'i', ')'})

	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.UnsafePointer:
		ptr := v.Pointer()
		if ptr == 0 {
			str = "<nil>"
			length = len(str)
			break
		}
		p.buf[0] = '0'
		p.buf[1] = 'x'
		length = 2 + strconv.FormatUint(p.buf[2:], uint(ptr), 16)

	case reflect.String:
		str = v.String()
		length = len(str)

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
