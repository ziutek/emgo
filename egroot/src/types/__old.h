typedef __INT8_TYPE__ int8;
typedef __INT16_TYPE__ int16;
typedef __INT32_TYPE__ int32;
typedef __INT64_TYPE__ int64;

typedef __UINT8_TYPE__ byte;
typedef __UINT8_TYPE__ uint8;
typedef __UINT16_TYPE__ uint16;
typedef __UINT32_TYPE__ uint32;
typedef __UINT64_TYPE__ uint64;
typedef unsigned int uint;

typedef __UINTPTR_TYPE__ uintptr;

typedef void* unsafe_Pointer;

typedef float float32;
typedef double float64;

typedef float _Complex complex64;
typedef double _Complex complex128;

#ifndef EG_ALLOW_C_TYPES

// Forget about C types
#define	unsigned	XXunsigned / / /
#define	signed		XXsigned / / /
#define	char		XXchar / / /
#define	short		XXshort / / /
#define	long		XXlong / / /
#define	float		XXfloat / / /
#define	double		XXdouble / / /
#define _Complex	XX_Complex / / /

#endif

typedef uint8 bool;

#define true (1)
#define false (0)

#define complex(re, im) ((re)+1.i*(im))
#define real(c) (__real__(c))
#define imag(c) (__imag__(c))

typedef struct {
	byte *str;
	uint len;
} string;

#define __EGSTR(s) ((string){(byte *)s, sizeof(s)-1})

typedef struct {
	unsafe_Pointer array;
	uint len;
	uint cap;
} __slice;

#define len(v) (v.len)

#define __ALEN(a) (sizeof(a) / sizeof((a)[0]))

#define __SLICE_LOW(expr, typ, low)       \
	__slice __s = expr;                   \
	uint __low = low;                     \
	__s.array = ((typ)__s.array) + __low
	
#define __SLICEL(expr, typ, low) \
({                               \
	__SLICE_LOW(expr, typ, low); \
	__s;                         \
})

#define __SLICELH(expr, typ, low, high) \
({                                      \
	__SLICE_LOW(expr, typ, low);        \
	__s.len = high - __low;             \
	__s;                                \
})

#define __SLICELHM(expr, typ, low, high, max) \
({                                            \
	__SLICE_LOW(expr, typ, low);              \
	__s.len = high - __low;                   \
	__s.cap = max - __low;                    \
	__s;                                      \
})
	
#undef __SLICE_LOW

#define __SLICEH(expr, high) \
({                                \
	__slice __s = expr;           \
	__s.len = high;               \
	__s;                          \
})
	
#define __SLICEM(expr, max) "Go 1.2 doesn't allow [::max]"
	
#define __SLICEHM(expr, high, max) \
({                                      \
	__slice __s = expr;                 \
	__s.len = high;                     \
	__s.cap = max;                      \
	__s;                                \
})

	
#define __ASLICEL(expr, low) \
	(__slice){               \
		&(expr)[low],        \
		__ALEN(expr)-low,    \
		__ALEN(expr)-low,    \
	}

#define __ASLICELH(expr, low, high) \
	(__slice){                      \
		&(expr)[low],               \
		high-low,                   \
		__ALEN(expr)-low            \
	}
	
#define __ASLICELHM(expr, low, high, max) \
	(__slice){                            \
		&(expr)[low],                     \
		high-low,                         \
		max-low                           \
	}
	
#define __ASLICEH(expr, high) (__slice){expr, high, __ALEN(expr)}
	
#define __ASLICEM(expr, max) "Go 1.2 doesn't allow [::max]"
	
#define __ASLICEHM(expr, high, max) (__slice){expr, high, max}
	
#define __SLICPY(typ, dst, src)                               \
	runtime_Copy(                                             \
		dst.array, src.array,                                 \
		(dst.len < src.len ? dst.len : src.len) * sizeof(typ) \
	)
	
//#define __STRCPY(dst, src)

static inline void panic(string s) {
	for (;;) {
	}
}
