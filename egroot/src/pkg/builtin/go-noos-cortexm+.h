// +build noos
// +build cortexm0 cortexm3 cortexm4 cortexm4f

#define _GORUN(call)                    		    \
	void func() {									\
		call;										\
		goexit();									\
	}												\
	register void (*r0)() asm("r0") = func;			\
	asm volatile ("svc 0" :: "r" (r0), "r" (r1))
	
#define GO(call) do {				\
	register int r1 asm("r1") = 0;	\
	_GORUN(call);					\
} while(0)

#define GOWAIT(call) do {			\
	register int r1 asm("r1") = 1;	\
	_GORUN(call);					\
} while(0)

static inline
void goexit() {
	asm volatile ("svc 1");
}

static inline
void goready() {
	asm volatile ("svc 2");
}
