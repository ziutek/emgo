#define ithead builtin$ItHead
#define tinfo  builtin$Type

typedef struct {
	ithead h$;
	string(*Error) (ival *);
} error;

enum {
	Invalid = 0,
	Bool,
	Int,
	Int8,
	Int16,
	Int32,
	Int64,
	Uint,
	Uint8,
	Uint16,
	Uint32,
	Uint64,
	Uintptr,
	Float32,
	Float64,
	Complex64,
	Complex128,
	Array,
	Chan,
	Func,
	Interface,
	Map,
	Ptr,
	Slice,
	String,
	Struct,
	UnsafePointer,
};

#define TINFO(i) (((const ithead*)(i).itab$)->Type)

#define IASSIGN(expr, etyp, ityp) INTERFACE(          \
	expr,                                             \
	builtin$GetItable((void*)&ityp, (void*)&etyp).arr \
)

#define ICONVERTI(iexpr, ityp) ({                            \
	interface e = iexpr;                                     \
	(interface){                                             \
		e.val$,                                              \
		builtin$GetItable((void*)&ityp, (void*)TINFO(e)).arr \
	);                                                       \
})

#define ICONVERTE(iexpr) ({        \
	interface e = iexpr;           \
	(interface){e.val$, TINFO(e)}; \
})