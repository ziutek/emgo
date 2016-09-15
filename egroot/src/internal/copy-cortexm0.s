// +build cortexm0

.syntax unified

// func Memmove(dst, src unsafe.Pointer, n uintptr)
.global internal$Memmove

// unsafe$Pointer memmove(unsafe$Pointer dst, unsafe$Pointer, src, uintptr n)
.global memmove

// unsafe$Pointer memcpy(unsafe$Pointer dst, unsafe$Pointer src, uintptr n)
.global memcpy

.thumb_func
internal$Memmove:
.thumb_func
memmove:
.thumb_func
memcpy:
	cmp  r2, 0
	bne  0f
	bx   lr
0:
	cmp  r1, r0
	blo  2f

// Forward copy
1:
	mov  ip, r0
0:
	ldrb  r3, [r1]
	strb  r3, [r0]
	adds  r1, 1
	adds  r0, 1
	subs  r2, 1
	bne   0b

	mov  r0, ip
	bx   lr

// Backward copy:
2:
	add  r0, r2
	add  r1, r2
0:
	subs  r1, 1
	subs  r0, 1
	ldrb  r3, [r1]
	strb  r3, [r0]
	subs  r2, 1
	bne   0b

	bx  lr
