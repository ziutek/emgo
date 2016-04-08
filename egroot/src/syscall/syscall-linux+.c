// +build linux

static uintptr
syscall$syscallPath(uintptr trap, string path, uintptr a2, uintptr a3) {
	byte p[path.len + 1];
	memcpy(p, path.str, path.len);
	p[path.len] = 0;
	return internal$Syscall3(trap, (uintptr) (p), a2, a3);
}
