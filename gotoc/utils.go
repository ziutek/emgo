package gotoc

import (
	"bytes"
	"fmt"
	"go/ast"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"

	"code.google.com/p/go.tools/go/types"
	"code.google.com/p/go.tools/go/types/typeutil"
)

func notImplemented(n ast.Node) {
	fmt.Fprintf(os.Stderr, "not implemented: %v <%T>\n", n, n)
	os.Exit(1)
}

func upath(path string) string {
	return strings.Replace(path, "/", "_", -1)
}

func tmpname(w *bytes.Buffer) string {
	return "__" + strconv.Itoa(w.Len())
}

func write(s string, ws ...io.Writer) error {
	for _, w := range ws {
		if _, err := io.WriteString(w, s); err != nil {
			return err
		}
	}
	return nil
}

type tupleNamer struct {
	th    typeutil.Hasher
	mutex sync.Mutex
	known map[uint32]struct{}
}

func newTupleNamer() *tupleNamer {
	return &tupleNamer{th: typeutil.MakeHasher(), known: make(map[uint32]struct{})}
}

func (tn *tupleNamer) name(h uint32) string {
	return "__tuple" + strconv.FormatUint(uint64(h), 10)
}

func (tn *tupleNamer) DeclName(t *types.Tuple) (string, bool) {
	tn.mutex.Lock()
	h := tn.th.Hash(t)
	_, known := tn.known[h]
	if !known {
		tn.known[h] = struct{}{}
	}
	tn.mutex.Unlock()
	return tn.name(h), known
}

func (tn *tupleNamer) Name(t *types.Tuple) string {
	tn.mutex.Lock()
	h := tn.th.Hash(t)
	tn.mutex.Unlock()
	return tn.name(h)
}
