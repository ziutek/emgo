int runtime$init();
int main$init();
int main$main();

void
_start(int argc, byte ** argv, byte ** env) {
	runtime$init();
	main$init();
	main$main();
	runtime$linux$exit(0);
}
