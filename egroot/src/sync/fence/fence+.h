extern inline __attribute__ ((always_inline))
void
sync$fence$Compiler() {
	asm volatile ("":::"memory");
}
