typedef struct {
	uintptr val$;
	uintptr typ$;
} interface;

#define INTERFACE(v, t) (interface){.val$ = (uintptr)(v), .typ$ = t}

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