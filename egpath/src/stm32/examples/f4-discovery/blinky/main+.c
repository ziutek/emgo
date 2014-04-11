__attribute__ ((noinline))
int main_F(main_Args a, int b) {
	void *cur = __builtin_frame_address(0);
	void *prev = __builtin_frame_address(1);
	return (int)prev - (int)cur;
}
