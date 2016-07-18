package gotoc

import (
	"io"
	"strings"
)

func Upath(s string) (ret string) {
	for {
		i := strings.IndexAny(s, "/.+-")
		if i == -1 {
			break
		}
		ret += s[:i]
		switch s[i] {
		case '/':
			ret += "$"
		case '.':
			ret += "$0$"
		case '-':
			ret += "$1$"
		case '+':
			ret += "$2$"
		}
		s = s[i+1:]
	}
	ret += s
	return
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
