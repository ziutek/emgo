type union {
	uintptr  ptr;
	int      i;
	int64    i64;
	string   str;
	slice    sli;
	float32  f32;
	float64  f64;
	complex64 c64
} ival;

typedef struct {
	ival    val$;
	uintptr typ$;
} interface;

#define NILI (interface){}

#define EQUALI(lhs, rhs) ({               \
	typeof(lhs) a = lhs;                  \
	typeof(rhs) b = rhs;                  \
	a.typ$ == b.typ$ && a.val$ == b.val$; \
})

typedef struct {
	interface;
	string (*Error)(uintptr);
} error;