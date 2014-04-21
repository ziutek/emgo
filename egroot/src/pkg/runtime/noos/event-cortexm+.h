// +build cortexm0 cortexm3 cortexm4 cortexm4f

__attribute__ ((always_inline))
extern inline
void runtime$noos$Event$Wait(runtime$noos$Event e) {
	register runtime$noos$Event r0 asm("r0") = e;
	asm volatile ("svc 3" :: "r" (r0) : "memory");
}