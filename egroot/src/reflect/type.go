package reflect

import (
	"builtin"
)

type Kind byte

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
	Array
	Chan
	Func
	Interface
	Map
	Ptr
	Slice
	String
	Struct
	UnsafePointer
)

var kindNames = [...]string{
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
	"array",
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
	if k < 0 || int(k) >= len(kindNames) {
		k = 0
	}
	return kindNames[k]
}

type Type struct {
	b *builtin.Type
}

func TypeOf(i interface{}) Type {
	return ValueOf(i).Type()
}

func (t Type) Kind() Kind {
	return Kind(t.b.Kind())
}

func (t Type) String() string {
	name := t.b.Name()
	if len(name) == 0 {
		name = t.Kind().String()
	}
	return name
}
