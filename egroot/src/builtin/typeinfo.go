package builtin

type Type struct {
	name string
	size uintptr
	kind byte
	elem []*Type `C:"const"`
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

type ITHead struct {
	*Type `C:"const"`
}
