// +build cortexm0

.syntax unified

// func Memcmp(p1, p2 unsafe.Pointer, n uintptr) int
.global internal$Memcmp

.thumb_func
internal$Memcmp:
	cmp   r2, 0
	bne   0f
	eors  r0, r0  // Reurn 0 (equal) for empty strings.
	bx    lr
0:
	mov  ip, r4
	mov  r4, r0
1:
	ldrb  r0, [r4]
	ldrb  r3, [r1]
	adds  r4, 1
	adds  r1, 1
	subs  r0, r3
	bne   2f
	subs  r2, 1
	bne   1b
2:
	mov  r4, ip
	bx   lr





















































