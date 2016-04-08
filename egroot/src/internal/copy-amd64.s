// +build amd64

// func Memmove(dst, src unsafe.Pointer, n uintptr)
.global internal$Memmove

// unsafe$Pointer memmove(unsafe$Pointer dst, unsafe$Pointer, src, uint n)
.global memmove

// unsafe$Pointer memcpy(unsafe$Pointer dst, unsafe$Pointer src, uint n)
.global memcpy

internal$Memmove:
memmove:
memcpy:
	mov  %rdi, %rax
	cmp  $0, %rdx
	je   1f

	mov  %rdx, %rcx
	cmp  %rsi, %rdi
	jg   0f

	// Backward copy:
	std  
	add  %rdx, %rdi
	add  %rdx, %rsi
0:
	rep    
	movsb  
	cld    
1:
	ret  







