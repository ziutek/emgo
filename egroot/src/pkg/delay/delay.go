package delay

// Loop can be used to perform short active delay implemented
// in C as follows:
//	while (n > 0) {
//		asm volatile ("nop");
//		--n;
//	}
func Loop(n int)
