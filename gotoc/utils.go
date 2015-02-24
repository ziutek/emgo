package gotoc

import (
	"io"
	"strings"
)

func upath(path string) string {
	return strings.Replace(path, "/", "$", -1)
}

func write(s string, ws ...io.Writer) error {
	for _, w := range ws {
		if _, err := io.WriteString(w, s); err != nil {
			return err
		}
	}
	return nil
}

func indent(n int, s string) string {
	nt := "\n" + strings.Repeat("\t", n)
	return strings.Replace(s, "\n", nt, -1)
}