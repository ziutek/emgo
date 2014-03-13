// All external symbols as byte to prevent compiler to optimize
// any runtime align checks.
extern byte HeapStackBegin, HeapStackEnd, HeapStackSize, HeapSize;

static inline uintptr runtime_noos_heapSize() {
	return (uintptr)&HeapSize;
}

static inline __slice runtime_noos_heapStack() {
	uint size = (uint)&HeapStackSize;
	return (__slice){&HeapStackBegin, size, size};
}

static inline __slice runtime_noos_alloc(unsafe_Pointer sptr, __slice bs, int n, uintptr size) {
	__slice *p = (__slice*)sptr;
	p->arr = bs.arr;
	p->len = (uint)n;
	p->cap = (uint)n;
	uint blen = (uint)n * size;
	return (__slice){bs.arr + blen, bs.len - blen, bs.cap - blen};
} 

/*static inline __slice runtime_noos_sliceU8(unsafe_Pointer p, n uint) {
	return (__slice){p, n, n};
}

static inline __slice runtime_noos_sliceU16(unsafe_Pointer p, n uint) {
	runtime_noos_checkAlignment(p, 2);
	return (__slice){p, n, n};
}

static inline __slice runtime_nnos_sliceU32(unsafe_Pointer p, n uint) {
	runtime_noos_checkAlignment(p, 4);
	return (__slice){p, n, n};
}*/