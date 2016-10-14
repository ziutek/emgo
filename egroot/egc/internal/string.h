typedef struct {
	byte *str;
	uintptr len;
} string;

#define EGSTR(str) {str, sizeof(str)-1}
#define EGSTL(str) ((string){str, sizeof(str)-1})

#define SSLICEL(strx, lowx) ({ \
	string s = strx;           \
	uintptr l = lowx;          \
	s.str += l;                \
	s.len -= l;                \
	s;                         \
})

#define SSLICELC(strx, lowx) ( { \
	string s = strx;             \
	uintptr l = lowx;            \
	if (l > s.len) panicIndex(); \
	s.str += l;                  \
	s.len -= l;                  \
	s;                           \
})

#define SSLICELH(strx, lowx, highx) ({     \
	uintptr l = lowx;                      \
	(string){(strx).str + l, (highx) - l}; \
})

#define SSLICELHC(strx, lowx, highx) ({   \
	string s = strx;                      \
	uintptr l = lowx;                     \
	uintptr h = highx;                    \
	if (l > h || h > s.len) panicIndex(); \
	s.str += l;                           \
	s.len = h - l;                        \
	s;                                    \
})

#define SSLICEH(strx, highx) ((string){(strx).str, highx})

#define SSLICEHC(strx, highx) ({ \
	string s = strx;             \
	uintptr h = highx;           \
	if (h > s.len) panicIndex(); \
	s.len = h;                   \
	s;                           \
})

#define STRCPY(dstx, srcx) ({                 \
	slice d = dstx;                           \
	string s = srcx;                          \
	uint n = (d.len < s.len) ? d.len : s.len; \
	memmove(d.arr, s.str, n);                 \
	n;                                        \
})

#define BYTES(strx) ({                                        \
	string s = strx;                                          \
	slice b = (slice){__builtin_alloca(s.len), s.len, s.len}; \
	__builtin_memcpy(b.arr, s.str, s.len);                    \
	b;                                                        \
})

#define STRIDX(strx, idx) ((strx).str[idx])

#define STRIDXC(strx, idx) ({     \
	string s = strx;              \
	uintptr i = idx;              \
	if (i >= s.len) panicIndex(); \
	s.str + i;                    \
})[0]
