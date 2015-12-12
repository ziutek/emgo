package io

type stringWriter interface {
	WriteString(string) (int, error)
}

func WriteString(w Writer, s string) (int, error) {
	if sw, ok := w.(stringWriter); ok {
		return sw.WriteString(s)
	}
	return w.Write([]byte(s))
}

/*
type Drop struct{}

func (d Drop) Write(b []byte) (int, error)       { return len(b), nil }
func (d Drop) WriteString(s string) (int, error) { return len(s), nil }
*/
