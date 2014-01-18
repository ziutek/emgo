void delay_Loop(int i) {
	while (i > 0) {
		asm volatile ("nop");
		--i;
	}
}
