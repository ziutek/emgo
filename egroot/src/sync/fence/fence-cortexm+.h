// +build !smp
// +build cortexm0 cortexm0p cortexm3 cortexm4 cortexm4f cortexm7f cortexm7d

inline __attribute__((always_inline))
void
sync$fence$RW() {
	asm volatile ("dsb":::"memory");
}

inline __attribute__((always_inline))
void
sync$fence$RW_SMP() {
	asm volatile ("":::"memory");
}

inline __attribute__((always_inline))
void
sync$fence$RDP() {
}

inline __attribute__((always_inline))
void
sync$fence$RDP_SMP() {
}

inline __attribute__((always_inline))
void
sync$fence$R() {
	asm volatile ("dsb":::"memory");
}

inline __attribute__((always_inline))
void
sync$fence$R_SMP() {
	asm volatile ("":::"memory");
}

inline __attribute__((always_inline))
void
sync$fence$W() {
	asm volatile ("dsb":::"memory");
}

inline __attribute__((always_inline))
void
sync$fence$W_SMP() {
	asm volatile ("":::"memory");
}
