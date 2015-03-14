package reflect

import (
	"builtin"
)

type Value struct {
	val complex128
	typ *builtin.Type
}

func valueOf(i interface{}) Value

func ValueOf(i interface{}) Value {
	return valueOf(i)
}

func (v Value) Type() Type {
	return Type{v.typ}
}
