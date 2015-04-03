package fmt

import (
	"io"
	"reflect"
	"strconv"
)

type printer struct {
	io.Writer
	width int
	prec  int
	flags string
}

func (p *printer) parse(verb string) (needWidth, needPrec bool) {
	if len(verb) == 0 {
		p.width = -1
		p.prec = -1
		p.flags = ""
		return
	}
	return
}

func (p *printer) WriteString(s string) (int, error) {
	// io.WriteString uses WriteString method if p.Writer imlements it.
	return io.WriteString(p.Writer, s)
}

func (p *printer) Width() (int, bool) {
	return p.width, p.width != -1
}

func (p *printer) Precision() (int, bool) {
	return p.prec, p.prec != -1
}

func (p *printer) Flag(c int) bool {
	for i := 0; i < len(p.flags); i++ {
		if int(p.flags[i]) == c {
			return true
		}
	}
	return false
}

func (p *printer) format(i interface{}) (int, error) {
	var (
		buf [65]byte
		n   int
	)
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Invalid:
		return p.WriteString("<nil>")
	case reflect.String:
		return p.WriteString(v.String())
	case reflect.Int:
		n = strconv.FormatInt(buf[:], int(v.Int()), 10)
	case reflect.Int8, reflect.Int16, reflect.Int32:
		n = strconv.FormatInt32(buf[:], int32(v.Int()), 10)
	case reflect.Int64:
		n = strconv.FormatInt64(buf[:], int64(v.Int()), 10)
	case reflect.Uint:
		n = strconv.FormatUint(buf[:], uint(v.Uint()), 10)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32:
		n = strconv.FormatUint32(buf[:], uint32(v.Uint()), 10)
	case reflect.Uint64:
		n = strconv.FormatUint64(buf[:], uint64(v.Uint()), 10)
	}
	return p.Write(buf[n:])
}

func Fprint(w io.Writer, a ...interface{}) (int, error) {
	var n int
	p := printer{Writer: w}
	p.parse("")
	for _, v := range a {
		m, err := p.format(v)
		n += m
		if err != nil {
			return n, err
		}
	}
	return n, nil
}
