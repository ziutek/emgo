void
sync$fence$Compiler() {
	asm volatile ("":::"memory");
}
