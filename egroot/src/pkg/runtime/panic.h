__attribute__ ((noreturn))
static inline void panic(string s) {
	runtime_Panic(s);
	for (;;) {
	}
}