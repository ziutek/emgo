package builtin

type Type struct {
	size uintptr
	name string
}

type ITHead struct {
	typ *Type `C:"const"`
	ptr bool
}

func (ith *ITHead) Size() uintptr {
	return ith.typ.size
}

func (ith *ITHead) Name() string {
	return ith.typ.name
}

func (ith *ITHead) Ptr() bool {
	return ith.ptr
}
