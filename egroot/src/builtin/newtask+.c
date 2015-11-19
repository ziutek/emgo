uintptr
builtin$b2u(bool b) {
	return (uintptr) (b);
}

uintptr
builtin$f2u(void (*f) ()) {
	return CAST(uintptr, f);
}
