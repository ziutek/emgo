
void (*runtime$noos$utof64(uintptr u)) (int64) {
	return CAST(void (*)(int64), u);
}

int64(*runtime$noos$utofr64(uintptr u)) () {
	return CAST(int64(*)(), u);
}
