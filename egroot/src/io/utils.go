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
