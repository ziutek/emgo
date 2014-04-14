package gotoc

import (
	"fmt"
	"go/ast"
	"io"
	"os"
	"strings"
)

func notImplemented(n ast.Node) {
	fmt.Fprintf(os.Stderr, "not implemented: %T\n", n)
	os.Exit(1)
}

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