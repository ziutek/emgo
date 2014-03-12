// +build none

// All external symbols as byte to prevent compiler to optimize
// any runtime align checks.
extern byte DataRAM, DataLoad, DataSize;
extern byte BSSRAM, BSSSize;
extern byte FreeStart, FreeEnd, FreeSize, HeapSize;

__attribute__ ((always_inline))
extern inline uintptr runtime_freeStart() {
	return (uintptr)&FreeStart;
}

__attribute__ ((always_inline))
extern inline uintptr runtime_freeEnd() {
	return (uintptr)&FreeEnd;
}

__attribute__ ((always_inline))
extern inline uintptr runtime_freeSize() {
	return (uintptr)&FreeSize;
}

__attribute__ ((always_inline))
extern inline uintptr runtime_HeapSize() {
	return (uintptr)&HeapSize;
}

__attribute__ ((always_inline))
extern inline void runtime_setSlice(unsafe_Pointer sptr, unsafe_Pointer addr, uint len, uint cap) {
	__slice *p = (__slice*)sptr;
	p->arr = addr;
	p->len = len;
	p->cap = cap;
}