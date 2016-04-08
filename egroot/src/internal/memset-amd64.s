// +build amd64

// func Memset(s unsafe.Pointer, b byte, n uintptr)
.global internal$Memset

// unsafe$Pointer memset(unsafe$Pointer s, byte b, uintptr n)
.global memset

internal$Memset:
memset:
	mov    %rsi, %rax
	mov    %rdx, %rcx
	mov    %rdi, %rsi
	rep    
	stosb  
	mov    %rsi, %rax
	ret    
