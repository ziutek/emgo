typedef struct {
	unsafe$Pointer arr;
	uintptr len;
	uintptr cap;
} slice;


#define SLICEL(slx, typ, lowx) ({ \
	slice s = slx;                \
	uintptr l = lowx;             \
	s.arr = ((typ)s.arr) + l;     \
	s.len -= l;                   \
	s.cap -= l;                   \
	s;                            \
})

#define SLICELC(slx, typ, lowx) ({ \
	slice s = slx;                 \
	uintptr l = lowx;              \
	if (l > s.len) panicIndex();   \
	s.arr = ((typ)s.arr) + l;      \
	s.len -= l;                    \
	s.cap -= l;                    \
	s;                             \
})

#define SLICELH(slx, typ, lowx, highx) ({ \
	slice s = slx;                        \
	uintptr l = lowx;                     \
	s.arr = ((typ)s.arr) + l;             \
	s.len = (highx) - l;                  \
	s.cap -= l;                           \
	s;                                    \
})

#define SLICELHC(slx, typ, lowx, highx) ({ \
	slice s = slx;                         \
	uintptr l = lowx;                      \
	uintptr h = highx;                     \
	if (l > h || h > s.cap) panicIndex();  \
	s.arr = ((typ)s.arr) + l;              \
	s.len = h - l;                         \
	s.cap -= l;                            \
	s;                                     \
})

#define SLICELHM(slx, typ, lowx, highx, maxx) ({            \
	uintptr l = lowx;                                       \
	(slice){((typ)(slx).arr) + l, (highx) - l, (maxx) - l}; \
})

#define SLICELHMC(slx, typ, lowx, highx, maxx) ({  \
	slice s = slx;                                 \
	uintptr l = lowx;                              \
	uintptr h = highx;                             \
	uintptr m = maxx;                              \
	if (l > h || h > m || m > s.cap) panicIndex(); \
	s.arr = ((typ)s.arr) + l;                      \
	s.len = h - l;                                 \
	s.cap = m - l;                                 \
	s;                                             \
})

#define SLICEH(slx, highx) ({ \
	slice s = slx;            \
	s.len = highx;            \
	s;                        \
})

#define SLICEHC(slx, highx) ({   \
	slice s = slx;               \
	uintptr h = highx            \
	if (h > s.cap) panicIndex(); \
	s.len = h;                   \
	s;                           \
})

// #define SLICEM(slx, maxx) Go 1.2 doesn't allow [::max].
	
#define SLICEHM(slx, highx, maxx) ((slice){(slx).arr, highx, maxx})             

#define SLICEHMC(slx, highx, maxx) ({     \
	slice s = slx;                        \
	uintptr h = highx;                    \
	uintptr m = maxx;                     \
	if (h > m || m > s.cap) panicIndex(); \
	s.len = h;                            \
	s.cap = m;                            \
	s;                                    \
})

#define _ALEN(arx) (sizeof(arx) / sizeof((arx).arr[0]))

#define ASLICEL(arx, lowx) ({           \
	uintptr l = lowx;                   \
	uintptr newl = _ALEN(arx) - l;      \
	(slice){(arx).arr + l, newl, newl}; \
})

#define ASLICELC(arx, lowx) ({          \
	uintptr l = lowx;                   \
	if (l > _ALEN(arx)) panicIndex();   \
	uintptr newl = _ALEN(arx) - l;      \
	(slice){(arx).arr + l, newl, newl}; \
})

#define ASLICELH(arx, lowx, highx) ({                    \
	uintptr l = lowx;                                    \
	(slice){(arx).arr + l, (highx) - l, _ALEN(arx) - l}; \
})

#define ASLICELHC(arx, lowx, highx) ({              \
	uintptr l = lowx;                               \
	uintptr h = highx;                              \
	if (l > h || h > _ALEN(arx)) panicIndex();      \
	(slice){(arx).arr + l, h - l, _ALEN(arx) - l};  \
})

#define ASLICELHM(arx, lowx, highx, maxx) ({         \
	uintptr l = lowx;                                \
	(slice){(arx).arr + l, (highx) - l, (maxx) - l}; \
})
	
#define ASLICELHMC(arx, lowx, highx, maxx) ({           \
	uintptr l = lowx;                                   \
	uintptr h = highx;                                  \
	uintptr m = maxx;                                   \
	if (l > h || h > m || m > _ALEN(arx)) panicIndex(); \
	(slice){(arx).arr + l, h - l, m - l};               \
})
	
#define ASLICE(arx) ((slice){(arx).arr, _ALEN(arx), _ALEN(arx)})

#define CSLICE(len, ptrx) ((slice){(ptrx), len, len})

#define ASLICEH(arx, highx) ((slice){(arx).arr, highx, _ALEN(arx)})

#define ASLICEHC(arx, highx) ({        \
	uintptr h = highx;                 \
	if (h > _ALEN(arx)) panicIndex();  \
	(slice){(arx).arr, h, _ALEN(arx)}; \
})

// #define ASLICEM(arx, max) Go 1.2 doesn't allow [::max].
	
#define ASLICEHM(arx, highx, maxx) ((slice){(arx).arr, highx, maxx})

#define ASLICEHMC(arx, highx, maxx) ({         \
	uintptr h = highx;                         \
	uintptr m = maxx;                          \
	if (h > m || m > _ALEN(arx)) panicIndex(); \
	(slice){(arx).arr, h, m};                  \
})

#define SLICPY(typ, dstx, srcx) ({            \
	slice d = dstx;                           \
	slice s = srcx;                           \
	uint n = (d.len < s.len) ? d.len : s.len; \
	memmove(d.arr, s.arr, n * sizeof(typ));   \
	n;                                        \
})

#define NILSLICE ((slice){})

#define SLIDX(typ, slx, idx) (((typ)(slx).arr)[idx])

#define SLIDXC(typ, slx, idx)  ({ \
	slice s = slx;                \
	uintptr i = idx;              \
	if (i >= s.len) panicIndex(); \
	(typ)s.arr + i;               \
})[0]

#define AIDX(arx, idx) ((arx).arr[idx])

#define AIDXC(arx, idx) ({             \
	uintptr i = idx;                   \
	if (i >= _ALEN(arx)) panicIndex(); \
	(arx).arr + i;                     \
})[0]
