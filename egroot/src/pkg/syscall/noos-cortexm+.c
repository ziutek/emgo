// +build noos
// +build cortexm0 cortexm3 cortexm4 cortexm4f

static inline
uintptr$$syscall$Errno syscall$syscall0(uintptr trap) {
	register uintptr r0 asm("r0") = trap;
	register uintptr r1 asm("r1");
	asm volatile (
		"svc 0"
		: "=r" (r0), "=r" (r1)
		: "0" (r0)
		: "memory"
	);
	return (uintptr$$syscall$Errno){r0, r1};
}

static inline
uintptr$$syscall$Errno syscall$syscall1(uintptr trap, uintptr a1) {
	register uintptr r0 asm("r0") = trap;
	register uintptr r1 asm("r1") = a1;
	asm volatile (
		"svc 1"
		: "=r" (r0), "=r" (r1)
		: "0" (r0), "1" (r1)
		: "memory"
	);
	return (uintptr$$syscall$Errno){r0, r1};
}

static inline
uintptr$$syscall$Errno syscall$syscall2(uintptr trap, uintptr a1, uintptr a2) {
	register uintptr r0 asm("r0") = trap;
	register uintptr r1 asm("r1") = a1;
	register uintptr r2 asm("r2") = a2;
	asm volatile (
		"svc 2"
		: "=r" (r0), "=r" (r1)
		: "0" (r0), "1" (r1), "r" (r2)
		: "memory"
	);
	return (uintptr$$syscall$Errno){r0, r1};
}
