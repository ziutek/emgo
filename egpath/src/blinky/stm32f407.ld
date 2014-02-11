

MEMORY {
	flash (rx)	: ORIGIN = 0x08000000, LENGTH = 1024K
	ccram (rwx)	: ORIGIN = 0x10000000, LENGTH =   64K
 	ram (rwx)	: ORIGIN = 0x20000000, LENGTH =  128K
}

PROVIDE(mainStack = ORIGIN(ram) + LENGTH(ram));

SECTIONS {
	.text : {
		KEEP(*(.vectors))    /* Interrupt vector table */
		*(.text .text.*)     /* Code */   
		*(.rodata .rodata.*) /* Read only data */
    	/* } > ram */
	} >flash
    
	.data : {
		. = ALIGN(4);
		dataStart = .;
		
		*(.data .data.*)
		
		. = ALIGN(4);
		dataEnd = .;
	/* } >ram */
	} >ram AT>flash
    
	dataStartFlash = LOADADDR(.data);

	.bss : {
		. = ALIGN(4);
		bssStart = .;
    	
		*(.bss .bss.*)
        
		. = ALIGN(4);
		bssEnd = .;
	} >ram
}  

ENTRY(cortexm_startup_Start)
