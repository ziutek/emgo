typedef struct {
	uintptr val$;
	uint32 typ$;
} interface;

#define INTERFACE(v, t) (interface){.val$ = (uintptr)(v), .typ$ = t}

typedef struct {
	interface;
	string (*Error)(uintptr);
} error;