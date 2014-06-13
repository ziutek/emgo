package builtin

import "unsafe"

const (
	ChanOK uintptr = iota
	ChanClosed
	ChanAgain
)

// A Chan is internal representation of chan T type.
type Chan struct {
	C unsafe.Pointer
	M *ChanMethods
}

// A ChanMethods is set of methods used internally to perform send, receive and
// close operations.
type ChanMethods struct {
	Send       func(c, e unsafe.Pointer) (p unsafe.Pointer, d uintptr)
	Recv       func(c, e unsafe.Pointer) (p unsafe.Pointer, d uintptr)
	TrySend    func(c, e, w unsafe.Pointer) (p unsafe.Pointer, d uintptr)
	TryRecv    func(c, e, w unsafe.Pointer) (p unsafe.Pointer, d uintptr)
	CancelSend func(c, w unsafe.Pointer)
	CancelRecv func(c, w unsafe.Pointer)
	Done       func(c unsafe.Pointer, d uintptr)
	Close      func(c unsafe.Pointer)
}

// MakeChan is used internally to implement make(chan T, cap) operation.
var MakeChan func(cap int, size, align uintptr) Chan

// A Comm represents send statement or receive expression. It is used
// internally to implement select statement.
type Comm struct {
	Case   unsafe.Pointer
	C      unsafe.Pointer
	E      unsafe.Pointer
	Try    func(c, e, w unsafe.Pointer) (p unsafe.Pointer, d uintptr)
	Cancel func(c, w unsafe.Pointer)
}

// Select is used internelly to implement select statement. After Select comms
// can be modified (shufled).
var Select func(comms []*Comm, dflt unsafe.Pointer) (jmp, p unsafe.Pointer, d uintptr)
