extern byte HeapBegin, HeapSize;

static inline
slice runtime$noos$heap() {
	uint size = (uint)&HeapSize;
	return (slice){&HeapBegin, size, size};
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

static
slice runtime$noos$allocBottom(unsafe$Pointer sptr, slice bs, int n, uintptr elSize, uintptr elAlign, uintptr sliAlign) {
	uint blen = (uint)n * (uint)alignUp(elSize, elAlign);
	slice *p = (slice*)sptr;
	
	p->arr = (unsafe$Pointer)alignUp((uintptr)bs.arr, sliAlign);
	p->len = (uint)n;
	p->cap = (uint)n;
	memset(p->arr, 0, blen);
	
	blen += (uint)(p->arr - bs.arr);
	if (blen > bs.len) {
		return NILSLICE;
	}
	return (slice){bs.arr + blen, bs.len - blen, bs.cap - blen};
} 

static
slice runtime$noos$allocTop(unsafe$Pointer sptr, slice bs, int n, uintptr elSize, uintptr elAlign, uintptr sliAlign) {
	uint blen = (uint)n * (uint)alignUp(elSize, elAlign);
	slice *p = (slice*)sptr;
	
	unsafe$Pointer arr = bs.arr + bs.len - blen;
	p->arr = (unsafe$Pointer)alignDown((uintptr)arr, sliAlign);
	p->len = (uint)n;
	p->cap = (uint)n;
	memset(p->arr, 0, blen);
	
	blen += (uint)(arr - p->arr);
	if (blen > bs.len) {
		return NILSLICE;
	}
	return (slice){bs.arr, bs.len - blen, bs.len - blen};
} 
