void memmove(unsafe$Pointer dst, unsafe$Pointer src, uint n);
void memcpy(unsafe$Pointer dst, unsafe$Pointer src, uint n);
void memset(unsafe$Pointer s, byte b, uint n);
bool equals(string s1, string s2);

__attribute__ ((noreturn))
void panic(interface i);

__attribute__ ((noreturn))
void panicIC();

#define NEW(typ) (typ*) internal$Alloc(1, sizeof(typ), __alignof__(typ))

#define MAKESLI(typ, l) (slice){							\
	internal$Alloc(l, sizeof(typ), __alignof__(typ)), l, l	\
}

#define MAKESLIC(typ, l, c) (slice){						\
	internal$Alloc(c, sizeof(typ), __alignof__(typ)), l, c	\
}