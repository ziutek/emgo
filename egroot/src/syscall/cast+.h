uintptr syscall$b2u(bool b) {
	return (uintptr)(b);
}

uintptr syscall$f2u(void (*f)()) {
	return CAST(uintptr, f);
}

