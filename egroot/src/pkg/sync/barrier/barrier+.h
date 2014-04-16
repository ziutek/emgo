__attribute__ ((always_inline)) extern inline 
void sync$barrier$Compiler() {
	asm volatile ("":::"memory");
}

__attribute__ ((always_inline)) extern inline 
void sync$barrier$Memory() {
	__sync_synchronize();
}
