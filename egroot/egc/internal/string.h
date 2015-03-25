typedef struct {
	byte *str;
	uint len;
} string;

#define EGSTR(s) ((string){(byte*)(s), sizeof(s)-1})
#define EGISTR(s) {(byte*)(s), sizeof(s)-1}

// TODO: #define NEWSTR(b)

#define _SSLICE_LOW(expr, low) \
	string s = expr;           \
	uint l = low;              \
	s.str = s.str + l

#define SSLICEL(expr, low) ({ \
	_SSLICE_LOW(expr, low);   \
	s.len -= l;               \
	s;                        \
})

#define SSLICELH(expr, low, high) ({ \
	_SSLICE_LOW(expr, low);          \
	s.len = high - l;                \
	s;                               \
})

#define SSLICEH(expr, high) ({ \
	string s = expr;           \
	s.len = high;              \
	s;                         \
})

#define STRCPY(dst, src) ({                          \
	int n = (dst.len < src.len) ? dst.len : src.len; \
	memmove(dst.arr, src.str, n);                    \
	n;                                               \
})

#define BYTES(s) ({                                           \
	slice b = (slice){__builtin_alloca(s.len), s.len, s.len}; \
	__builtin_memcpy(b.arr, s.str, s.len);                    \
	b;                                                        \
})
