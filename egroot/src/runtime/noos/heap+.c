extern byte HeapBegin, HeapEnd;

static inline
uintptr runtime$noos$heapBegin() {
	return (uintptr)(&HeapBegin);
}

static inline
uintptr runtime$noos$heapEnd() {
	return (uintptr)(&HeapEnd);
}
