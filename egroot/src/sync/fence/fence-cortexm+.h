// +build !smp
// +build cortexm0 cortexm0p cortexm3 cortexm4 cortexm4f cortexm7f cortexm7d

void
sync$fence$RW() {
	asm volatile ("dsb":::"memory");
}

void
sync$fence$RW_SMP() {
	asm volatile ("":::"memory");
}

void
sync$fence$RD() {
}

void
sync$fence$RD_SMP() {
}

void
sync$fence$R() {
	asm volatile ("dsb":::"memory");
}

void
sync$fence$R_SMP() {
	asm volatile ("":::"memory");
}

void
sync$fence$W() {
	asm volatile ("dsb":::"memory");
}

void
sync$fence$W_SMP() {
	asm volatile ("":::"memory");
}
