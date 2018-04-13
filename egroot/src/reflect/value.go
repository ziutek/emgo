package reflect

import (
	"internal"
	"mem"
	"unsafe"
)

const (
	// flagIndir is set if Value.val doesn't store real value of type
	// Value.typ but pointer to it.
	flagIndir byte = 1 << iota
)

type Value struct {
	val   ival
	typ   Type
	flags byte
}

type emptyI struct {
	val ival
	typ Type
}

// ValueOf returns a new Value initialized to the concrete value stored in i.
// ValueOf(nil) returns the zero Value.
func ValueOf(i interface{}) Value {
	e := *(*emptyI)(unsafe.Pointer(&i))
	return Value{val: e.val, typ: e.typ}
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

func (v *Value) asptr() unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Pointer(&v.val))
}

// ptrto always returns pointer to real value that v stores or refers to.
func (v *Value) ptrto() unsafe.Pointer {
	if v.flags&flagIndir == 0 {
		return unsafe.Pointer(&v.val)
	}
	return v.asptr()
}

func (v Value) Bool() bool {
	if v.Kind() != Bool {
		panic(badKind)
	}
	return *(*bool)(v.ptrto())
}

// Int returns underlying value of v as an int64.
// It panics if kind of v isn't Int, Int8, Int16, Int32, Int64.
func (v Value) Int() int64 {
	pt := v.ptrto()
	switch v.Kind() {
	case Int:
		return int64(*(*int)(pt))
	case Int8:
		return int64(*(*int8)(pt))
	case Int16:
		return int64(*(*int16)(pt))
	case Int32:
		return int64(*(*int32)(pt))
	case Int64:
		return *(*int64)(pt)
	}
	panic(badKind)
}

// Uint returns underlying value of v as an uint64.
// It panics if kind of v isn't Uint, Uint8, Uint16, Uint32, Uint64.
func (v Value) Uint() uint64 {
	pt := v.ptrto()
	switch v.Kind() {
	case Uint:
		return uint64(*(*uint)(pt))
	case Uintptr:
		return uint64(*(*uintptr)(pt))
	case Uint8:
		return uint64(*(*uint8)(pt))
	case Uint16:
		return uint64(*(*uint16)(pt))
	case Uint32:
		return uint64(*(*uint32)(pt))
	case Uint64:
		return *(*uint64)(pt)
	}
	panic(badKind)
}

// Float returns underlying value of v as a float64.
// It panics if kind of v isn't Float32, Float64.
func (v Value) Float() float64 {
	pt := v.ptrto()
	switch v.Kind() {
	case Float32:
		return float64(*(*float32)(pt))
	case Float64:
		return *(*float64)(pt)
	}
	panic(badKind)
}

// Complex returns underlying value of v as a complex128.
// It panics if kind of v isn't Complex64, Complex128.
func (v Value) Complex() complex128 {
	pt := v.ptrto()
	switch v.Kind() {
	case Complex64:
		return complex128(*(*complex64)(pt))
	case Complex128:
		return *(*complex128)(pt)
	}
	panic(badKind)
}

// Pointer returns underlying value of v as an uintptr.
// It panics if kind of v is not Chan, Func, Map, Ptr, Slice or UnsafePointer.
func (v Value) Pointer() uintptr {
	pt := v.ptrto()
	switch v.Kind() {
	case Ptr, UnsafePointer, Func, Slice:
		return *(*uintptr)(pt)
	case Chan:
		return uintptr((*internal.Chan)(pt).C)
	case Map:
		// BUG: Not implemented
	}
	panic(badKind)
}

// IsNil returns true if underlying value of v is nil. It panics if kind of v
// isn't Interface, Chan, Func, Map, Ptr, Slice or UnsafePointer
func (v Value) IsNil() bool {
	if v.Kind() == Interface {
		return !(*emptyI)(v.ptrto()).typ.IsValid()
	}
	return v.Pointer() == 0
}

func (v Value) Elem() Value {
	switch v.Kind() {
	case Ptr:
		v.typ = v.typ.Elem()
		if v.flags&flagIndir == 0 {
			v.flags |= flagIndir
		} else {
			*(*unsafe.Pointer)(unsafe.Pointer(&v.val)) = *(*unsafe.Pointer)(v.asptr())
		}
		return v
	case Interface:
		// TODO
		break
	}
	panic(badKind)
}

// String returns underlying value of v as a string.
// It panics if kind of v isn't String.
func (v Value) String() string {
	if v.Kind() != String {
		panic(badKind)
	}
	return *(*string)(v.ptrto())
}

const badIndex = "reflect: index out of range"

func (v Value) Index(i int) Value {
	if uint(i) >= uint(v.Len()) {
		panic(badIndex)
	}
	switch k := v.Kind(); k {
	case Slice, Array:
		var ptr uintptr
		if k == Array {
			ptr = uintptr(v.ptrto())
		} else {
			ptr = v.Pointer()
		}
		r := Value{typ: v.Type().Elem(), flags: v.flags | flagIndir}
		ptr += mem.AlignUp(r.typ.Size(), r.typ.Align()) * uintptr(i)
		*(*unsafe.Pointer)(unsafe.Pointer(&r.val)) = unsafe.Pointer(ptr)
		return r
	case String:
		return ValueOf(v.String()[i])
	}
	panic(badKind)
}

func (v Value) Len() int {
	pt := v.ptrto()
	switch v.Kind() {
	case Array:
		return v.Type().Len()
	case Slice:
		return len(*(*[]byte)(pt))
	case Chan:
		return len(*(*chan byte)(pt))
	case Map:
		// BUG: Not implemented
		return -1
	case String:
		return len(*(*string)(pt))
	}
	panic(badKind)
}

func (v Value) Cap() int {
	pt := v.ptrto()
	switch v.Kind() {
	case Array:
		return v.Type().Len()
	case Slice:
		return cap(*(*[]byte)(pt))
	case Chan:
		return cap(*(*chan byte)(pt))
	}
	panic(badKind)
}

// Interfce returns underlying value of v as interfce{}. It returns nil if v
// isn't valid or underlying value can't be assigned to interface{}.
func (v Value) Interface() interface{} {
	if !v.IsValid() {
		return nil
	}
	ei := emptyI{typ: v.Type()}
	size := ei.typ.Size()
	if size > unsafe.Sizeof(ei.val) {
		return nil
	}
	internal.Memmove(unsafe.Pointer(&ei.val), v.ptrto(), size)
	return *(*interface{})(unsafe.Pointer(&ei))
}

func (v Value) NumField() int {
	return v.Type().NumField()
}

func (v Value) Field(i int) Value {
	t := v.Type()
	if t.Kind() != Struct {
		panic(badKind)
	}
	if uint(i) >= uint(t.NumField()) {
		panic(badIndex)
	}
	rt := t.Field(i).Type()
	if !rt.IsValid() {
		return Value{}
	}
	ptr := uintptr(v.ptrto())
	for k := 0; k < i; k++ {
		ptr += mem.AlignUp(t.Field(k).Size(), t.Field(k+1).Align())
	}
	r := Value{typ: rt, flags: v.flags | flagIndir}
	*(*unsafe.Pointer)(unsafe.Pointer(&r.val)) = unsafe.Pointer(ptr)
	return r
}
