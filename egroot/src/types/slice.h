typedef struct {
	unsafe$Pointer arr;
	uint len;
	uint cap;
} slice;

#define _SLICE_LOW(expr, typ, low) \
	slice s = expr;               \
	uint l = low;                 \
	s.arr = ((typ)s.arr) + l
	
#define SLICEL(expr, typ, low) ({ \
	_SLICE_LOW(expr, typ, low);    \
	s.len -= l;                   \
	s.cap -= l;                   \
	s;                            \
})

#define SLICELH(expr, typ, low, high) ({ \
	_SLICE_LOW(expr, typ, low);           \
	s.len = high - l;                    \
	s.cap -= l;                          \
	s;                                   \
})

#define SLICELHM(expr, typ, low, high, max) ({ \
	_SLICE_LOW(expr, typ, low);                 \
	s.len = high - l;                          \
	s.cap = max - l;                           \
	s;                                         \
})

#define __SLICEH(expr, high) ({ \
	slice __s = expr;           \
	s.len = high;               \
	s;                          \
})
	
// #define SLICEM(expr, max) Go 1.2 doesn't allow [::max].
	
#define SLICEHM(expr, high, max) ({ \
	slice s = expr;                 \
	s.len = high;                   \
	s.cap = max;                    \
	s;                              \
})

	
#define ASLICEL(expr, low) \
	(slice){               \
		&(expr)[low],      \
		ALEN(expr)-low,    \
		ALEN(expr)-low,    \
	}

#define ASLICELH(expr, low, high) \
	(slice){                      \
		&(expr)[low],             \
		high-low,                 \
		ALEN(expr)-low            \
	}
	
#define ASLICELHM(expr, low, high, max) \
	(slice){                            \
		&(expr)[low],                   \
		high-low,                       \
		max-low                         \
	}
	
#define ASLICE(expr) (slice){expr, ALEN(expr), ALEN(expr)}

#define ASLICEH(expr, high) (slice){expr, high, ALEN(expr)}
	
// #define ASLICEM(expr, max) Go 1.2 doesn't allow [::max].
	
#define ASLICEHM(expr, high, max) (slice){expr, high, max}

#define SLICPY(typ, dst, src) ({                     \
	int n = (dst.len < src.len) ? dst.len : src.len; \
	memmove(dst.arr, src.arr, n * sizeof(typ));      \
	n;                                               \
})

#define NILSLICE (slice){0, 0, 0}