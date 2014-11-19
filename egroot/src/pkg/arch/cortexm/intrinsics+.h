__attribute__ ((always_inline))
extern inline
uint32 arch$cortexm$APSR() {
	uintptr r;
	asm volatile ("mrs %0, apsr" : "=r" (r));
	return r;
}

__attribute__ ((always_inline))
extern inline
void arch$cortexm$SetAPSR(uint32 r) {
	asm volatile ("msr apsr, %0" :: "r" (r) : "apsr");
}

__attribute__ ((always_inline))
extern inline
uint32 arch$cortexm$IPSR() {
	uint32 r;
	asm volatile ("mrs %0, ipsr" : "=r" (r));
	return r;
}

__attribute__ ((always_inline))
extern inline
void arch$cortexm$SetIPSR(uint32 r) {
	asm volatile ("msr ipsr, %0" :: "r" (r) : "ipsr");
}

__attribute__ ((always_inline))
extern inline
uint32 arch$cortexm$EPSR() {
	uint32 r;
	asm volatile ("mrs %0, epsr" : "=r" (r));
	return r;
}

__attribute__ ((always_inline))
extern inline
uint32 cortexm$IEPSR() {
	uint32 r;
	asm volatile ("mrs %0, iepsr" : "=r" (r));
	return r;
}

__attribute__ ((always_inline))
extern inline
void arch$cortexm$SetIEPSR(uint32 r) {
	asm volatile ("msr iepsr, %0" :: "r" (r) : "iepsr");
}

__attribute__ ((always_inline))
extern inline
uint32 arch$cortexm$IAPSR() {
	uint32 r;
	asm volatile ("mrs %0, iapsr" : "=r" (r));
	return r;
}

__attribute__ ((always_inline))
extern inline
void arch$cortexm$SetIAPSR(uint32 r) {
	asm volatile ("msr iapsr, %0" :: "r" (r) : "apsr");
}

__attribute__ ((always_inline))
extern inline
uint32 arch$cortexm$EAPSR() {
	uint32 r;
	asm volatile ("mrs %0, eapsr" : "=r" (r));
	return r;
}

__attribute__ ((always_inline))
extern inline
void arch$cortexm$SetEAPSR(uint32 r) {
	asm volatile ("msr eapsr, %0" :: "r" (r) : "eapsr");
}

__attribute__ ((always_inline))
extern inline
uint32 arch$cortexm$PSR() {
	uint32 r;
	asm volatile ("mrs %0, psr" : "=r" (r));
	return r;
}

__attribute__ ((always_inline))
extern inline
void arch$cortexm$SetPSR(uint32 r) {
	asm volatile ("msr psr, %0" :: "r" (r) : "psr");
}

__attribute__ ((always_inline))
extern inline
uintptr arch$cortexm$MSP() {
	uintptr p;
	asm volatile ("mrs %0, msp" : "=r" (p));
	return p;
}

__attribute__ ((always_inline))
extern inline
void arch$cortexm$SetMSP(unsafe$Pointer p) {
	asm volatile ("msr msp, %0" :: "r" (p) : "sp");
}

__attribute__ ((always_inline))
extern inline
uintptr arch$cortexm$PSP() {
	uintptr p;
	asm volatile ("mrs %0, psp" : "=r" (p));
	return p;
}

__attribute__ ((always_inline))
extern inline
void arch$cortexm$SetPSP(unsafe$Pointer p) {
	asm volatile ("msr psp, %0" :: "r" (p) : "sp");
}

__attribute__ ((always_inline))
extern inline
uint32 arch$cortexm$LR() {
	uint32 r;
	asm volatile ("mov %0, lr" : "=r" (r));
	return r;
}

__attribute__ ((always_inline))
extern inline
void arch$ortexm$SetLR(uint32 r) {
	asm volatile ("mov lr, %0" :: "r" (r) : "lr");
}

__attribute__ ((always_inline))
extern inline
arch$cortexm$Control arch$cortexm$Ctrl() {
	arch$cortexm$Control c;
	asm volatile ("mrs %0, control" : "=r" (c));
	return c;
}

__attribute__ ((always_inline))
extern inline
void arch$cortexm$SetCtrl(arch$cortexm$Control c) {
	asm volatile ("msr control, %0" :: "r" (c));
}

__attribute__ ((always_inline))
extern inline
void arch$cortexm$SEV() {
	asm volatile ("sev");
}

__attribute__ ((always_inline))
extern inline
void arch$cortexm$ISB() {
	asm volatile ("isb");
}

#define arch$cortexm$SVC(imm) asm volatile ("svc %0" :: "i" (imm))

#define arch$cortexm$BKPT(imm) asm volatile ("bkpt %0" :: "i" (imm))