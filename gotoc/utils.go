package gotoc

import (
	"fmt"
	"go/ast"
	"io"
	"os"
	"strings"

	"code.google.com/p/go.tools/go/types"
)

func notImplemented(n ast.Node, tl ...types.Type) {
	fmt.Fprintf(os.Stderr, "not implemented: %T\n", n)
	for _, t := range tl {
		fmt.Fprintf(os.Stderr, "	in case of: %T\n", t)
	}
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
