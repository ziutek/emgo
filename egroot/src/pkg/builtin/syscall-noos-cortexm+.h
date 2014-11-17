// +build noos
// +build cortexm0 cortexm3 cortexm4 cortexm4f

#define builtin$Syscall0(trap) ({  \
	register uintptr r0 asm("r0"); \
	register uintptr r1 asm("r1"); \
	asm volatile (                 \
		"svc %2"                   \
		: "=r" (r0), "=r" (r1)     \
		: "i" (trap)               \
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
		: "i" (trap)                    \
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
		: "i" (trap)                      \
		: "memory"                        \
	);                                    \
	(uintptr$$uintptr){r0, r1};           \
})

// uint64 in register
//
// ARM EABI tells that 64bit operand is stored in even:odd register pair. But
// I'm unsure, does `register uint64 r asm("r0")` means that r ocupies r0:r1
// even if r isn't an operand nor a return value of function.
#define builtin$Syscall0u64(trap) ({ \
	register uint64 r asm("r0");     \
	asm volatile (                   \
		"svc %1"                     \
		: "=r" (r)                   \
		: "i" (trap)                 \
		: "memory"                   \
	);                               \
	r;                               \
})