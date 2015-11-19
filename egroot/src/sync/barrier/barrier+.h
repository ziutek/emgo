void
sync$barrier$Compiler() {
	asm volatile ("":::"memory");
}

void
sync$barrier$Memory() {
	__sync_synchronize();
}
