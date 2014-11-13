// +build noos
// +build cortexm0 cortexm3 cortexm4 cortexm4f

#define builtin$Syscall0(trap) ({  \
	register uintptr r0 asm("r0"); \
	register uintptr r1 asm("r1"); \
	asm volatile (                 \
		"svc %2"                   \
		: "=r" (r0), "=r" (r1)     \
		: "I" (trap)               \
		: "memory"                 \
	);                             \
	(uintptr$$uintptr){r0, r1};    \
})

#define builtin$Syscall1(trap, a1) ({   \
	register uintptr r0 asm("r0") = a1; \
	register uintptr r1 asm("r1");      \
	asm volatile (                      \
		"svc %2"                        \
		: "+r" (r0), "=r" (r1)          \
		: "I" (trap)                    \
		: "memory"                      \
	);                                  \
	(uintptr$$uintptr){r0, r1};         \
})

#define builtin$Syscall2(trap, a1, a2) ({ \
	register uintptr r0 asm("r0") = a1;   \
	register uintptr r1 asm("r1") = a2;   \
	asm volatile (                        \
		"svc %2"                          \
		: "+r" (r0), "+r" (r1)            \
		: "I" (trap)                      \
		: "memory"                        \
	);                                    \
	(uintptr$$uintptr){r0, r1};           \
})


/*
// Version that uses r0 for syscall number.

__attribute__ ((always_inline))
extern inline
uintptr$$uintptr builtin$Syscall0(uintptr trap) {
	register uintptr r0 asm("r0") = trap;
	register uintptr r1 asm("r1");
	asm volatile (
		"svc 0"
		: "+r" (r0), "=r" (r1)
		:
		: "memory"
	);
	return (uintptr$$uintptr){r0, r1};
}

__attribute__ ((always_inline))
extern inline
uintptr$$uintptr builtin$Syscall1(uintptr trap, uintptr a1) {
	register uintptr r0 asm("r0") = trap;
	register uintptr r1 asm("r1") = a1;
	asm volatile (
		"svc 1"
		: "+r" (r0), "+r" (r1)
		:
		: "memory"
	);
	return (uintptr$$uintptr){r0, r1};
}

__attribute__ ((always_inline))
extern inline
uintptr$$uintptr builtin$Syscall2(uintptr trap, uintptr a1, uintptr a2) {
	register uintptr r0 asm("r0") = trap;
	register uintptr r1 asm("r1") = a1;
	register uintptr r2 asm("r2") = a2;
	asm volatile (
		"svc 2"
		: "+r" (r0), "+r" (r1)
		: "r" (r2)
		: "memory"
	);
	return (uintptr$$uintptr){r0, r1};
}
*/