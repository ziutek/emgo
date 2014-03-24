// All external symbols as byte to prevent compiler to optimize
// any runtime align checks.
extern byte HeapStackBegin, HeapStackEnd, HeapStackSize, StackSize;

static inline uintptr alignUp(uintptr p, uintptr a) {
	uintptr mask = a - 1;
	if (p&mask != 0) {
		p = (p + a) &~ mask;
	}
	return p;
}

static inline uintptr alignDown(uintptr p, uintptr a) {
	return p &~ (a - 1);
}

static inline uintptr runtime_noos_stackSize() {
	return (uintptr)&StackSize;
}

static inline uintptr runtime_noos_heapStackEnd() {
	return (uintptr)&HeapStackEnd;
}

static inline __slice runtime_noos_heapStack() {
	uint size = (uint)&HeapStackSize;
	return (__slice){&HeapStackBegin, size, size};
}

static __slice runtime_noos_allocBottom(unsafe_Pointer sptr, __slice bs, int n, uintptr size, uintptr align) {
	uint blen = (uint)n * size;
	__slice *p = (__slice*)sptr;
	
	p->arr = (unsafe_Pointer)alignUp((uintptr)bs.arr, align);
	p->len = (uint)n;
	p->cap = (uint)n;
	blen += (uint)(p->arr - bs.arr);
	
	if (blen > bs.len) {
		return __NILSLICE;
	}
	return (__slice){bs.arr + blen, bs.len - blen, bs.cap - blen};
} 

static __slice runtime_noos_allocTop(unsafe_Pointer sptr, __slice bs, int n, uintptr size, uintptr align) {
	uint blen = (uint)n * size;
	__slice *p = (__slice*)sptr;
	
	unsafe_Pointer arr = bs.arr + bs.len - blen;
	p->arr = (unsafe_Pointer)alignDown((uintptr)arr, align);
	p->len = (uint)n;
	p->cap = (uint)n;
	blen += (uint)(arr - p->arr);
	
	if (blen > bs.len) {
		return __NILSLICE;
	}
	return (__slice){bs.arr, bs.len - blen, bs.cap};
} 
