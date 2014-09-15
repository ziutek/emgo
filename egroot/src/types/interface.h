/*typedef union {
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

typedef union {
	unsafe$Pointer ptr;
	complex128     c128;
} ival;

typedef struct {
	ival    val$;
	uintptr typ$;
} interface;

#define EQUALI(lhs, rhs) ({                         \
	typeof(lhs) a = (lhs);                          \
	typeof(rhs) b = (rhs);                          \
	a.typ$ == b.typ$ && a.val$.c128 == b.val$.c128; \
})


#define INTERFACE(e, tid) ({              \
	union {typeof(e) in; ival out;} cast; \
	cast.in = (e);                        \
	(interface){cast.out, tid};           \
})

#define NILI (interface){}

typedef struct {
	interface;
	string(*Error) (ival *);
} error;
