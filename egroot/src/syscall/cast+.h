uintptr
syscall$ftou(void (*f) ()) {
	return (uintptr)(f);
}

uintptr
syscall$f64btou(void (*f) (int64, bool)) {
	return (uintptr)(f);
}

uintptr
syscall$fr64tou(int64(*f) ()) {
	return (uintptr)(f);
}
