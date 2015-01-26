typedef union {
	unsafe$Pointer ptr;
	complex128     c128;
} ival;

typedef struct {
	ival       val$;
 	const void *itab$;
} interface;

#define EQUALI(lhs, rhs) ({                           \
	typeof(lhs) a = (lhs);                            \
	typeof(rhs) b = (rhs);                            \
	a.itab$ == b.itab$ && a.val$.c128 == b.val$.c128; \
})


#define INTERFACE(e, itab) ({                  \
	union {typeof(e) in; ival out;} cast = {}; \
	cast.in = (e);                             \
	(interface){cast.out, itab};               \
})

#define NILI (interface){}

typedef struct {
	string (*Error)(ival *);
} error;
