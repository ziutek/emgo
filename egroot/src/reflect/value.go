package reflect

type Value struct {
	val complex128
	typ Type
}

func valueOf(i interface{}) Value

// ValueOf returns a new Value initialized to the concrete value stored in i. ValueOf(nil) returns the zero Value.
func ValueOf(i interface{}) Value {
	return valueOf(i)
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