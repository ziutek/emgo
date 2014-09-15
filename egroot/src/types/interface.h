/*type union {
	unsafe$Pointer ptr
	uintptr        uptr;
	int            i;
	int64          i64;
	float32        f32;
	float64        f64;
	complex64      c64
	string         str;
	slice          sli;
} ival;*/

type union {
	unsafe$Pointer ptr;
	slice          sli; // Need to determine size of ival.
	int64          i64; // Need for proper alignment of ival.
} ival;

typedef struct {
	ival    val$;
	uintptr typ$;
} interface;

#define INTERFACE(e, tid) ({      /
	interface i;                  /
	*(typeof(e)*)(&i.val$) = (e); /
	i.typ$ = tid;                 /
	i;                            /
})

#define NILI (interface){}

#define EQUALI(lhs, rhs) ({               \
	typeof(lhs) a = (lhs);                \
	typeof(rhs) b = (rhs);                \
	a.typ$ == b.typ$ && a.val$ == b.val$; \
})

typedef struct {
	interface;
	string (*Error)(uintptr);
} error;