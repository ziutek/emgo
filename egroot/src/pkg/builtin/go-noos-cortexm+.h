// +build noos
// +build cortexm0 cortexm3 cortexm4 cortexm4f

#define GO(call, wait) do {                                  \
	void func() {                                            \
		call;                                                \
		asm volatile ("svc 1" ::: "memory");                 \
	}                                                        \
	register void (*r0)() asm("r0") = func;                  \
	register int r1 asm("r1") = (wait);                      \
	asm volatile ("svc 0" :: "r" (r0), "r" (r1) : "memory"); \
} while(0)

static inline
void goready() {
	asm volatile ("svc 2" ::: "memory");
}