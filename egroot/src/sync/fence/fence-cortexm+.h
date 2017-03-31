// +build !smp
// +build cortexm0 cortexm0p cortexm3 cortexm4 cortexm4f cortexm7f cortexm7d

extern inline __attribute__ ((always_inline))
void
sync$fence$RW() {
	asm volatile ("dsb":::"memory");
}

extern inline __attribute__ ((always_inline))
void
sync$fence$RW_SMP() {
	asm volatile ("":::"memory");
}

extern inline __attribute__ ((always_inline))
void
sync$fence$RDP() {
}

extern inline __attribute__ ((always_inline))
void
sync$fence$RDP_SMP() {
}

extern inline __attribute__ ((always_inline))
void
sync$fence$R() {
	asm volatile ("dsb":::"memory");
}

extern inline __attribute__ ((always_inline))
void
sync$fence$R_SMP() {
	asm volatile ("":::"memory");
}

extern inline __attribute__ ((always_inline))
void
sync$fence$W() {
	asm volatile ("dsb":::"memory");
}

extern inline __attribute__ ((always_inline))
void
sync$fence$W_SMP() {
	asm volatile ("":::"memory");
}
