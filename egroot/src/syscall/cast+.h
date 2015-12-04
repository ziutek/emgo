uintptr syscall$btou(bool b) {
	return (uintptr)(b);
}

uintptr syscall$ftou(void (*f)()) {
	return CAST(uintptr, f);
}

uintptr syscall$f32tou(uint32 (*f)()) {
	return CAST(uintptr, f);
}