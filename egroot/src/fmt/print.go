package fmt

import "io"

func Fprint(w io.Writer, a ...interface{}) (int, error) {
	p := printer{W: w}
	p.parse("")
	for _, v := range a {
		p.format('v', v)
		if p.Err != nil {
			break
		}
	}
	return p.N, p.Err
}

func Fprintln(w io.Writer, a ...interface{}) (int, error) {
	p := printer{W: w}
	p.parse("")
	for i, v := range a {
		if i > 0 {
			p.Write([]byte{' '})
			if p.Err != nil {
				return p.N, p.Err
			}
		}
		p.format('v', v)
		if p.Err != nil {
			return p.N, p.Err
		}
	}
	p.Write([]byte{'\n'})
	return p.N, p.Err
}

var DefaultWriter io.Writer

func Print(a ...interface{}) (int, error) {
	return Fprint(DefaultWriter, a...)
}

func Println(a ...interface{}) (int, error) {
	return Fprintln(DefaultWriter, a...)
}
