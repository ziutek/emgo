package hdcfb

import "unsafe"

func (sl *Slice) Write(s []byte) (int, error) {
	return sl.WriteString(*(*string)(unsafe.Pointer(&s)))
}

func (sl *SyncSlice) Write(s []byte) (int, error) {
	return sl.WriteString(*(*string)(unsafe.Pointer(&s)))
}
