void memmove(unsafe$Pointer dst, unsafe$Pointer src, uintptr n);
void memcpy(unsafe$Pointer dst, unsafe$Pointer src, uintptr n);
void memset(unsafe$Pointer s, byte b, uintptr n);

__attribute__ ((noreturn))
void panicIC();

#define NEW(typ) ((typ*)internal$Alloc(1, sizeof(typ), __alignof__(typ)))

#define MAKESLI(typ, lx) ({						                     \
	uintptr l = lx;                                                  \
	(slice){internal$Alloc(l, sizeof(typ), __alignof__(typ)), l, l}; \
})

#define MAKESLIC(typ, lx, cx) ({                                     \
	uintptr l = lx;                                                  \
	uintptr c = cx;		                                             \
	(slice){internal$Alloc(c, sizeof(typ), __alignof__(typ)), l, c}; \
})

#define NEWSTR(bx) ({                                        \
	slice b = bx;                                            \
	string s = (string){internal$Alloc(b.len, 1, 1), b.len}; \
	internal$Memmove(s.str, b.arr, b.len);                   \
	s;                                                       \
})

#define EQUALA(a1, a2) \
	(internal$Memcmp((a1).arr, (a2).arr, sizeof((a1).arr)) == 0)

inline __attribute__((always_inline))
int_
cmpstr(string s1, string s2) {
	int_ n = len(s1);
	if (n > len(s2)) {
		n = len(s2);
	}
	int_ ret = internal$Memcmp(s1.str, s2.str, n);
	if (ret != 0 || len(s1) == len(s2)) {
		return ret;
	}
	if (len(s1) < len(s2)) {
		return -1;
	}
	return 1;
}