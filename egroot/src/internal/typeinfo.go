package internal

import "unsafe"

type Method struct {
	_ byte
}

type StructField struct {
	name string
	Type *Type
}

func (f StructField) Name() string {
	if f.Type == nil {
		return ""
	}
	return f.name
}

func (f StructField) Exported() bool {
	return f.Type != nil && f.name[0] >= 'A' && f.name[0] <= 'Z'
}

type savedSize struct{ align, size uintptr }

func (f StructField) SavedSize() uintptr {
	return (*savedSize)(unsafe.Pointer(&f.name)).size
}

func (f StructField) SavedAlign() uintptr {
	return (*savedSize)(unsafe.Pointer(&f.name)).align
}

type Type struct {
	name    string
	kind    int
	elems   unsafe.Pointer
	elemN   uintptr
	methods unsafe.Pointer
	methodN uintptr
	fns     [0]unsafe.Pointer
}

func (t *Type) Kind() int {
	if t.kind < 0 {
		// Array
		return -1
	}
	return t.kind
}

func (t *Type) Len() int {
	if t.kind < 0 {
		return -1 - t.kind
	}
	return 0
}

func (t *Type) Name() string {
	return t.name
}

type elemKey struct{ elem, key *Type }

func (t *Type) Elem() *Type {
	return (*Type)(t.elems)
}

func (t *Type) Key() *Type {
	return (*Type)(unsafe.Pointer(t.elemN))
}

func (t *Type) Fields() []StructField {
	return (*[1 << 24]StructField)(t.elems)[:t.elemN]
}

func (t *Type) Methods() []*Method {
	return (*[1 << 28]*Method)(t.methods)[:t.methodN]
}

func (t *Type) Fns() []unsafe.Pointer {
	return (*[1 << 28]unsafe.Pointer)(unsafe.Pointer(&t.fns))[:t.methodN]
}

type ItHead struct {
	typ *Type
}

func (ith *ItHead) Type() *Type {
	return ith.typ
}

type Itable struct {
	head ItHead
	fns  [0]unsafe.Pointer
}

func (it *Itable) Head() *ItHead {
	return &it.head
}

func (it *Itable) Fns() []unsafe.Pointer {
	return (*[1 << 28]unsafe.Pointer)(unsafe.Pointer(&it.fns))[:it.head.typ.methodN]
}

// NewItable allocates and initializes new itable.
// etyp must be assignable to ityp.
func NewItable(ityp, etyp *Type) *Itable {
	const hlen = (unsafe.Sizeof(ItHead{}) + unsafe.Sizeof(uintptr(0)) - 1) /
		unsafe.Sizeof(uintptr(0))
	sli := make([]uintptr, hlen+ityp.methodN)
	itab := (*Itable)(unsafe.Pointer(&sli[0]))
	itab.head.typ = etyp
	e := 0
	itabFns := itab.Fns()
	etypFns := etyp.Fns()
	etm := etyp.Methods()
	for i, m := range ityp.Methods() {
		for etm[e] != m {
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
	if t == nil {
		return false
	}
	if ityp.methodN == 0 {
		return true
	}
	if t.methodN < ityp.methodN {
		return false
	}
	k := uintptr(0)
	tm := t.Methods()
	for _, im := range ityp.Methods() {
		for {
			if k >= t.methodN {
				return false
			}
			m := tm[k]
			k++
			if m == im {
				break
			}
		}
	}
	return true
}

// ItableFor should return itable for given interface and non-interface type
// pair. It is always called with etyp assignable to ityp. To allow
// assign/assert to interfaces in interrupt handlers ItableFor must be
// implemented in nonblocking way.
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
