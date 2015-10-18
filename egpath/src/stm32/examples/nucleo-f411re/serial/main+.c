uintptr$$uintptr main$slisiz() {
	return (uintptr$$uintptr){sizeof(slice), __alignof__(slice)};
}

uintptr$$uintptr main$int64siz() {
	return (uintptr$$uintptr){sizeof(int64), __alignof__(int64)};
}

uintptr$$uintptr main$cplx128siz() {
	return (uintptr$$uintptr){sizeof(complex128), __alignof__(complex128)};
}

uintptr$$uintptr main$interfacesiz() {
	return (uintptr$$uintptr){sizeof(interface), __alignof__(interface)};
}

uintptr main$sssize() {
	return sizeof(main$SS);
}