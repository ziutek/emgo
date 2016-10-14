typedef struct {
	unsafe$Pointer arr;
	uint len;
	uint cap;
} slice;


#define SLICEL(slx, typ, lowx) ({ \
	slice s = slx;                \
	uint l = lowx;                \
	s.arr = ((typ)s.arr) + l;     \
	s.len -= l;                   \
	s.cap -= l;                   \
	s;                            \
})

#define SLICELC(slx, typ, lowx) ({ \
	slice s = slx;                 \
	uint l = lowx;                 \
	if (l > s.len) panicIndex();   \
	s.arr = ((typ)s.arr) + l;      \
	s.len -= l;                    \
	s.cap -= l;                    \
	s;                             \
})

#define SLICELH(slx, typ, lowx, highx) ({ \
	slice s = slx;                        \
	uint l = lowx;                        \
	s.arr = ((typ)s.arr) + l;             \
	s.len = (highx) - l;                  \
	s.cap -= l;                           \
	s;                                    \
})

#define SLICELHC(slx, typ, lowx, highx) ({ \
	slice s = slx;                         \
	uint l = lowx;                         \
	uint h = highx;                        \
	if (l > h || h > s.cap) panicIndex();  \
	s.arr = ((typ)s.arr) + l;              \
	s.len = h - l;                         \
	s.cap -= l;                            \
	s;                                     \
})

#define SLICELHM(slx, typ, lowx, highx, maxx) ({            \
	uint l = lowx;                                          \
	(slice){((typ)(slx).arr) + l, (highx) - l, (maxx) - l}; \
})

#define SLICELHMC(slx, typ, lowx, highx, maxx) ({  \
	slice s = slx;                                 \
	uint l = lowx;                                 \
	uint h = highx;                                \
	uint m = maxx;                                 \
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
	uint h = highx;              \
	if (h > s.cap) panicIndex(); \
	s.len = h;                   \
	s;                           \
})

// #define SLICEM(slx, maxx) Go 1.2 doesn't allow [::max].
	
#define SLICEHM(slx, highx, maxx) ((slice){(slx).arr, highx, maxx})             

#define SLICEHMC(slx, highx, maxx) ({     \
	slice s = slx;                        \
	uint h = highx;                       \
	uint m = maxx;                        \
	if (h > m || m > s.cap) panicIndex(); \
	s.len = h;                            \
	s.cap = m;                            \
	s;                                    \
})

#define _ALEN(a) (sizeof(a->arr) / sizeof(a->arr[0]))

#define ASLICEL(arx, lowx) ({        \
	typeof(arx) a = arx;             \
	uint l = lowx;                   \
	uint newl = _ALEN(a) - l;        \
	(slice){a->arr + l, newl, newl}; \
})

#define ASLICELC(arx, lowx) ({       \
	typeof(arx) a = arx;             \
	uint l = lowx;                   \
	if (l > _ALEN(a)) panicIndex();  \
	uint newl = _ALEN(a) - l;        \
	(slice){a->arr + l, newl, newl}; \
})

#define ASLICELH(arx, lowx, highx) ({               \
	typeof(arx) a = arx;                            \
	uint l = lowx;                                  \
	(slice){a->arr + l, (highx) - l, _ALEN(a) - l}; \
})

#define ASLICELHC(arx, lowx, highx) ({         \
	typeof(arx) a = arx;                       \
	uint l = lowx;                             \
	uint h = highx;                            \
	if (l > h || h > _ALEN(a)) panicIndex();   \
	(slice){a->arr + l, h - l, _ALEN(a) - l};  \
})

#define ASLICELHM(arx, lowx, highx, maxx) ({      \
	typeof(arx) a = arx;                          \
	uint l = lowx;                                \
	(slice){a->arr + l, (highx) - l, (maxx) - l}; \
})
	
#define ASLICELHMC(arx, lowx, highx, maxx) ({         \
	typeof(arx) a = arx;                              \
	uint l = lowx;                                    \
	uint h = highx;                                   \
	uint m = maxx;                                    \
	if (l > h || h > m || m > _ALEN(a)) panicIndex(); \
	(slice){a->arr + l, h - l, m - l};                \
})
	
#define ASLICE(arx) ({                   \
	typeof(arx) a = arx;                 \
	(slice){a->arr, _ALEN(a), _ALEN(a)}; \
})

#define CSLICE(len, ptrx) ((slice){(ptrx), len, len})

#define ASLICEH(arx, highx) ({         \
	typeof(arx) a = arx;               \
	((slice){a->arr, highx, _ALEN(a)}; \
})

#define ASLICEHC(arx, highx) ({     \
	typeof(arx) a = arx;            \
	uint h = highx;                 \
	if (h > _ALEN(a)) panicIndex(); \
	(slice){a->arr, h, _ALEN(a)};   \
})
	
#define ASLICEHM(arx, highx, maxx) ((slice){(arx)->arr, highx, maxx})

#define ASLICEHMC(arx, highx, maxx) ({       \
	typeof(arx) a = arx;                     \
	uint h = highx;                          \
	uint m = maxx;                           \
	if (h > m || m > _ALEN(a)) panicIndex(); \
	(slice){a->arr, h, m};                   \
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
	uint i = idx;                 \
	if (i >= s.len) panicIndex(); \
	(typ)s.arr + i;               \
})[0]

#define AIDX(arx, idx) ((arx)->arr[idx])

#define AIDXC(arx, idx) ({           \
	typeof(arx) a = arx;             \
	uint i = idx;                    \
	if (i >= _ALEN(a)) panicIndex(); \
	a->arr + i;                      \
})[0]
