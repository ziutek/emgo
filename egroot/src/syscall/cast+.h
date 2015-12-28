uintptr
syscall$ftou(void (*f) ()) {
	return CAST(uintptr, f);
}

uintptr
syscall$f64tou(void (*f) (int64)) {
	return CAST(uintptr, f);
}

uintptr
syscall$fr64tou(int64(*f) ()) {
	return CAST(uintptr, f);
}
