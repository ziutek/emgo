// +build cortexm0

.syntax unified

// unsafe$Pointer memmove(unsafe$Pointer dst, unsafe$Pointer, src, uint n)
.global memmove

// unsafe$Pointer memcpy(unsafe$Pointer dst, unsafe$Pointer src, uint n)
.global memcpy

.thumb_func
memmove:
.thumb_func
memcpy:
	cmp   r2, 0
	bne  0f	
	bx lr
0:	
	mov  ip, r0
	cmp  r1, r0
	blo  2f

// Forward copy
1:
	ldrb  r3, [r1]
	strb  r3, [r0]
	adds  r1, 1
	adds  r0, 1
	subs  r2, 1
	bne   1b
	
	mov  r0, ip
	bx   lr	

// Backward copy:
2:
	ldrb  r3, [r1]
	strb  r3, [r0]
	subs  r1, 1
	subs  r0, 1
	subs  r2, 1
	bne   2b
	
	mov  r0, ip
	bx   lr	
