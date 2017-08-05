typedef union {
	unsafe$Pointer ptr;
	uint64         u128[2];
	slice          sli;

	// Additional fields, useful for debugging:
	int8       i8;
	int16      i16;
	int32      i32;
	int64      i64;
	float32    f32;
	float64    f64;
	complex64  c64;
	complex128 c128;
	string     str;
} ival;

typedef struct {
	ival val$;
	const void *itab$;
} interface;

#define EQUALI(lhs, rhs) ({                    \
	typeof(lhs) a = (lhs);                     \
	typeof(rhs) b = (rhs);                     \
	sizeof(a.val$.sli) > sizeof(a.val$.u128) ? \
		a.itab$ == b.itab$ &&                  \
		a.val$.sli.arr == b.val$.sli.arr &&    \
		a.val$.sli.len == b.val$.sli.len &&    \
		a.val$.sli.cap == b.val$.sli.cap       \
	:                                          \
		a.itab$ == b.itab$ &&                  \
		a.val$.u128[0] == b.val$.u128[0] &&    \
		a.val$.u128[1] == b.val$.u128[1]       \
	;                                          \
})


#define _CAST(t, e) ({                              \
	union {typeof(e) in; typeof(t) out;} cast = {}; \
	cast.in = (e);                                  \
	cast.out;                                       \
})
#define INTERFACE(e, itab) ((interface){_CAST(ival, e), itab})
#define IVAL(i, typ) _CAST(typ, (i).val$)

#define NILI ((interface){})
#define ISNILI(e) ((e).itab$ == nil)
