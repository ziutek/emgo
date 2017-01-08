package hdcfb

import "unsafe"

// Write: see WriteString.
func (sl *Slice) Write(s []byte) (int, error) {
	return sl.WriteString(*(*string)(unsafe.Pointer(&s)))
}

// Write: see Slice.WriteString.
func (sl *SyncSlice) Write(s []byte) (int, error) {
	return sl.WriteString(*(*string)(unsafe.Pointer(&s)))
}
