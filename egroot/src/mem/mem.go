// Package mem provides more controll to memory allocation (not implemented!).
package mem

import (
	"internal"
)

// Type is a bitfield that describes limitations of some memory type. Zero
// value means no limitations.
type Type int

const (
	// TypeNoDMA bit means memory that can't be accessed by DMA.
	TypeNoDMA Type = internal.MemNoDMA

	// Stack is special type, only allowed to be returned in case of len or cap
	// parameter of make function.
	TypeStack Type = internal.MemStack
)

// Typer is an interface that can be implemented by any dynamically allocated
// object. It is used by new and make function to determine type of memory that
// satisfies object's needs.
//
// Whether type implements Typer is determined at compilation time so there is
// no any additional runtime overhead for "ordinary" allocation and small
// (function call, maybe inlined) overhead for allocation of memory of specific
// type.
//
// new and make functions requires MemoryType method that accepts a pointer
// receiver. It is called just before allocation with nil value as receiver.
//
// In case of make, Typer can be implemented by object type parameter and only
// one of len and cap parameters. If len or cap implements Typer only their
// MemoryType method is called so it can be used to overwrite type of memory
// required be allocated object.
type Typer interface {
	MemoryType() Type
}

// NoDMA can be used to cast len or cap parameter to permit make to allocate
// object in memory region that doesn't support DMA.
type NoDMA int

func (_ *NoDMA) MemoryType() Type {
	return TypeNoDMA
}

// Stack can be used to cast len or cap parameter to force make to allocate
// object on current stack. Allocated memory is always fried at function return.
type Stack int

func (_ *Stack) MemoryType() Type {
	return TypeStack
}
