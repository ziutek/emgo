// Package stack allows allocate memory on stack using
// C alloca function.

func alignUp(p, a uintptr) uintptr {
	a--
	return (p + a) &^ a
}

func alloca() unsafe.Pointer

func Alloc(n int, size, align uintptr) unsafe.Pointer {
	size = alignUp(size, align) * uintptr(n)
	
}