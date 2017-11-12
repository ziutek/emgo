package ili9341

// DCI stands for Display Controller Interface / Data and Control Interface.
type DCI interface {
	Cmd(b byte)       // Cmd invokes a command (8-bit word size).
	WriteByte(b byte) // WriteByte pass one byte of data (8-bit word size).

	SetWordSize(size int) // SetWordSize changes the data word to size bits.

	Cmd2(w uint16)        // Cmd2 invokes two commands (16-bit word size).
	WriteWord(w uint16)   // Word passes one word of data (16-bit word size).
	Write(data []uint16)  // Write passes many words of data (16-bit word size).
	Fill(w uint16, n int) // Fill passes a word n times (16-bit word size).

	Err() error // Err returns and clears internal error variable.
}
