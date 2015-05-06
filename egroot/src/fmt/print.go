package fmt

import "io"

func Fprint(w io.Writer, a ...interface{}) (int, error) {
	p := printer{W: w}
	p.Parse("")
	var n int
	for _, v := range a {
		n += p.format(v)
		if p.Err != nil {
			break
		}
	}
	return n, p.Err
}

func Fprintln(w io.Writer, a ...interface{}) (int, error) {
	p := printer{W: w}
	p.Parse("")
	var n int
	for i, v := range a {
		if i > 0 {
			n += p.write([]byte{' '})
			if p.Err != nil {
				return n, p.Err
			}
		}
		n += p.format(v)
		if p.Err != nil {
			return n, p.Err
		}
	}
	n += p.write([]byte{'\n'})
	return n, p.Err
}

var DefaultWriter io.Writer

func Print(a ...interface{}) (int, error) {
	return Fprint(DefaultWriter, a...)
}

func Println(a ...interface{}) (int, error) {
	return Fprintln(DefaultWriter, a...)
}
