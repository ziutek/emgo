// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl f10x_cl

.syntax unified

// This code allows to boot program from SRAM on STM32F1xx.

.global bootSRAM

.section .text.bootcode

.thumb_func
bootSRAM:
	ldr  r0, =VectorsStart  // Load address of vector table.
	ldr  r1, =0xE000ED08    // Load address of NVIC VTOR register.
	str  r0, [r1]           // Set VTOR to VectorsStart.
	ldr  sp, [r0]           // Load SP from exception vector 0 (initial SP)
	ldr  pc, [r0, 4]        // Jump to exception vector 1 (reset code).

