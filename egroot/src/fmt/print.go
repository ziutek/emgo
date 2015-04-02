package fmt

import (
	"io"
	"reflect"
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

func (p *printer) format(v interface{}) (int, error) {

}


func Fprint(w io.Writer, a ...interface{}) (n int, err error) {
	p := printer{Writer: w}
	p.parse("")
	for _, v := range a {
		m, err := p.format(v)
		n += m
		if err != nil {
			return
		}
	}
	return
}
