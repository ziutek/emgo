package reflect

import (
	"builtin"
	"unsafe"
)

type Value struct {
	val complex128
	typ Type
}

// ValueOf returns a new Value initialized to the concrete value stored in i.
// ValueOf(nil) returns the zero Value.
func ValueOf(i interface{}) Value {
	return *(*Value)(unsafe.Pointer(&i))
}

// Zero returns value that represents zero value of type t.
func Zero(t Type) Value {
	return Value{typ: t}
}

// IsValid returns true if v represents a value. It returns false if v is zero
// Value.
func (v Value) IsValid() bool {
	return v.typ.IsValid()
}

// Type returns type of v.
func (v Value) Type() Type {
	return v.typ
}

// Kind returns kind od v. If v is zero Value, Kind returns Invalid.
func (v Value) Kind() Kind {
	return v.typ.Kind()
}

func (v Value) Bool() bool {
	if v.Kind() != Bool {
		panic("reflect: not bool")
	}
	return *(*bool)(unsafe.Pointer(&v.val))
}

// Int returns underlying value of v as an int64.
// It panics if kind of v isn't Int, Int8, Int16, Int32, Int64.
func (v Value) Int() int64 {
	p := unsafe.Pointer(&v.val)
	switch v.Kind() {
	case Int:
		return int64(*(*int)(p))
	case Int8:
		return int64(*(*int8)(p))
	case Int16:
		return int64(*(*int16)(p))
	case Int32:
		return int64(*(*int32)(p))
	case Int64:
		return *(*int64)(p)
	}
	panic("reflect: not signed int")
}

// Uint returns underlying value of v as an uint64.
// It panics if kind of v isn't Uint, Uint8, Uint16, Uint32, Uint64.
func (v Value) Uint() uint64 {
	p := unsafe.Pointer(&v.val)
	switch v.Kind() {
	case Uint:
		return uint64(*(*uint)(p))
	case Uint8:
		return uint64(*(*uint8)(p))
	case Uint16:
		return uint64(*(*uint16)(p))
	case Uint32:
		return uint64(*(*uint32)(p))
	case Uint64:
		return *(*uint64)(p)
	}
	panic("reflect: not unsigned int")
}

// Float returns underlying value of v as a float64.
// It panics if kind of v isn't Float32, Float64.
func (v Value) Float() float64 {
	p := unsafe.Pointer(&v.val)
	switch v.Kind() {
	case Float32:
		return float64(*(*float32)(p))
	case Float64:
		return *(*float64)(p)
	}
	panic("reflect: not float")
}

// Complex returns underlying value of v as a complex128.
// It panics if kind of v isn't Complex64, Complex128.
func (v Value) Complex() complex128 {
	p := unsafe.Pointer(&v.val)
	switch v.Kind() {
	case Complex64:
		return complex128(*(*complex64)(p))
	case Complex128:
		return *(*complex128)(p)
	}
	panic("reflect: not complex")
}

// Pointer returns underlying value of v as an uintptr.
// It panic if kind of v isn't Chan, Func, Map, Ptr, Slice or UnsafePointer.
func (v Value) Pointer() uintptr {
	p := unsafe.Pointer(&v.val)
	switch v.Kind() {
	case Func, Ptr, UnsafePointer:
		return *(*uintptr)(p)
	case Chan:
		return uintptr((*builtin.Chan)(p).C)
	case Map:
		// BUG: Not implemented
	case Slice:
		s := *(*[]byte)(p)
		return uintptr(unsafe.Pointer(&s[0]))
	}
	panic("reflect: not chan, func, ptr, slice")
}

// String returns underlying value of v as a string.
// It panic if kind of v isn't String.
func (v Value) String() string {
	if v.Kind() != String {
		panic("reflect: not string")
	}
	return *(*string)(unsafe.Pointer(&v.val))
}

const notASCMS = "reflect: not array, slice, chan, map, string"

func panicASCMS() {
	panic(notASCMS)
}

func panicASC() {
	panic(notASCMS[:25])
}

func (v Value) Len() int {
	p := unsafe.Pointer(&v.val)
	switch v.Kind() {
	case Array:
		return v.Type().Len()
	case Slice:
		return len(*(*[]byte)(p))
	case Chan:
		return len(*(*chan byte)(p))
	case Map:
		// BUG: Not implemented
		return -1
	case String:
		return len(*(*string)(p))
	}
	panicASCMS()
	return 0
}

/*func (v Value) Index() Value {
	return
}*/
