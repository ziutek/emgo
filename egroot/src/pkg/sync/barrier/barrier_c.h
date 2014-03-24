__attribute__ ((always_inline))
extern inline void sync_barrier_Compiler() {
	asm volatile ("":::"memory");
}

__attribute__ ((always_inline))
extern inline void sync_barrier_Memory() {
	__sync_synchronize();
}
