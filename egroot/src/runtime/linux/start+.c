int runtime$init();
int main$init();
int main$main();

void
start(uintptr argc, void **argv) {
	internal$Argv = (uintptr) argv;
	argv += argc + 1;
	internal$Env = (uintptr) argv;
	while (*argv != nil) {
		argv++;
	}
	argv++;
	internal$Auxv = (uintptr) argv;

	runtime$init();
	main$init();
	main$main();
	runtime$linux$exit(0);
}
