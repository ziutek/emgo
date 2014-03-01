__attribute__ ((always_inline))
extern inline void sync_Barrier() {
	asm volatile ("":::"memory");
}

__attribute__ ((always_inline))
extern inline void sync_Memory() {
	__sync_synchronize();
}
