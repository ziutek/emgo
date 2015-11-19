void
sync$fence$Compiler() {
	asm volatile ("":::"memory");
}

void
sync$fence$Memory() {
	__sync_synchronize();
}
