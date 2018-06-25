
inline __attribute__((always_inline))
uintptr
stm32$hal$sdmmc$burstCopyPTM(uintptr p, uintptr m) {
	asm volatile (
		"ldm %1,  {r0-r7}\n\t"
		"stm %0!, {r0-r7}"
		: "+&r" (m) 
		: "r" (p)
		: "r0", "r1", "r2", "r3", "r4", "r5", "r6", "r7", "memory"
	);
	return m;
}

inline __attribute__((always_inline))
uintptr
stm32$hal$sdmmc$burstCopyMTP(uintptr m, uintptr p) {
	asm volatile (
		"ldm %0!, {r0-r7}\n\t"
		"stm %1,  {r0-r7}"
		: "+&r" (m)
		: "r" (p) 
		: "r0", "r1", "r2", "r3", "r4", "r5", "r6", "r7", "memory"
	);
	return m;
}
