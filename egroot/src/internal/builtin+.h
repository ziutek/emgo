void memmove(unsafe$Pointer dst, unsafe$Pointer src, uint n);
void memcpy(unsafe$Pointer dst, unsafe$Pointer src, uint n);
void memset(unsafe$Pointer s, byte b, uint n);
bool equals(string s1, string s2);

__attribute__ ((noreturn))
void panic(interface i);

__attribute__ ((noreturn))
void panicIC();

#define NEW(typ) (typ*) internal$Alloc(1, sizeof(typ), __alignof__(typ))

#define MAKESLI(typ, lx) ({						                     \
	uintptr l = lx;                                                  \
	(slice){internal$Alloc(l, sizeof(typ), __alignof__(typ)), l, l}; \
})

#define MAKESLIC(typ, lx, cx) (slice){                               \
	uintptr l = lx;                                                  \
	uintptr c = cx;		                                             \
	(slice){internal$Alloc(c, sizeof(typ), __alignof__(typ)), l, c}; \
}

#define NEWSTR(bx) ({                                        \
	slice b = bx;                                            \
	string s = (string){internal$Alloc(b.len, 1, 1), b.len}; \
	internal$Memmove(s.str, b.arr, b.len);                   \
	s;                                                       \
})