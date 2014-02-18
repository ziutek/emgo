_MainStack = ORIGIN(RAM) + LENGTH(RAM);

SECTIONS {
	.text : ALIGN(4) {
		KEEP(*(.vectors))    /* Interrupt vector table */
		*(.text .text.*)     /* Code */   
	} >loadTo
	
	.rodata : ALIGN(4) {
		*(.rodata .rodata.*) /* Read only data */
	} >loadTo
    
	.data : ALIGN(4) {
		*(.data .data.*)
	} >RAM AT>loadTo
	
	_DataRAM = ADDR(.data);
	_DataLoad = LOADADDR(.data);
	_DataLen = SIZEOF(.data);
	
	.bss : ALIGN(4) {
		*(.bss .bss.*)
	} >RAM
	
	_BSSRAM = ADDR(.bss);
	_BSSLen = SIZEOF(.bss);
}  

ENTRY(cortexm_startup_Start)
