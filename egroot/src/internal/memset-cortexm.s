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
	// Use ip as dst. r0 will be returned unmodified.
	mov  ip, r0
	and  r1, 0xff
	orr  r1, r1, r1, lsl 8
	orr  r1, r1, r1, lsl 16

	cmp    r2, 4
	itt    lo
	movlo  r3, r2
	blo    5f

	// calculate number of bytes to set
	// to make s (ip) word aligned
	ands   r3, ip, 3
	itt    ne
	rsbne  r3, 4
	bne    5f

	// set words
6:
	subs   r2, 4
	itt    hs
	strhs  r1, [ip], 4
	bhs    6b

	adds  r2, 4
	beq   9f
	mov   r3, r2

	// head/tail set
5:
	// set up to 3 bytes
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
	bne  6b
9:
	bx  lr
