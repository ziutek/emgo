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

type discard int

func (_ discard) Write(b []byte) (int, error)       { return len(b), nil }
func (_ discard) WriteString(s string) (int, error) { return len(s), nil }

// Discard implements Writer and StringWriter and simply
// discards all data written to it.
const Discard discard = 0
