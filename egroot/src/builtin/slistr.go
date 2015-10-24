package builtin

type Slice struct {
	Addr, Len, Cap uintptr
}

type String struct {
	Addr, Len uintptr
}
