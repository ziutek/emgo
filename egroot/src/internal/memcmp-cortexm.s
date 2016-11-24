// +build cortexm3 cortexm4 cortexm4f cortexm7f cortexm7d

.syntax unified

// func Memcmp(p1, p2 unsafe.Pointer, n uintptr) int
.global internal$Memcmp

.thumb_func
internal$Memcmp:
	mov  ip, r0  // To improve 16-bit/32-bit instruction ratio.

	// Go to tail check if n < 4.
	cmp  r2, 4
	blo  6f

	// Calculate the number of bytes to check to make p1 (ip) word aligned.
	ands  r0, 3
	beq   5f
	rsb   r0, 4

	// Perform head check (up to 3 bytes).
	subs  r2, r0
	tbb   [pc, r0]
0:
	.byte  (4f-0b)/2
	.byte  (3f-0b)/2
	.byte  (2f-0b)/2
	.byte  (1f-0b)/2
1:
	ldrb  r0, [ip], 1
	ldrb  r3, [r1], 1
	subs  r0, r3
	bne   7f
2:
	ldrb  r0, [ip], 1
	ldrb  r3, [r1], 1
	subs  r0, r3
	bne   7f
3:
	ldrb  r0, [ip], 1
	ldrb  r3, [r1], 1
	subs  r0, r3
	bne   7f
4:

5:
	// Check words.
	subs  r2, 4
	blo   1f
0:
	ldr   r0, [ip], 4
	ldr   r3, [r1], 4
	cmp   r0, r3
	bne   8f
	subs  r2, 4
	bhs   0b
1:
	adds  r2, 4
	beq   9f

6:
	// Perform tail check (up to 3 bytes).
	tbb  [pc, r2]
0:
	.byte  (4f-0b)/2
	.byte  (3f-0b)/2
	.byte  (2f-0b)/2
	.byte  (1f-0b)/2
1:
	ldrb  r0, [ip], 1
	ldrb  r3, [r1], 1
	subs  r0, r3
	bne   7f
2:
	ldrb  r0, [ip], 1
	ldrb  r3, [r1], 1
	subs  r0, r3
	bne   7f
3:
	ldrb  r0, [ip], 1
	ldrb  r3, [r1], 1
	subs  r0, r3
4:

7:
	bx  lr

8:
	// Non-equal words. Convert litle-endian words to big-endian.
	rev   r0, r0
	rev   r3, r3
9:
	subs  r0, r3
	bx    lr

































