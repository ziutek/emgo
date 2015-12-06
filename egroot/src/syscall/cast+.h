uintptr
syscall$btou(bool b) {
	return (uintptr) (b);
}

uintptr
syscall$ftou(void (*f) ()) {
	return CAST(uintptr, f);
}

uintptr
syscall$f64tou(void (*f) (uint64)) {
	return CAST(uintptr, f);
}

uintptr
syscall$fr64tou(uint64(*f) ()) {
	return CAST(uintptr, f);
}
