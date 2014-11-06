// +build noos
// +build cortexm0 cortexm3 cortexm4 cortexm4f

__attribute__ ((always_inline))
extern inline
uintptr$$syscall$Errno syscall$Syscall0(uintptr trap) {
	register uintptr r0 asm("r0") = trap;
	register uintptr r1 asm("r1");
	asm volatile ("svc 0" : "=r" (r0), "=r" (r1) : "r" (r0) : "memory");
	return (uintptr$$syscall$Errno){r0, r1};
}

