arm-none-eabi-gdb --tui \
	-ex "target extended-remote localhost:4242" \
	-ex "set remote hardware-breakpoint-limit 6" \
	-ex "set remote hardware-watchpoint-limit 4" \
	main.elf
