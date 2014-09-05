__attribute__ ((always_inline))
extern inline
unsafe$Pointer stack$alloc(int n, uintptr size) {
	unsafe$Pointer p = __builtin_alloca(size);
	__builtin_memset(p, 0, size);
	return p;
}

__attribute__ ((always_inline))
extern inline
slice stack$bytes(int n) {
	return (slice){stack$alloc(n, 1), n, n};
}

__attribute__ ((always_inline))
extern inline
slice stack$ints8(int n) {
	return (slice){stack$alloc(n, 1), n, n};
}

__attribute__ ((always_inline))
extern inline
slice stack$ints16(int n) {
	return (slice){stack$alloc(n, 2), n, n};
}

__attribute__ ((always_inline))
extern inline
slice stack$uints16(int n) {
	return (slice){stack$alloc(n, 2), n, n};
}
 
__attribute__ ((always_inline))
extern inline
slice stack$ints32(int n) {
	return (slice){stack$alloc(n, 4), n, n};
}

__attribute__ ((always_inline))
extern inline
slice stack$uints32(int n) {
	return (slice){stack$alloc(n, 4), n, n};
}