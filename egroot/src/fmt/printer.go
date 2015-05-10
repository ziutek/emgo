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

	width int
	prec  int
	flags string
	buf   [65]byte
}

func (p *printer) Parse(verb string) {
	if len(verb) == 0 {
		p.width = -2
		p.prec = -2
		p.flags = ""
		return
	}
	return
}

func (p *printer) write(b []byte) (n int) {
	if p.Err != nil {
		return
	}
	n, p.Err = p.W.Write(b)
	return
}

func (p *printer) Write(b []byte) (int, error) {
	return p.write(b), p.Err
}

func (p *printer) writeString(s string) (n int) {
	if p.Err != nil {
		return
	}
	// io.WriteString uses WriteString method if p.W imlements it.
	n, p.Err = io.WriteString(p.W, s)
	return
}

func (p *printer) WriteString(s string) (int, error) {
	return p.writeString(s), p.Err
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

func (p *printer) padSpaces(length int) int {
	width, wok := p.Width()
	if !wok {
		return 0
	}
	width -= length
	if width <= 0 {
		return 0
	}
	spaces := [8]byte{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '}
	m := width
	for {
		if m > len(spaces) {
			m -= p.write(spaces[:])
		} else {
			m -= p.write(spaces[:m])
			break
		}
		if p.Err != nil {
			break
		}
	}
	return width - m
}

func (p *printer) format(i interface{}) (n int) {
	switch f := i.(type) {
	case nil:
		i = "<nil>"
	case error:
		i = f.Error()
	case Formatter:

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
		n += p.write([]byte{'('})
		n += p.format(real(c))
		if imag(c) >= 0 {
			n += p.write([]byte{'+'})
		}
		n += p.format(imag(c))
		n += p.write([]byte{'i', ')'})

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
		n = p.padSpaces(length)
	}
	if length != 0 {
		if len(str) != 0 {
			n += p.writeString(str)
		} else {
			n += p.write(p.buf[:length])
		}
	}
	if left {
		n += p.padSpaces(length)
	}
	return
}
