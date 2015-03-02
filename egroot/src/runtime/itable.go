package runtime

import (
	"sync/atomic"
	"sync/barrier"
	"builtin"
	"unsafe"
)

type itable []unsafe.Pointer

func (it itable) head() *ItHead {
	return (*ItHead)(unsafe.Pointer(&it[0]))
}

const headLen = int(unsafe.Sizeof(builtin.Itable{}) -
	unsafe.Sizeof(builtin.Itable{}.Methods)) / unsafe.Sizeof(uintptr(0))

func (it itable) methods() []unsafe.Pointer {
	return itable[headLen:]
}

// makeItable makes new Itable. etyp must be assignable to ityp.
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

func hash(ityp, etyp *Type) int {
	h := uintptr(unsafe.Pointer(ityp)) ^ uintptr(unsafe.Pointer(etyp))
	return int(h) & (len(itHashTab) - 1)
}

type itListElem struct {
	ityp, etyp *Type
	next       unsafe.Pointer // *itListElem
	itab       itable
}

var itHashTab [1 << 3]unsafe.Pointer // *itListElem

// GetItable returns itable for given interface and non-interface type pair. If
// etyp isn't assignable to ityp behavior of GetItable is undefined.
func GetItable(ityp, etyp *Type) itable {
	// Try find itable in hash table.
	list := &itHashTab[hash(ityp, etyp)]
	for {
		elem := (*itListElem)(atomic.LoadPointer(list))
		if elem == nil {
			break
		}
		if elem.ityp == ityp && elem.etyp == etyp {
			return elem.itab
		}
		list = &elem.next
	}
	// Not found. Make and add new one to the list.
	newel := new(itListElem)
	newel.ityp = ityp
	newel.etyp = etyp
	newel.itab = makeItable(ityp, etyp)
	barrier.Memory()
	for {
		if atomic.CompareAndSwap(list, nil, unsafe.Pointer(newel)) {
			return newel.itab
		}

	}
}
