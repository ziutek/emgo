package fmt

import "io"

func Fprint(w io.Writer, a ...Formatter) (n int, err error) {
	var m int
	for _, v := range a {
		m, err = v.Format(w)
		n += m
		if err != nil {
			return
		}
	}
	return
}

func findVerb(f string) (start int, verb byte, flags string) {
	for f[start] != '%' {
		start++
		if start >= len(f) {
			return
		}
	}
	end := start + 1
	for ; end < len(f); end++ {
		c := f[end]
		if c >= '0' && c <= '9' {
			continue
		}
		switch c {
		case '+', '-', ' ', '*', '#':
			continue
		}
		break
	}
	flags = f[start+1 : end]
	if end != len(f) {
		verb = f[end]
	}
	return
}

func ferr(w io.Writer, verb byte, info string, a Formatter) (int, error) {
	if a == nil {
		a = Bytes(nil)
	}
	if verb == 0 {
		return Fprint(w, Str("%!("), Str(info), a, Rune(')'))
	}
	return Fprint(w, Str("%!"), Bytes{verb, '('}, Str(info), a, Rune(')'))
}

func Fprintf(w io.Writer, f string, a ...Formatter) (int, error) {
	var i, n int
	for {
		first, verb, flags := findVerb(f)
		m, err := io.WriteString(w, f[:first])
		n += m
		if err != nil {
			return n, err
		}
		if first == len(f) {
			break
		}
		switch verb {
		case 'v':
			if i < len(a) {
				m, err = a[i].Format(w)
			}
			i++
		case '%':
			m, err = w.Write([]byte{'%'})
		case 0:
			// Unfinished format.
			m, err = ferr(w, 0, "UNFINISHED", nil)
		default:
			// Unkonown format
			m, err = ferr(w, verb, "UNKNOWN", nil)
		}
		n += m
		if err != nil {
			return n, err
		}
		if i > len(a) {
			m, err = ferr(w, verb, "MISSING", nil)
			n += m
			if err != nil {
				return n, err
			}
		}
		f = f[first+2+len(flags):]
	}
	for ; i < len(a); i++ {
		m, err := ferr(w, 0, "EXTRA ", a[i])
		n += m
		if err != nil {
			return n, err
		}
	}
	return n, nil
}
