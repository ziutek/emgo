package builtin

import "unsafe"

type Chan struct {
	P unsafe.Pointer
	M *ChanMethods
}

type ChanMethods struct {
	Send    func(c, e unsafe.Pointer) (p unsafe.Pointer, d uintptr)
	Recv    func(c, e unsafe.Pointer) (p unsafe.Pointer, d uintptr)
	TrySend func(c, e unsafe.Pointer) (p unsafe.Pointer, d uintptr)
	TryRecv func(c, e unsafe.Pointer) (p unsafe.Pointer, d uintptr)
	Done    func(c unsafe.Pointer, d uintptr)
	Close   func(c unsafe.Pointer)
}

var MakeChan func(cap int, size, align uintptr) Chan
