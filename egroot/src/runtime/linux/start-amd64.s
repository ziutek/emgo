// +build amd64

.global _start

_start:
	// Move argc and argv to rsi, rdx.
	pop  %rdi
	mov  %rsp, %rsi
	jmp  start






