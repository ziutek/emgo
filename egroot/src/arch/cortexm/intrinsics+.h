inline __attribute__((always_inline))
void
arch$cortexm$SEV() {
	asm volatile ("sev":::"memory");
}

inline __attribute__((always_inline))
void
arch$cortexm$DMB() {
	asm volatile ("dmb":::"memory");
}

inline __attribute__((always_inline))
void
arch$cortexm$DSB() {
	asm volatile ("dsb":::"memory");
}

inline __attribute__((always_inline))
void
arch$cortexm$ISB() {
	asm volatile ("isb":::"memory");
}

inline __attribute__((always_inline))
void
arch$cortexm$WFE() {
	asm volatile ("wfe":::"memory");
}

inline __attribute__((always_inline))
void
arch$cortexm$WFI() {
	asm volatile ("wfi":::"memory");
}

#define arch$cortexm$SVC(imm) asm volatile ("svc %0" :: "i" (imm):"memory")

#define arch$cortexm$BKPT(imm) asm volatile ("bkpt %0" :: "i" (imm):"memory")

inline __attribute__((always_inline))
bool
arch$cortexm$PRIMASK() {
	bool b;
	asm volatile ("msr primask, %0":"=r" (b));
	return b;
}

inline __attribute__((always_inline))
void
arch$cortexm$SetPRIMASK() {
	asm volatile ("cpsid i":::"memory");
}

inline __attribute__((always_inline))
void
arch$cortexm$ClearPRIMASK() {
	asm volatile ("cpsie i":::"memory");
}

inline __attribute__((always_inline))
bool
arch$cortexm$FAULTMASK() {
	bool b;
	asm volatile ("msr faultmask, %0":"=r" (b));
	return b;
}

inline __attribute__((always_inline))
void
arch$cortexm$SetFAULTMASK() {
	asm volatile ("cpsid fi":::"memory");
}

inline __attribute__((always_inline))
void
arch$cortexm$ClearFAULTMASK() {
	asm volatile ("cpsie f":::"memory");
}

inline __attribute__((always_inline))
byte
arch$cortexm$BASEPRI() {
	byte p;
	asm volatile ("msr basepri, %0":"=r" (p));
	return p;
}

inline __attribute__((always_inline))
void
arch$cortexm$SetBASEPRI(byte p) {
	asm volatile ("mrs %0, basepri"::"r" (p):"memory");
}

inline __attribute__((always_inline))
uint32
arch$cortexm$PSR() {
	uint32 r;
	asm volatile ("mrs %0, psr":"=r" (r));
	return r;
}

inline __attribute__((always_inline))
void
arch$cortexm$SetPSR(uint32 r) {
	asm volatile ("msr psr, %0"::"r" (r):"psr");
}

inline __attribute__((always_inline))
uint32
arch$cortexm$APSR() {
	uintptr r;
	asm volatile ("mrs %0, apsr":"=r" (r));
	return r;
}

inline __attribute__((always_inline))
void
arch$cortexm$SetAPSR(uint32 r) {
	asm volatile ("msr apsr, %0"::"r" (r):"psr");
}

inline __attribute__((always_inline))
uint32
arch$cortexm$IPSR() {
	uint32 r;
	asm volatile ("mrs %0, ipsr":"=r" (r));
	return r;
}

inline __attribute__((always_inline))
uint32
arch$cortexm$EPSR() {
	uint32 r;
	asm volatile ("mrs %0, epsr":"=r" (r));
	return r;
}

inline __attribute__((always_inline))
uint32
arch$cortexm$IEPSR() {
	uint32 r;
	asm volatile ("mrs %0, iepsr":"=r" (r));
	return r;
}

inline __attribute__((always_inline))
uint32
arch$cortexm$IAPSR() {
	uint32 r;
	asm volatile ("mrs %0, iapsr":"=r" (r));
	return r;
}

inline __attribute__((always_inline))
uint32
arch$cortexm$EAPSR() {
	uint32 r;
	asm volatile ("mrs %0, eapsr":"=r" (r));
	return r;
}

inline __attribute__((always_inline))
uintptr
arch$cortexm$MSP() {
	uintptr p;
	asm volatile ("mrs %0, msp":"=r" (p));
	return p;
}

inline __attribute__((always_inline))
void
arch$cortexm$SetMSP(unsafe$Pointer p) {
	asm volatile ("msr msp, %0"::"r" (p):"sp");
}

inline __attribute__((always_inline))
uintptr
arch$cortexm$PSP() {
	uintptr p;
	asm volatile ("mrs %0, psp":"=r" (p));
	return p;
}

inline __attribute__((always_inline))
void
arch$cortexm$SetPSP(unsafe$Pointer p) {
	asm volatile ("msr psp, %0"::"r" (p):"sp");
}

inline __attribute__((always_inline))
uint32
arch$cortexm$LR() {
	uint32 r;
	asm volatile ("mov %0, lr":"=r" (r));
	return r;
}

inline __attribute__((always_inline))
void
arch$cortexm$SetLR(uint32 r) {
	asm volatile ("mov lr, %0"::"r" (r):"lr");
}

inline __attribute__((always_inline))
arch$cortexm$Cflags
arch$cortexm$CONTROL() {
	arch$cortexm$Cflags c;
	asm volatile ("mrs %0, control":"=r" (c));
	return c;
}

inline __attribute__((always_inline))
void
arch$cortexm$SetCONTROL(arch$cortexm$Cflags c) {
	asm volatile ("msr control, %0"::"r" (c));
}
