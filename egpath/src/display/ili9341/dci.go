package ili9341

type DCI interface {
	Cmd(b byte)
	Byte(b byte)

	SetWordSize(size int)

	Cmd16(w uint16)
	Word(w uint16)
	Data(data []uint16)
	Fill(w uint16, n int)

	Err() error
}
