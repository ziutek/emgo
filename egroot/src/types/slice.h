typedef struct {
	unsafe_Pointer arr;
	uint len;
	uint cap;
} __slice;

#define __SLICE_LOW(expr, typ, low)       \
	__slice __s = expr;                   \
	uint __low = low;                     \
	__s.arr = ((typ)__s.arr) + __low
	
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
	
#define __ASLICE(expr) (__slice){expr, __ALEN(expr), __ALEN(expr)}

#define __ASLICEH(expr, high) (__slice){expr, high, __ALEN(expr)}
	
#define __ASLICEM(expr, max) "Go 1.2 doesn't allow [::max]"
	
#define __ASLICEHM(expr, high, max) (__slice){expr, high, max}

#define __SLICPY(typ, dst, src)                               \
	runtime_Copy(                                             \
		dst.arr, src.arr,                                     \
		(dst.len < src.len ? dst.len : src.len) * sizeof(typ) \
	)
