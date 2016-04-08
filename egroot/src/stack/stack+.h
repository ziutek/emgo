__attribute__ ((always_inline))
static inline
unsafe$Pointer stack$alloc(int n, uintptr size) {
	size *= n;
	unsafe$Pointer p = __builtin_alloca(size);
	__builtin_memset(p, 0, size);
	return p;
}

__attribute__ ((always_inline))
static inline
uintptr stack$Alloc(int n, uintptr size) {
	return (uintptr)(stack$alloc(n, size));
}


#define _DEFFUNC(name, typ)                                \
	__attribute__ ((always_inline))                        \
	static inline                                          \
	slice stack$##name(int n) {                            \
		return (slice){stack$alloc(n, sizeof(typ)), n, n}; \
	}

_DEFFUNC(Bytes, byte)
_DEFFUNC(Ints, int)
_DEFFUNC(Uints, uint)
_DEFFUNC(Uintptrs, uintptr)
_DEFFUNC(Pointers, unsafe$Pointer)
_DEFFUNC(Bools, bool)
_DEFFUNC(Interfaces, interface)

#undef _DEFFUNC

#define _DEFFUNC(name, bits)                          \
	__attribute__ ((always_inline))                   \
	static inline                                     \
	slice stack$##name##bits(int n) {                 \
		return (slice){stack$alloc(n, bits/8), n, n}; \
	}

_DEFFUNC(Ints, 8)
_DEFFUNC(Ints, 16)
_DEFFUNC(Ints, 32)
_DEFFUNC(Ints, 64)

_DEFFUNC(Uints, 16)
_DEFFUNC(Uints, 32)
_DEFFUNC(Uints, 64)

_DEFFUNC(Floats, 32)
_DEFFUNC(Floats, 64)

_DEFFUNC(Complexs, 64)
_DEFFUNC(Complexs, 128)

#undef _DEFFUNC