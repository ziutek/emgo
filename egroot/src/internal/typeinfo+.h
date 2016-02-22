#define ithead internal$ItHead

struct minfo {
	internal$Method;
};

struct tinfo {
	internal$Type;
	unsafe$Pointer imethods[];
};

typedef struct {
	ithead h$;
	string (*Error)(ival *);
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
	Chan,
	Func,
	Interface,
	Map,
	Ptr,
	Slice,
	String,
	Struct,
	UnsafePointer,
	
	Array = -1
};

//#define TINFO(i) (tinfo*)(((const ithead*)((i).itab$))->typ)

#define TINFO(i) ({                  \
	const ithead *ith = (i).itab$;   \
	const tinfo *ti = nil;           \
	if (ith != nil) {                \
		ti = (const tinfo*)ith->typ; \
	}                                \
	ti;                              \
})
	

#define IASSIGN(expr, etyp, ityp) INTERFACE(        \
	expr,                                           \
	internal$ItableFor((tinfo*)&ityp, (tinfo*)&etyp) \
)

#define ICONVERTIE(iexpr) ({       \
	interface e = iexpr;           \
	(interface){e.val$, TINFO(e)}; \
})

#define ICONVERTEI(iexpr, ityp) ({                          \
	interface e = iexpr;                                    \
	(interface){                                            \
		e.val$,                                             \
		internal$ItableFor((tinfo*)&ityp, (tinfo*)(e.itab$)) \
	};                                                      \
})


#define ICONVERTII(iexpr, ityp) ({                         \
	interface e = iexpr;                                   \
	(interface){                                           \
		e.val$,                                            \
		internal$ItableFor((tinfo*)&ityp, (tinfo*)TINFO(e)) \
	};                                                     \
})

static inline
bool implements(const internal$Type* t, const internal$Type * it) {
	return internal$Type$Implements((internal$Type*)t, (internal$Type*)it);
}