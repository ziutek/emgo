package sdcard

import (
	"reflect"
	"unsafe"
)

// Data can be used to access bytes sent/received using SD card data transfers.
// It ensures 8-byte alignment required by Host.SetupData method. The bit order
// of its 64-bit elements is hardware depend. Use Bytes method to return Data as
// correctly ordered string of bytes. Additionally it helps to operate on
// typical 512 byte blocks.
type Data struct {
	words []uint64
}

// MakeDataBytes allocates and initializes Data object that can store n bytes.
func MakeDataBytes(n int) Data {
	return Data{make([]uint64, (n+7)/8)}
}

// MakeDataBlocks allocates and initializes Data object that can store n
// 512-byte blocks.
func MakeDataBlocks(n int) Data {
	return Data{make([]uint64, n*64)}
}

// AsData returns Data object that refers to d.
func AsData(d []uint64) Data {
	return Data{d}
}

// Words returns d as []uint64.
func (d Data) Words() []uint64 {
	return d.words
}

// Bytes returns d as []byte.
func (d Data) Bytes() []byte {
	h := (*reflect.SliceHeader)(unsafe.Pointer(&d.words))
	h.Len *= 8
	h.Cap *= 8
	return *(*[]uint8)(unsafe.Pointer(h))
}

// ByteSlice aligns m, n down to 8-byte boundary and the returns slice of d that
// contains (n>>3 - m>>3) * 8 bytes from aligned m to aligned n.
/*func (d Data) ByteSlice(m, n int) []byte {
	return Data{d.words[m>>3 : n>>3]}.Bytes()
}*/

// Block returns slice of d that contains n-th 512-byte block.
func (d Data) Block(n int) Data {
	n *= 64
	return Data{d.words[n : n+64]}
}

// Block returns the slice of d that contains n-m 512-byte blocks from m to n.
func (d Data) BlockSlice(m, n int) Data {
	return Data{d.words[m*64 : n*64]}
}

// NumBlocks returns the number of full 512-byte blocks that can fit into d.
func (d Data) NumBlocks() int {
	return len(d.words) >> 6
}
