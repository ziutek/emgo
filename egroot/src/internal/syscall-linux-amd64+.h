// +build linux
// +build amd64

static inline uintptr
internal$Syscall1(uintptr trap, uintptr a1) {
	register uintptr rax asm("%rax") = trap;
	register uintptr rdi asm("%rdi") = a1;
	asm volatile (
		"syscall"
		:"+r" (rax)
		:"r"(rdi)
		:"%rcx", "%r11", "memory"
	);
	return rax;
}

static inline uintptr
internal$Syscall2(uintptr trap, uintptr a1, uintptr a2) {
	register uintptr rax asm("%rax") = trap;
	register uintptr rdi asm("%rdi") = a1;
	register uintptr rsi asm("%rsi") = a2;
	asm volatile (
		"syscall"
		:"+r" (rax)
		:"r"(rdi), "r"(rsi)
		:"%rcx", "%r11", "memory"
	);
	return rax;
}

static inline uintptr
internal$Syscall3(uintptr trap, uintptr a1, uintptr a2, uintptr a3) {
	uintptr t0 = trap;
	uintptr t1 = a1;
	uintptr t2 = a2;
	uintptr t3 = a3;
	register uintptr rax asm("%rax") = t0;
	register uintptr rdi asm("%rdi") = t1;
	register uintptr rsi asm("%rsi") = t2;
	register uintptr rdx asm("%rdx") = t3;
	asm volatile (
		"syscall"
		:"+r" (rax)
		:"r"(rdi), "r"(rsi), "r"(rdx)
		:"%rcx", "%r11", "memory"
	);
	return rax;
}

static inline uintptr
internal$Syscall4(uintptr trap, uintptr a1, uintptr a2, uintptr a3, uintptr a4) {
	register uintptr rax asm("%rax") = trap;
	register uintptr rdi asm("%rdi") = a1;
	register uintptr rsi asm("%rsi") = a2;
	register uintptr rdx asm("%rdx") = a3;
	register uintptr r10 asm("%r10") = a4;
	asm volatile (
		"syscall"
		:"+r" (rax)
		:"r"(rdi), "r"(rsi), "r"(rdx), "r"(r10)
		:"%rcx", "%r11", "memory"
	);
	return rax;
}

static inline uintptr
internal$Syscall5(uintptr trap, uintptr a1, uintptr a2, uintptr a3, uintptr a4,
	uintptr a5) {
	register uintptr rax asm("%rax") = trap;
	register uintptr rdi asm("%rdi") = a1;
	register uintptr rsi asm("%rsi") = a2;
	register uintptr rdx asm("%rdx") = a3;
	register uintptr r10 asm("%r10") = a4;
	register uintptr r8 asm("%r8") = a5;
	asm volatile (
		"syscall"
		:"+r" (rax)
		:"r"(rdi), "r"(rsi), "r"(rdx), "r"(r10), "r"(r8)
		:"%rcx", "%r11", "memory"
	);
	return rax;
}

static inline uintptr
internal$Syscall6(uintptr trap, uintptr a1, uintptr a2, uintptr a3, uintptr a4,
	uintptr a5, uintptr a6) {
	register uintptr rax asm("%rax") = trap;
	register uintptr rdi asm("%rdi") = a1;
	register uintptr rsi asm("%rsi") = a2;
	register uintptr rdx asm("%rdx") = a3;
	register uintptr r10 asm("%r10") = a4;
	register uintptr r8 asm("%r8") = a5;
	register uintptr r9 asm("%r9") = a6;
	asm volatile (
		"syscall"
		:"+r" (rax)
		:"r"(rdi), "r"(rsi), "r"(rdx), "r"(r10), "r"(r8), "r"(r9)
		:"%rcx", "%r11", "memory"
	);
	return rax;
}
