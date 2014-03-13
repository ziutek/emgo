__attribute__ ((always_inline))
extern inline uint32 cortexm_APSR() {
	uintptr r;
	asm volatile ("mrs %0, apsr" : "=r" (r));
	return r;
}

__attribute__ ((always_inline))
extern inline void cortexm_SetAPSR(uint32 r) {
	asm volatile ("msr apsr, %0" :: "r" (r) : "apsr");
}

__attribute__ ((always_inline))
extern inline uint32 cortexm_IPSR() {
	uint32 r;
	asm volatile ("mrs %0, ipsr" : "=r" (r));
	return r;
}

__attribute__ ((always_inline))
extern inline void cortexm_SetIPSR(uint32 r) {
	asm volatile ("msr ipsr, %0" :: "r" (r) : "ipsr");
}

__attribute__ ((always_inline))
extern inline uint32 cortexm_EPSR() {
	uint32 r;
	asm volatile ("mrs %0, epsr" : "=r" (r));
	return r;
}

__attribute__ ((always_inline))
extern inline uint32 cortexm_IEPSR() {
	uint32 r;
	asm volatile ("mrs %0, iepsr" : "=r" (r));
	return r;
}

__attribute__ ((always_inline))
extern inline void cortexm_SetIEPSR(uint32 r) {
	asm volatile ("msr iepsr, %0" :: "r" (r) : "iepsr");
}

__attribute__ ((always_inline))
extern inline uint32 cortexm_IAPSR() {
	uint32 r;
	asm volatile ("mrs %0, iapsr" : "=r" (r));
	return r;
}

__attribute__ ((always_inline))
extern inline void cortexm_SetIAPSR(uint32 r) {
	asm volatile ("msr iapsr, %0" :: "r" (r) : "apsr");
}

__attribute__ ((always_inline))
extern inline uint32 cortexm_EAPSR() {
	uint32 r;
	asm volatile ("mrs %0, eapsr" : "=r" (r));
	return r;
}

__attribute__ ((always_inline))
extern inline void cortexm_SetEAPSR(uint32 r) {
	asm volatile ("msr eapsr, %0" :: "r" (r) : "eapsr");
}

__attribute__ ((always_inline))
extern inline uint32 cortexm_PSR() {
	uint32 r;
	asm volatile ("mrs %0, psr" : "=r" (r));
	return r;
}

__attribute__ ((always_inline))
extern inline void cortexm_SetPSR(uint32 r) {
	asm volatile ("msr psr, %0" :: "r" (r) : "psr");
}

__attribute__ ((always_inline))
extern inline uintptr cortexm_MSP() {
	uintptr p;
	asm volatile ("mrs %0, msp" : "=r" (p));
	return p;
}

__attribute__ ((always_inline))
extern inline void cortexm_SetMSP(unsafe_Pointer p) {
	asm volatile ("msr msp, %0" :: "r" (p) : "sp");
}

__attribute__ ((always_inline))
extern inline uintptr cortexm_PSP() {
	uintptr p;
	asm volatile ("mrs %0, psp" : "=r" (p));
	return p;
}

__attribute__ ((always_inline))
extern inline void cortexm_SetPSP(unsafe_Pointer p) {
	asm volatile ("msr psp, %0" :: "r" (p) : "sp");
}

__attribute__ ((always_inline))
extern inline cortexm_Control cortexm_Ctrl() {
	cortexm_Control c;
	asm volatile ("mrs %0, control" : "=r" (c));
	return c;
}

__attribute__ ((always_inline))
extern inline void cortexm_SetCtrl(cortexm_Control c) {
	asm volatile ("msr control, %0" :: "r" (c));
}

__attribute__ ((always_inline))
extern inline void cortexm_SEV() {
	asm volatile ("sev");
}

__attribute__ ((always_inline))
extern inline void cortexm_SVC(byte n) {
	asm volatile ("svc %0" :: "i" (n));
}

__attribute__ ((always_inline))
extern inline void cortexm_ISB() {
	asm volatile ("isb");
}