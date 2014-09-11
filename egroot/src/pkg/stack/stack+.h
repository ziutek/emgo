__attribute__ ((always_inline))
extern inline
unsafe$Pointer stack$alloc(int n, uintptr size) {
	unsafe$Pointer p = __builtin_alloca(size);
	__builtin_memset(p, 0, size);
	return p;
}

#define _DEFFUNC(typ)                                      \
	__attribute__ ((always_inline))                        \
	extern inline                                          \
	slice stack$##typ##s(int n) {                            \
		return (slice){stack$alloc(n, sizeof(typ)), n, n}; \
	}
	
_DEFFUNC(int)
_DEFFUNC(uint)
_DEFFUNC(uintptr)
_DEFFUNC(bool)
_DEFFUNC(interface)

#undef _DEFFUNC

#define _DEFFUNC(typ, bits)                           \
	__attribute__ ((always_inline))                   \
	extern inline                                     \
	slice stack$##typ##s##bits(int n) {               \
		return (slice){stack$alloc(n, bits/8), n, n}; \
	}                                                 \

_DEFFUNC(int, 8)
_DEFFUNC(int, 16)
_DEFFUNC(int, 32)
_DEFFUNC(int, 64)

_DEFFUNC(uint, 8)
_DEFFUNC(uint, 16)
_DEFFUNC(uint, 32)
_DEFFUNC(uint, 64)

_DEFFUNC(float, 32)
_DEFFUNC(float, 64)

_DEFFUNC(complex, 64)
_DEFFUNC(complex, 128)

#undef _DEFFUNC