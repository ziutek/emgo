// +build cortexm3 cortexm4 cortexm4f

.syntax unified

// func Memset(s unsafe.Pointer, b byte, n uintptr)
.global internal$Memset

// unsafe$Pointer memset(unsafe$Pointer s, byte b, uintptr n)
.global memset

.thumb_func
internal$Memset:
.thumb_func
memset:
	// Fast path for 0, 1 bytes.
	cmp     r2, 1
	blo     0f  // n == 0
	itt     eq
	strbeq  r1, [r0]  // n == 1
	bxeq    lr

	// Setup full word of bytes.
	and  r1, 0xff
	orr  r1, r1, r1, lsl 8
	orr  r1, r1, r1, lsl 16

	// Fast path for 2, 3, 4 bytes.
	cmp     r2, 4
	bhi     1f  // n > 4
	itt     eq
	streq   r1, [r0]  // n == 4
	bxeq    lr
	cmp     r2, 2
	itt     eq
	strheq  r1, [r0]  // n == 2
	bxeq    lr
	strh    r1, [r0]  // n == 3
	strb    r1, [r0, 2]
0:
	bx  lr
1:

	// Use ip as dst. r0 will be returned unmodified.
	mov  ip, r0

	// Calculate the number of bytes to copy to make dst (ip) word aligned.
	ands  r3, ip, 3
	beq   5f
	rsb   r3, 4

	// Head set (up to 3 bytes).
	subs  r2, r3
	tbb   [pc, r3]
0:
	.byte  (4f-0b)/2
	.byte  (3f-0b)/2
	.byte  (2f-0b)/2
	.byte  (1f-0b)/2
1:
	strb  r1, [ip], 1
2:
	strb  r1, [ip], 1
3:
	strb  r1, [ip], 1
4:

5:
	// Set words.
	subs   r2, 4
	itt    hs
	strhs  r1, [ip], 4
	bhs    5b

	// Restore the number of remaining bytes.
	adds  r2, 4

	// Tail set (up to 3 bytes).
	tbb  [pc, r2]
0:
	.byte  (4f-0b)/2
	.byte  (3f-0b)/2
	.byte  (2f-0b)/2
	.byte  (1f-0b)/2
1:
	strb  r1, [ip], 1
2:
	strb  r1, [ip], 1
3:
	strb  r1, [ip], 1
4:

	bx  lr




