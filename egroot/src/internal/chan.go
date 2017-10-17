package internal

import "unsafe"

const ChanOK uintptr = 0

// TODO: Channels were implemented before interfaces. Consider use
// `type Chan interface { ... }` to implement channels.

// A Chan is internal representation of chan T type.
type Chan struct {
	C unsafe.Pointer
	M *ChanMethods
}

// A ChanMethods is set of methods used internally to implement channel
// operations.
type ChanMethods struct {
	Send       func(c, e unsafe.Pointer) (p unsafe.Pointer, d uintptr)
	Recv       func(c, e unsafe.Pointer) (p unsafe.Pointer, d uintptr)
	TrySend    func(c, e, w unsafe.Pointer) (p unsafe.Pointer, d uintptr)
	TryRecv    func(c, e, w unsafe.Pointer) (p unsafe.Pointer, d uintptr)
	CancelSend func(c, w unsafe.Pointer)
	CancelRecv func(c, w unsafe.Pointer)
	Done       func(c unsafe.Pointer, d uintptr)
	Close      func(c unsafe.Pointer)
	Len        func(c unsafe.Pointer) int
	Cap        func(c unsafe.Pointer) int
}

// MakeChan is used internally to implement make(chan T, cap) operation.
var MakeChan func(cap int, size, align uintptr) *Chan

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
// can be modified (shuffled).
var Select func(comms []*Comm, dflt unsafe.Pointer) (jmp, p unsafe.Pointer, d uintptr)

// TimeChan is used to wait for specified time (mainly for deadline/timeout in
// select statement).
var TimeChan Chan
