typedef union {
	unsafe$Pointer ptr;
	complex128 c128;
} ival;

typedef struct {
	ival val$;
	const void *itab$;
} interface;

#define EQUALI(lhs, rhs) ({                           \
	typeof(lhs) a = (lhs);                            \
	typeof(rhs) b = (rhs);                            \
	a.itab$ == b.itab$ && a.val$.c128 == b.val$.c128; \
})

#define CAST(t, e) ({                               \
	union {typeof(e) in; typeof(t) out;} cast = {}; \
	cast.in = (e);                                  \
	cast.out;                                       \
})

#define INTERFACE(e, itab) ((interface){CAST(ival, e), itab})
#define IVAL(i, typ) CAST(typ, (i).val$)
#define NILI ((interface){})
