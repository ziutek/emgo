package reflect

import (
	"builtin"
)

type Kind int

const (
	Invalid Kind = iota
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	Chan
	Func
	Interface
	Map
	Ptr
	Slice
	String
	Struct
	UnsafePointer
	
	Array Kind = -1
)

var kindNames = [...]string{
	"array",
	"invalid",
	"bool",
	"int",
	"int8",
	"int16",
	"int32",
	"int64",
	"uint",
	"uint8",
	"uint16",
	"uint32",
	"uint64",
	"uintptr",
	"float32",
	"float64",
	"complex64",
	"complex128",
	"chan",
	"func",
	"interface",
	"map",
	"ptr",
	"slice",
	"string",
	"struct",
	"unsafe.Pointer",
}

func (k Kind) String() string {
	if k++; k < 0 || int(k) >= len(kindNames) {
		k = 1
	}
	return kindNames[k]
}

type Type struct {
	b *builtin.Type
}

// TypeOf returns the reflection type of value in i. TypeOf(nil) returns the
// zero Type.
func TypeOf(i interface{}) Type {
	return ValueOf(i).Type()
}

// IsValid returns true if t represents a type. It returns false if t is zero
// Type.
func (t Type) IsValid() bool {
	return t.b != nil
}

// Kind returns specific kind of t.
func (t Type) Kind() Kind {
	if !t.IsValid() {
		return Invalid
	}
	return Kind(t.b.Kind())
}

// Len returns length of array.
func (t Type) Len() int {
	if t.Kind() != Array {
		panic("reflect: not array")
	}
	return t.b.Len()
} 

// Name returns name of type within its package. It returns empty string if t
// is not valid or represents unnamed type.
func (t Type) Name() string {
	if !t.IsValid() {
		return ""
	}
	name := t.b.Name()
	for i := len(name); i != 0; i-- {
		if name[i-1] == '.' {
			return name[i:]
		}
	}
	return name
}

// String returns string representation of type.
func (t Type) String() string {
	if !t.IsValid() {
		return Invalid.String()
	}
	if name := t.b.Name(); len(name) != 0 {
		return name
	}
	return t.Kind().String()
}

func (t Type) NumElems() int {
	if !t.IsValid() {
		return 0
	}
	return len(t.b.Elems())
}

func (t Type) Elem(i int) Type {
	return Type{t.b.Elems()[i]}
}

