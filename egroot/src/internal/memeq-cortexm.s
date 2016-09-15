// +build cortexm3 cortexm4 cortexm4f

.syntax unified

// func Memeq(p1, p2 unsafe.Pointer, n uintptr) bool
.global internal$Memeq

// bool memeq(unsafe$Pointer p1, unsafe$Pointer p2, uintptr n)
.global memeq

.thumb_func
internal$Memeq:
.thumb_func
memeq:
	// Go to head/tail check if n < 4.
	cmp    r2, 4
	itt    lo
	movlo  r3, r2
	blo    5f

	// Calculate number of bytes to check
	// to make p1 (r0) word aligned.
	ands   r3, r0, 3
	itt    ne
	rsbne  r3, 4
	bne    5f

	// Check words.
6:
	subs  r2, 4
	blo   1f
0:
	ldr   r3, [r0], 4
	ldr   ip, [r1], 4
	cmp   r3, ip
	bne   9f
	subs  r2, 4
	bhs   0b
1:
	adds  r2, 4
	beq   8f
	mov   r3, r2

	// Head/tail check.
5:
	// Check up to 3 bytes.
	tbb  [pc, r3]
0:
	.byte  (4f-0b)/2
	.byte  (3f-0b)/2
	.byte  (2f-0b)/2
	.byte  (1f-0b)/2
1:
	ldrb  r3, [r0], 1
	ldrb  ip, [r1], 1
	cmp   r3, ip
	bne   9f
2:
	ldrb  r3, [r0], 1
	ldrb  ip, [r1], 1
	cmp   r3, ip
	bne   9f
3:
	ldrb  r3, [r0], 1
	ldrb  ip, [r1], 1
	cmp   r3, ip
	bne   9f
4:
	subs  r2, r3
	bne   6b
8:
	adds  r0, r2, 1  // Set r0 to 1.
	bx    lr
9:
	subs  r0, r0  // Set r0 to 0.
	bx    lr














