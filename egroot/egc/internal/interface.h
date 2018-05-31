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

// Be careful when choosing a name for temporary variables in macro. See what
// gotoc uses for its temporary variables (eg. _i, _ok, _tmp).

/*
#define INTERFACE(e, itab) ({       \
	interface i = { .itab = itab }; \
	*((typeof(e) *)(&i.val)) = (e); \
	i;                              \
})
*/

#define INTERFACE(e, itab) ({   \
	ival v = {};                \
	*((typeof(e) *)(&v)) = (e); \
	(interface){v, itab};       \
})

#define IVAL(ie, typ) ({ interface i = (ie); *(typ *)(&i.val); })

#define NILI ((interface){})
#define ISNILI(ie) ((ie).itab == nil)
