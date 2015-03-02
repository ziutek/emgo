package builtin

import "unsafe"

type itable []unsafe.Pointer

const headLen = int(unsafe.Sizeof(ItHead{}) / unsafe.Sizeof(uintptr(0)))

func (it itable) head() *ItHead {
	return (*ItHead)(unsafe.Pointer(&it[0]))
}

func (it itable) methods() []unsafe.Pointer {
	return it[headLen:]
}

func mequal(a, b *Method) bool {
	if len(a.param) != len(b.param) {
		return false
	}
	if a.name != b.name {
		return false
	}
	for i, t := range a.param {
		if t != b.param[i] {
			return false
		}
	}
	return true
}

// makeItable makes new Itable. It requires that etyp can be assigned to ityp.
func makeItable(ityp, etyp *Type) itable {
	it := make(itable, headLen+len(ityp.mset))
	it.head().Type = etyp
	itm := it.methods()
	e := 0
	for i := range ityp.mset {
		for !mequal(&ityp.mset[i], &etyp.mset[e]) {
			e++
		}
		itm[i] = etyp.mset[e].fptr
		e++
	}
	return it
}

type itListElem struct {
	ityp, etyp *Type
	next       *itListElem
	itab       itable
}

var itHashTab [1 << 3]*itListElem

func hash(ityp, etyp *Type) int {
	h := uintptr(unsafe.Pointer(ityp)) ^ uintptr(unsafe.Pointer(etyp))
	return int(h) & (len(itHashTab) - 1)
}

// GetItable returns itable for given interface, non-interface type pair. It
// allocates new one if there is no itable for given pair.
func GetItable(ityp, etyp *Type) itable {

	return nil
}
