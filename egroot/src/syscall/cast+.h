uintptr
syscall$ftou(void (*f) ()) {
	return (uintptr)(f);
}

uintptr
syscall$f64tou(void (*f) (int64)) {
	return (uintptr)(f);
}

uintptr
syscall$fr64tou(int64(*f) ()) {
	return (uintptr)(f);
}
