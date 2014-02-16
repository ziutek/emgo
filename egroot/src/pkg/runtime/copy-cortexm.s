// +build cortexm3 cortexm4 cortexm4f

.syntax unified

.global runtime_Copy
.global memcpy
.global memmove


// TODO: Implement full aligned copy using LDM and STM
// and check difference on real Cortex-M application:
// microcontroler with small SRAM and instructions read
// from Flash (all Flash acceleration on).


// func Copy(dst, src unsafe.Pointer, n uint)

.thumb_func
runtime_Copy:
.thumb_func
memcpy:
.thumb_func
memmove:
	// TODO: Check is better to always use
	// forward copy on non-overlaping data.
	cmp r0, r1
	blo 10f

// Forward copy

	cmp		r2, #4
	itt		lo
	movlo	r3, r2
	blo		5f
	
	// calculate number of bytes to copy
	// to make dst (r0) word aligned
	ands	r3, r0, #3
	it		ne
	rsbne	r3, #4
5:
	// copy up to 3 bytes
	subs	r2, r3
	tbb		[pc, r3]
0:
	.byte	(4f-0b)/2
	.byte	(3f-0b)/2
	.byte	(2f-0b)/2
	.byte	(1f-0b)/2
1:
	ldrb	r3, [r1], #1
	strb	r3, [r0], #1
2:
	ldrb	r3, [r1], #1
	strb	r3, [r0], #1
3:
	ldrb	r3, [r1], #1
	strb	r3, [r0], #1
4:
	beq		20f
	// copy words
0:
	subs 	r2, #4
	ittt	hs
	ldrhs	r3, [r1], #4
	strhs	r3, [r0], #4
	bhs		0b
	
	add		r2,	#4
	beq		20f
	
	// tail copy
	
	mov		r3, r2
	b		5b

10:
// Backward copy:

	add		r1, r2
	add		r0, r2	

	cmp		r2, #4
	itt		lo
	movlo	r3, r2
	blo		5f
	
	// calculate number of bytes to copy
	// to make dst (r0) word aligned
	ands	r3, r0, #3
5:
	// copy up to 3 bytes
	subs	r2, r3
	tbb		[pc, r3]
0:
	.byte	(4f-0b)/2
	.byte	(3f-0b)/2
	.byte	(2f-0b)/2
	.byte	(1f-0b)/2
1:
	ldrb	r3, [r1], #-1
	strb	r3, [r0], #-1
2:
	ldrb	r3, [r1], #-1
	strb	r3, [r0], #-1
3:
	ldrb	r3, [r1], #-1
	strb	r3, [r0], #-1
4:
	beq		20f
	// copy words
0:
	subs 	r2, #4
	ittt	hs
	ldrhs	r3, [r1], #-4
	strhs	r3, [r0], #-4
	bhs		0b
	
	adds	r2,	#4
	beq		20f
	
	// tail copy
	mov		r3, r2
	b		5b	
	
20:
	bx lr
