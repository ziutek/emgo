void
sync$Fence() {
	asm volatile ("":::"memory");
}