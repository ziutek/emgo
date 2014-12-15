// +build noos

__attribute__ ((always_inline))
extern inline
uintptr syscall$b2p(bool b) {
	return (uintptr)(b);
}

__attribute__ ((always_inline))
extern inline
uintptr syscall$f2p(void (*f)()) {
	union {void (*in)(); uintptr out;} cast;
	cast.in = f;
	return cast.out;
}

