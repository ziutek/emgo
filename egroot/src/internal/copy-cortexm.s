// +build cortexm3 cortexm4 cortexm4f cortexm7f cortexm7d

.syntax unified

// TODO: Implement full aligned copy using LDM and STM and check does this
// improves performance.

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
	// Fast path for 0, 1, 2, 3, 4 bytes (uses only 16-bit instructions: 44 B).
	cmp     r2, 1
	blo     0f  // n == 0
	ittt    eq
	ldrbeq  r3, [r1]  // n == 1
	strbeq  r3, [r0]
	bxeq    lr
	cmp     r2, 4
	bhi     1f  // n > 4
	ittt    eq
	ldreq   r3, [r1]  // n == 4
	streq   r3, [r0]
	bxeq    lr
	cmp     r2, 2
	ittt    eq
	ldrheq  r3, [r1]  // n == 2
	strheq  r3, [r0]
	bxeq    lr
	ldrh    r3, [r1]  // n == 3
	ldrb    r2, [r1, 2]
	strh    r3, [r0]
	strb    r2, [r0, 2]
0:
	bx  lr
1:

	// Use ip as dst. r0 will be returned unmodified.
	mov  ip, r0

	cmp  r1, r0
	blo  10f

	// Forward copy

	// Calculate the number of bytes to copy to make dst (ip) word aligned.
	ands  r3, ip, 3
	beq   5f
	rsb   r3, 4

	// Head copy (up to 3 bytes).
	subs  r2, r3
	tbb   [pc, r3]
0:
	.byte  (4f-0b)/2
	.byte  (3f-0b)/2
	.byte  (2f-0b)/2
	.byte  (1f-0b)/2
1:
	ldrb  r3, [r1], 1
	strb  r3, [ip], 1
2:
	ldrb  r3, [r1], 1
	strb  r3, [ip], 1
3:
	ldrb  r3, [r1], 1
	strb  r3, [ip], 1
4:

5:
	// Copy words.
	subs   r2, 4
	ittt   hs
	ldrhs  r3, [r1], 4
	strhs  r3, [ip], 4
	bhs    5b

	// Restore the number of remaining bytes.
	adds  r2, 4

	// Tail copy (up to 3 bytes).
	tbb  [pc, r2]
0:
	.byte  (4f-0b)/2
	.byte  (3f-0b)/2
	.byte  (2f-0b)/2
	.byte  (1f-0b)/2
1:
	ldrb  r3, [r1], 1
	strb  r3, [ip], 1
2:
	ldrb  r3, [r1], 1
	strb  r3, [ip], 1
3:
	ldrb  r3, [r1], 1
	strb  r3, [ip], 1
4:

	bx  lr

// Backward copy:
10:
	add  r1, r2
	add  ip, r2

	// Calculate the number of bytes to copy to make dst (ip) word aligned.
	ands  r3, ip, 3
	beq   5f

	// Tail copy (up to 3 bytes).
	subs  r2, r3
	tbb   [pc, r3]
0:
	.byte  (4f-0b)/2
	.byte  (3f-0b)/2
	.byte  (2f-0b)/2
	.byte  (1f-0b)/2
1:
	ldrb  r3, [r1, -1]!
	strb  r3, [ip, -1]!
2:
	ldrb  r3, [r1, -1]!
	strb  r3, [ip, -1]!
3:
	ldrb  r3, [r1, -1]!
	strb  r3, [ip, -1]!
4:

5:
	// Copy words.
	subs   r2, 4
	ittt   hs
	ldrhs  r3, [r1, -4]!
	strhs  r3, [ip, -4]!
	bhs    5b

	// Restore the number of remaining bytes.
	adds  r2, 4

	// Head copy (up to 3 bytes).
	tbb  [pc, r2]
0:
	.byte  (4f-0b)/2
	.byte  (3f-0b)/2
	.byte  (2f-0b)/2
	.byte  (1f-0b)/2
1:
	ldrb  r3, [r1, -1]!
	strb  r3, [ip, -1]!
2:
	ldrb  r3, [r1, -1]!
	strb  r3, [ip, -1]!
3:
	ldrb  r3, [r1, -1]!
	strb  r3, [ip, -1]!
4:

	bx  lr























