package builtin

type Type struct {
	size uintptr
	name string
	kind byte
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
	ptr   bool
}

func (ith *ITHead) Ptr() bool {
	return ith.ptr
}
