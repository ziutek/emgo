package startup

// Start is default reset handler.
// This function is called by Cortex-M CPU after reset. You usually don't want to
// call it from user code (calling it results somewhat incomplete "Soft Reset"):
// 1. Data in RAM are initialized from Flash..
// 2. main package is initialized and main.main() is called.
// 3. Peripherals and CPU (eg. stack pointer) aren't affected.
// 4. This function never returns.
func Start()


func DefaultHandler() {
	for {
	}
}