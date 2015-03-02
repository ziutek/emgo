package builtin

import "unsafe"

type Method struct {
	name  string
	param []*Type `C:"const"`
	fptr  unsafe.Pointer
}

type Type struct {
	name string
	size uintptr
	kind byte
	elem []*Type `C:"const"`
	mset []Method
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

type ItHead struct {
	*Type `C:"const"`
	// ItHead size must be n * sizeof(uintptr)
}
