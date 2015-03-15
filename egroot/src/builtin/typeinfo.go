package builtin

import "unsafe"

type Method struct {
	_ byte
}

type Type struct {
	name    string
	size    uintptr
	kind    byte
	elems   []*Type `C:"const"`
	methods []*Method
	fns     [0]unsafe.Pointer
}

func (t *Type) Kind() byte {
	return t.kind
}

func (t *Type) Size() uintptr {
	return t.size
}

func (t *Type) Name() string {
	return t.name
}

func (t *Type) Elems() []*Type {
	return t.elems
}

func (t *Type) Methods() []*Method {
	return t.methods
}

func (t *Type) Fns() []unsafe.Pointer {
	return t.fns[:len(t.methods)]
}

type ItHead struct {
	typ *Type `C:"const"`
}

func (ith *ItHead) Type() *Type {
	// ith.typ is read-only (can be plced in ROM).
	return (*Type)(unsafe.Pointer(ith.typ))
}

type Itable struct {
	head ItHead
	fns  [0]unsafe.Pointer
}

func (it *Itable) Head() *ItHead {
	return &it.head
}

func (it *Itable) Fns() []unsafe.Pointer {
	return it.fns[:len(it.head.typ.methods)]
}

// NewItable allocates and initializes new itable.
// etyp must be assignable to ityp.
func NewItable(ityp, etyp *Type) *Itable {
	const hlen = int((unsafe.Sizeof(ItHead{}) + unsafe.Sizeof(uintptr(0)) - 1) / unsafe.Sizeof(uintptr(0)))
	sli := make([]uintptr, hlen+len(ityp.methods))
	itab := (*Itable)(unsafe.Pointer(&sli[0]))
	itab.head.typ = etyp
	e := 0
	itabFns := itab.Fns()
	etypFns := etyp.Fns()
	for i, m := range ityp.methods {
		for etyp.methods[e] != m {
			e++
		}
		itabFns[i] = etypFns[e]
		e++
	}
	return itab
}

// Implements returns true if t has all methods that ityp has. If ityp is
// interface type then Implements returns true if t implements ityp.
func (t *Type) Implements(ityp *Type) bool {
	if len(ityp.methods) == 0 {
		return true
	}
	if t == nil || len(t.methods) < len(ityp.methods) {
		return false
	}
	k := 0
	for _, im := range ityp.methods {
		for {
			if k >= len(t.methods) {
				return false
			}
			m := t.methods[k]
			k++
			if m == im {
				break
			}
		}
	}
	return true
}

// ItableFor should return itable for given interface and non-interface type
// pair. It is always called with etyp assignable to ityp.
// To allow assign/assert to interfaces in interrupt handlers ItableFor must
// be implemented in nonblocking way.
var ItableFor func(ityp, etyp *Type) *Itable

type generateBasicTinfos struct {
	_ *bool
	_ *int
	_ *int8
	_ *int16
	_ *int32
	_ *int64
	_ *uint
	_ *uint8
	_ *uint16
	_ *uint32
	_ *uint64
	_ *uintptr
	_ *float32
	_ *float64
	_ *complex64
	_ *complex128
	_ *string
	_ *unsafe.Pointer
	_ *error
	_ []bool
	_ []int
	_ []int8
	_ []int16
	_ []int32
	_ []int64
	_ []uint
	_ []uint8
	_ []uint16
	_ []uint32
	_ []uint64
	_ []uintptr
	_ []float32
	_ []float64
	_ []complex64
	_ []complex128
	_ []string
	_ []unsafe.Pointer
	_ []error
}
