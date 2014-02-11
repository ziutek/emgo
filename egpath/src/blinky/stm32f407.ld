

MEMORY {
	flash (rx)	: ORIGIN = 0x08000000, LENGTH = 1024K
	ccram (rwx)	: ORIGIN = 0x10000000, LENGTH =   64K
 	ram (rwx)	: ORIGIN = 0x20000000, LENGTH =  128K
}

PROVIDE(mainstack = ORIGIN(ram) + LENGTH(ram));

SECTIONS {
	.text : {
		KEEP(*(.vectors))    /* Interrupt vector table */
		*(.text .text.*)     /* Code */   
		*(.rodata .rodata.*) /* Read only data */
    	/* } > ram */
	} >flash
    
	.data : {
		. = ALIGN(4);
		__dataStart = .;
		
		*(.data .data.*)
		
		. = ALIGN(4);
		__dataEnd = .;
	/* } >ram */
	} >ram AT>flash
    
	__dataStartFlash = LOADADDR(.data);

	.bss : {
		. = ALIGN(4);
		__bssStart = .;
    	
		*(.bss .bss.*)
        
		. = ALIGN(4);
		__bssEnd = .;
	} >ram
}  

ENTRY(cortexm_startup_Start)
