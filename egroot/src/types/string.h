typedef struct {
	byte *str;
	uint len;
} string;

#define __EGSTR(s) ((string){(byte *)s, sizeof(s)-1})

// #define __NEWSTR(b) Memory allocation need.

#define __SSLICE_LOW(expr, low) \
	string __s = expr;          \
	uint __low = low;           \
	__s.str = __s.str + __low

#define __SSLICEL(expr, low) ({ \
	__SSLICE_LOW(expr, low);    \
	__s.len -= low;             \
	__s;                        \
})

#define __SSLICELH(expr, low, high) ({ \
	__SSLICE_LOW(expr, low);           \
	__s.len = high - __low;            \
	__s;                               \
})

#define __SSLICEH(expr, high) ({ \
	string __s = expr;           \
	__s.len = high;              \
	__s;                         \
})

#define __STRCPY(dst, src) ({                     \
	int __n = (dst.len < src.len) ? dst.len : src.len; \
	memmove(dst.arr, src.str, __n); \
	__n;                                               \
})
