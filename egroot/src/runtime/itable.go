package runtime

import (
	"builtin"
	"sync/atomic"
	"unsafe"
)

func hash(ityp, etyp *builtin.Type) int {
	h := uintptr(unsafe.Pointer(ityp)) ^ uintptr(unsafe.Pointer(etyp))
	return int(h) & (len(itHashTab) - 1)
}

type itListElem struct {
	ityp, etyp *builtin.Type
	next       unsafe.Pointer // *itListElem
	itab       *builtin.Itable
}

var itHashTab [1 << 3]unsafe.Pointer // *itListElem

// itableFor implements builtin.ItableFor.
func itableFor(ityp, etyp *builtin.Type) *builtin.Itable {
	// Find itable in hash table.
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
	newel.itab = builtin.NewItable(ityp, etyp)
	for {
		if atomic.CompareAndSwapPointer(list, nil, unsafe.Pointer(newel)) {
			return newel.itab
		}
		for {
			elem := (*itListElem)(atomic.LoadPointer(list))
			if elem == nil {
				break
			}
			if elem.ityp == ityp && elem.etyp == etyp {
				// BUG: newel allocated but not used (memory leak if no GC).
				return elem.itab
			}
			list = &elem.next
		}

	}
}

func init() {
	builtin.ItableFor = itableFor
}
