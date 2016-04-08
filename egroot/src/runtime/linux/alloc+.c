extern byte _end;

static uintptr
runtime$linux$end() {
	return (uintptr) (&_end);
}
