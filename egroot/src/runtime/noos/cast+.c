
void (*runtime$noos$utof64(uintptr u)) (uint64) {
	return CAST(void (*)(uint64), u);
}

uint64(*runtime$noos$utofr64(uintptr u)) () {
	return CAST(uint64(*)(), u);
}
