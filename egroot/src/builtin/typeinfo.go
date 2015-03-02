package builtin

import "unsafe"

type Method struct {
	name  string
	param []*Type `C:"const"`
}

type Type struct {
	name string
	size uintptr
	kind byte
	elem []*Type `C:"const"`
	mset []*Method
	fptr [1 << 28]unsafe.Pointer
}

func (t *Type) Kind() byte {
	return t.kind
}

func (t *Type) Size() uintptr {
	return t.size
}

func (t *Type) Name() string {
	return t.name
}

func (t *Type) Methods() []Method {
	return t.mset
}

type ItHead struct {
	*Type `C:"const"`
}

type Itable struct {
	Head ItHead
	Func [1 << 28]unsafe.Pointer
}

// GetItable should return itable for given interface and non-interface type
// pair. It is always called with etyp assignable to ityp.
var GetItable func(ityp, etyp *Type) *Itable
