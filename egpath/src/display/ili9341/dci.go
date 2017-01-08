package ili9341

// DCI stands for Display Controller Interface / Data and Control Interface.
type DCI interface {
	Cmd(b byte)  // Cmd invokes a command (8-bit word size).
	Byte(b byte) // Byte pass one byte of data (8-bit word size).

	SetWordSize(size int) // SetWordSize changes the data word size.

	Cmd16(w uint16)       // Cmd invokes a two commands (16-bit word size).
	Word(w uint16)        // Word passes one word of data (16-bit word size).
	Data(data []uint16)   // Data passes many words of data (16-bit word size).
	Fill(w uint16, n int) // Data passes a word n times (16-bit word size).

	Err() error // Err returns and clears internal error variable.
}
