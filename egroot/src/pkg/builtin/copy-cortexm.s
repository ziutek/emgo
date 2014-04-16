// +build cortexm0 cortexm3 cortexm4 cortexm4f

.syntax unified

.global memcpy
.global memmove


// TODO: Implement full aligned copy using LDM and STM
// and check difference on real Cortex-M application:
// a microcontroler with small SRAM and instructions read
// from Flash (all Flash acceleration on).

// void memmove(unsafe_Pointer dst, unsafe.Pointer, src, uint n)
// void memcpy(unsafe.Pointer dst, unsafe.Pointer src, uint n)

.thumb_func
memmove:
.thumb_func
memcpy:
	// TODO: Check will be better to always
	// use forward copy on non-overlaping data.
	cmp  r0, r1
	blo  10f

// Forward copy
	cmp    r2, 4
	itt    lo
	movlo  r3, r2
	blo    5f

	// calculate number of bytes to copy
	// to make dst (r0) word aligned
	ands   r3, r0, 3
	itt    ne
	rsbne  r3, 4
	bne    5f

	// copy words
6:
	subs   r2, 4
	ittt   hs
	ldrhs  r3, [r1], 4
	strhs  r3, [r0], 4
	bhs    6b

	adds  r2, 4
	beq   9f
	mov   r3, r2

	// head/tail copy
5:
	// copy up to 3 bytes
	subs  r2, r3
	tbb   [pc, r3]
0:
	.byte  (4f-0b)/2
	.byte  (3f-0b)/2
	.byte  (2f-0b)/2
	.byte  (1f-0b)/2
1:
	ldrb  r3, [r1], 1
	strb  r3, [r0], 1
2:
	ldrb  r3, [r1], 1
	strb  r3, [r0], 1
3:
	ldrb  r3, [r1], 1
	strb  r3, [r0], 1
4:
	bne  6b
9:
	bx  lr

// Backward copy:
10:
	add  r1, r2
	add  r0, r2

	cmp    r2, 4
	itt    lo
	movlo  r3, r2
	blo    5f

	// calculate number of bytes to copy
	// to make dst (r0) word aligned
	ands  r3, r0, 3
	bne   5f

	// copy words
6:
	subs   r2, 4
	ittt   hs
	ldrhs  r3, [r1, -4]!
	strhs  r3, [r0, -4]!
	bhs    6b

	adds  r2, 4
	beq   9f
	mov   r3, r2

	// head/tail copy
5:
	// copy up to 3 bytes
	subs  r2, r3
	tbb   [pc, r3]
0:
	.byte  (4f-0b)/2
	.byte  (3f-0b)/2
	.byte  (2f-0b)/2
	.byte  (1f-0b)/2
1:
	ldrb  r3, [r1, -1]!
	strb  r3, [r0, -1]!
2:
	ldrb  r3, [r1, -1]!
	strb  r3, [r0, -1]!
3:
	ldrb  r3, [r1, -1]!
	strb  r3, [r0, -1]!
4:
	bne  6b
9:
	bx  lr
