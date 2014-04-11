extern byte HeapBegin, HeapSize;

static inline
__slice runtime_noos_heap() {
	uint size = (uint)&HeapSize;
	return (__slice){&HeapBegin, size, size};
}

static inline
uintptr alignUp(uintptr p, uintptr a) {
	--a;
	return (p + a) & ~a;
}

static inline
uintptr alignDown(uintptr p, uintptr a) {
	return p & ~(a - 1);
}

static inline
__slice runtime_noos_allocBottom(unsafe_Pointer sptr, __slice bs, int n, uintptr elSize, uintptr elAlign, uintptr sliAlign) {
	uint blen = (uint)n * (uint)alignUp(elSize, elAlign);
	__slice *p = (__slice*)sptr;
	
	p->arr = (unsafe_Pointer)alignUp((uintptr)bs.arr, sliAlign);
	p->len = (uint)n;
	p->cap = (uint)n;
	blen += (uint)(p->arr - bs.arr);
	
	if (blen > bs.len) {
		return __NILSLICE;
	}
	return (__slice){bs.arr + blen, bs.len - blen, bs.cap - blen};
} 

static inline
__slice runtime_noos_allocTop(unsafe_Pointer sptr, __slice bs, int n, uintptr elSize, uintptr elAlign, uintptr sliAlign) {
	uint blen = (uint)n * (uint)alignUp(elSize, elAlign);
	__slice *p = (__slice*)sptr;
	
	unsafe_Pointer arr = bs.arr + bs.len - blen;
	p->arr = (unsafe_Pointer)alignDown((uintptr)arr, sliAlign);
	p->len = (uint)n;
	p->cap = (uint)n;
	blen += (uint)(arr - p->arr);
	
	if (blen > bs.len) {
		return __NILSLICE;
	}
	return (__slice){bs.arr, bs.len - blen, bs.len - blen};
} 
