typedef struct {
	uintptr val$;
	uintptr typ$;
} interface;

#define INTERFACE(v, t) (interface){.val$ = (uintptr)(v), .typ$ = t}

typedef struct {
	interface;
	string (*Error)(uintptr);
} error;