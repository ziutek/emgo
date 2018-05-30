#if __SIZEOF_POINTER__ == 8

// On 64-bit architecture ival must store slice (3 x 64-bit).
typedef struct { uintptr ptr, w1, w2; } ival;

#define EQUALI(lhs, rhs) ({   \
	typeof(lhs) a = (lhs);    \
	typeof(rhs) b = (rhs);    \
	a.itab    == b.itab &&    \
	a.val.ptr == b.val.ptr && \
	a.val.w1  == b.val.w1  && \
	a.val.w2  == b.val.w2;    \
})

#else

// On 32-bit architecture ival must store complex128.
typedef struct { uintptr ptr, w; uint64 dw; } ival;

#define EQUALI(lhs, rhs) ({   \
	typeof(lhs) a = (lhs);    \
	typeof(rhs) b = (rhs);    \
	a.itab    == b.itab &&    \
	a.val.ptr == b.val.ptr && \
	a.val.w   == b.val.w   && \
	a.val.dw  == b.val.dw;    \
})

#endif

typedef struct {
	ival val;
	const void *itab;
} interface;

/*
#define INTERFACE(e, itab) ({         \
	interface _i = { .itab = itab }; \
	*((typeof(e) *)(&_i.val)) = (e); \
	_i;                              \
})
*/

#define INTERFACE(e, itab) ({    \
	ival _v = {};                \
	*((typeof(e) *)(&_v)) = (e); \
	(interface){_v, itab};       \
})

#define IVAL(i, typ) ({ interface _i = (i); *(typ *)(&_i.val); })

#define NILI ((interface){})
#define ISNILI(e) ((e).itab == nil)
