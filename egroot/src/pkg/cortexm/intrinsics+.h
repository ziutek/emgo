__attribute__ ((always_inline))
extern inline uint32 cortexm$APSR() {
	uintptr r;
	asm volatile ("mrs %0, apsr" : "=r" (r));
	return r;
}

__attribute__ ((always_inline))
extern inline void cortexm$SetAPSR(uint32 r) {
	asm volatile ("msr apsr, %0" :: "r" (r) : "apsr");
}

__attribute__ ((always_inline))
extern inline uint32 cortexm$IPSR() {
	uint32 r;
	asm volatile ("mrs %0, ipsr" : "=r" (r));
	return r;
}

__attribute__ ((always_inline))
extern inline void cortexm$SetIPSR(uint32 r) {
	asm volatile ("msr ipsr, %0" :: "r" (r) : "ipsr");
}

__attribute__ ((always_inline))
extern inline uint32 cortexm$EPSR() {
	uint32 r;
	asm volatile ("mrs %0, epsr" : "=r" (r));
	return r;
}

__attribute__ ((always_inline))
extern inline uint32 cortexm$IEPSR() {
	uint32 r;
	asm volatile ("mrs %0, iepsr" : "=r" (r));
	return r;
}

__attribute__ ((always_inline))
extern inline void cortexm$SetIEPSR(uint32 r) {
	asm volatile ("msr iepsr, %0" :: "r" (r) : "iepsr");
}

__attribute__ ((always_inline))
extern inline uint32 cortexm$IAPSR() {
	uint32 r;
	asm volatile ("mrs %0, iapsr" : "=r" (r));
	return r;
}

__attribute__ ((always_inline))
extern inline void cortexm$SetIAPSR(uint32 r) {
	asm volatile ("msr iapsr, %0" :: "r" (r) : "apsr");
}

__attribute__ ((always_inline))
extern inline uint32 cortexm$EAPSR() {
	uint32 r;
	asm volatile ("mrs %0, eapsr" : "=r" (r));
	return r;
}

__attribute__ ((always_inline))
extern inline void cortexm$SetEAPSR(uint32 r) {
	asm volatile ("msr eapsr, %0" :: "r" (r) : "eapsr");
}

__attribute__ ((always_inline))
extern inline uint32 cortexm$PSR() {
	uint32 r;
	asm volatile ("mrs %0, psr" : "=r" (r));
	return r;
}

__attribute__ ((always_inline))
extern inline void cortexm$SetPSR(uint32 r) {
	asm volatile ("msr psr, %0" :: "r" (r) : "psr");
}

__attribute__ ((always_inline))
extern inline uintptr cortexm$MSP() {
	uintptr p;
	asm volatile ("mrs %0, msp" : "=r" (p));
	return p;
}

__attribute__ ((always_inline))
extern inline void cortexm$SetMSP(unsafe$Pointer p) {
	asm volatile ("msr msp, %0" :: "r" (p) : "sp");
}

__attribute__ ((always_inline))
extern inline uintptr cortexm$PSP() {
	uintptr p;
	asm volatile ("mrs %0, psp" : "=r" (p));
	return p;
}

__attribute__ ((always_inline))
extern inline void cortexm$SetPSP(unsafe$Pointer p) {
	asm volatile ("msr psp, %0" :: "r" (p) : "sp");
}

__attribute__ ((always_inline))
extern inline cortexm$Control cortexm$Ctrl() {
	cortexm$Control c;
	asm volatile ("mrs %0, control" : "=r" (c));
	return c;
}

__attribute__ ((always_inline))
extern inline void cortexm$SetCtrl(cortexm$Control c) {
	asm volatile ("msr control, %0" :: "r" (c));
}

__attribute__ ((always_inline))
extern inline void cortexm$SEV() {
	asm volatile ("sev");
}

__attribute__ ((always_inline))
extern inline void cortexm$ISB() {
	asm volatile ("isb");
}

#define cortexm$SVC(imm) asm volatile ("svc %0" :: "i" (imm))

#define cortexm$BKPT(imm) asm volatile ("bkpt %0" :: "i" (imm))