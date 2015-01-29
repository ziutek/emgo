static
void (*runtime$noos$p2f(uintptr p))() {
	union {uintptr in; void (*out)();} cast;
	cast.in = p;
	return cast.out;
}