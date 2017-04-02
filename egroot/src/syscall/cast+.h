inline __attribute__((always_inline))
uintptr
syscall$ftou(void (*f) ()) {
	return (uintptr)(f);
}

inline __attribute__((always_inline))
uintptr
syscall$f64btou(void (*f) (int64, bool)) {
	return (uintptr)(f);
}

inline __attribute__((always_inline))
uintptr
syscall$fr64tou(int64(*f) ()) {
	return (uintptr)(f);
}
