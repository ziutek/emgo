package reflect

import (
	"internal"
	"unsafe"
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

type kindInfo struct {
	name  string
	size  byte
	align byte
}

//emgo:const
var kinfos = [...]kindInfo{
	{"array", 0, 0},
	{"invalid", 0, 0},
	{"bool", byte(unsafe.Sizeof(bool(false))), byte(unsafe.Alignof(bool(false)))},
	{"int", byte(unsafe.Sizeof(int(0))), byte(unsafe.Alignof(int(0)))},
	{"int8", byte(unsafe.Sizeof(int8(0))), byte(unsafe.Alignof(int8(0)))},
	{"int16", byte(unsafe.Sizeof(int16(0))), byte(unsafe.Alignof(int16(0)))},
	{"int32", byte(unsafe.Sizeof(int32(0))), byte(unsafe.Alignof(int32(0)))},
	{"int64", byte(unsafe.Sizeof(int64(0))), byte(unsafe.Alignof(int64(0)))},
	{"uint", byte(unsafe.Sizeof(uint(0))), byte(unsafe.Alignof(uint(0)))},
	{"uint8", byte(unsafe.Sizeof(uint8(0))), byte(unsafe.Alignof(uint8(0)))},
	{"uint16", byte(unsafe.Sizeof(uint16(0))), byte(unsafe.Alignof(uint16(0)))},
	{"uint32", byte(unsafe.Sizeof(uint32(0))), byte(unsafe.Alignof(uint32(0)))},
	{"uint64", byte(unsafe.Sizeof(uint64(0))), byte(unsafe.Alignof(uint64(0)))},
	{"uintptr", byte(unsafe.Sizeof(uintptr(0))), byte(unsafe.Alignof(uintptr(0)))},
	{"float32", byte(unsafe.Sizeof(float32(0))), byte(unsafe.Alignof(float32(0)))},
	{"float64", byte(unsafe.Sizeof(float64(0))), byte(unsafe.Alignof(float64(0)))},
	{"complex64", byte(unsafe.Sizeof(complex64(0))), byte(unsafe.Alignof(complex64(0)))},
	{"complex128", byte(unsafe.Sizeof(complex128(0))), byte(unsafe.Alignof(complex128(0)))},
	{"chan", byte(unsafe.Sizeof((chan int)(nil))), byte(unsafe.Alignof((chan int)(nil)))},
	{"func", byte(unsafe.Sizeof(func() {})), byte(unsafe.Alignof(func() {}))},
	{"interface", byte(unsafe.Sizeof(interface{}(nil))), byte(unsafe.Alignof(interface{}(nil)))},
	{"map", byte(unsafe.Sizeof(map[int]int(nil))), byte(unsafe.Alignof(map[int]int(nil)))},
	{"ptr", byte(unsafe.Sizeof(uintptr(0))), byte(unsafe.Alignof(uintptr(0)))},
	{"slice", byte(unsafe.Sizeof([]byte(nil))), byte(unsafe.Alignof([]byte(nil)))},
	{"string", byte(unsafe.Sizeof("")), byte(unsafe.Alignof(""))},
	{"struct", 0, 0},
	{"unsafe.Pointer", byte(unsafe.Sizeof(uintptr(0))), byte(unsafe.Alignof(uintptr(0)))},
}

func (k Kind) String() string {
	if k++; uint(k) >= uint(len(kinfos)) {
		k = 1
	}
	return kinfos[k].name
}

type Type struct {
	b *internal.Type
}

const invalidT = "reflect: invalid type"

// Size returns the number of bytes that value of type t needs in memory.
func (t Type) Size() uintptr {
	k := t.Kind()
	switch k {
	case Array:
		return t.Elem().Size() * uintptr(t.Len())
	case Struct:
		var size uintptr
		n := t.NumField()
		for i := 0; i < n; i++ {
			size += t.Field(i).Size()
		}
		return size
	}
	if k++; uint(k) >= uint(len(kinfos)) {
		k = 1
	}
	return uintptr(kinfos[k].size)
}

func (t Type) Align() uintptr {
	k := t.Kind()
	switch k {
	case Array:
		return t.Elem().Align()
	case Struct:
		var align uintptr
		n := t.NumField()
		for i := 0; i < n; i++ {
			if a := t.Field(i).Align(); a > align {
				align = a
			}
		}
		return align
	}
	if k++; uint(k) >= uint(len(kinfos)) {
		k = 1
	}
	return uintptr(kinfos[k].align)
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

const badKind = "reflect: bad kind"

// Len returns length of array. It panics if kind of t isn't Array.
func (t Type) Len() int {
	if t.Kind() != Array {
		panic(badKind)
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

// Elem return type of element of t. It panics if kind of t isn't Array, Chan,
// Ptr, or Slice.
func (t Type) Elem() Type {
	if e := t.b.Elems(); len(e) == 1 && t.Kind() != Struct {
		return Type{e[0]}
	}
	panic(badKind)
}

// NumField returns number of fields in struct. It panics if kind of t isn't
// Struct.
func (t Type) NumField() int {
	if t.Kind() != Struct {
		panic(badKind)
	}
	return len(t.b.Elems())
}

func (t Type) Field(i int) Type {
	if t.Kind() != Struct {
		panic(badKind)
	}
	return Type{t.b.Elems()[i]}
}
