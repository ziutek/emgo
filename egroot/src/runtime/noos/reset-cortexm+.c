// +build cortexm0 cortexm3 cortexm4 cortexm4f

int_ runtime$init();
int_ main$init();
int_ main$main();

// All external symbols as byte to prevent compiler to optimize any runtime
// align checks.
extern byte DataRAM, DataLoad, DataSize;
extern byte BSSRAM, BSSSize;

void
runtime$noos$reset() {
	memmove(&DataRAM, &DataLoad, (uint) (&DataSize));
	memset(&BSSRAM, 0, (uint) (&BSSSize));

	runtime$init();
	main$init();
	main$main();

	for (;;);
}
