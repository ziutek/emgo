// panic.h
__attribute__ ((noreturn))
static inline void panic(string s) {
	runtime_Panic(s);
	for (;;) {
	}
}