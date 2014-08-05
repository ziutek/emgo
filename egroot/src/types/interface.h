typedef struct {
	uintptr val;
	uint32 typ;
} interface;

#define INTERFACE(e, t) (interface){(uintptr)(e), t}

#define NILI (interface){}

typedef struct {
	interface I$;
	string (*Error)(uintptr);
} error;