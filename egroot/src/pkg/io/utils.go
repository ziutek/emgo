package io

func WriteString(w Writer, s string) (int, error) {
	return w.Write([]byte(s))
}
