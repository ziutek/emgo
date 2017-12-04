package eve

// GE provides a convenient way to write Graphics Engine commands. Every command
// is a function call, so for better performance or lower RAM usage, use raw
// Writer with many Graphics Engine commands in array / slice.
type GE struct {
	DL
}

// DLStart starts a new display list.
func (ge GE) DLStart() {
	ge.aw32(CMD_DLSTART)
}

// Swap swaps the current display list.
func (ge GE) Swap() {
	ge.aw32(CMD_SWAP)
}

// ColdStart sets co-processor engine state to default values.
func (ge GE) ColdStart() {
	ge.aw32(CMD_COLDSTART)
}

// Interrupt triggers interrupt INT_CMDFLAG.
func (ge GE) Interrupt() {
	ge.aw32(CMD_INTERRUPT)
}

// Append appends more commands resident in RAM_G to the current display list
// memory address where the offset is specified in REG_CMD_DL.
func (ge GE) Append(addr, num int) {
	ge.aw32(CMD_APPEND)
	ge.wr32(uint32(addr))
	ge.wr32(uint32(num))
}

// RegRead reads a register value.
func (ge GE) RegRead(addr int) {
	ge.aw32(CMD_REGREAD)
	ge.wr32(uint32(addr))
}

// MemWrite writes the following bytes into memory.
func (ge GE) MemWrite(addr, num int) {
	ge.aw32(CMD_MEMWRITE)
	ge.wr32(uint32(addr))
	ge.wr32(uint32(num))
}

// Inflate decompresses the following compressed data into RAM_G.
func (ge GE) Inflate(addr int) {
	ge.aw32(CMD_INFLATE)
	ge.wr32(uint32(addr))
}

// LoadImage decompresses the following JPEG image data into a bitmap, in RAM_G
// (EVE2 supports also PNG).
func (ge GE) LoadImage(addr int, options uint32) {
	ge.aw32(CMD_LOADIMAGE)
	ge.wr32(uint32(addr))
	ge.wr32(options)
}

// MediaFIFO sets up a streaming media FIFO in RAM_G.
func (ge GE) MediaFIFO(addr, size int) {
	ge.aw32(CMD_MEDIAFIFO)
	ge.wr32(uint32(addr))
	ge.wr32(uint32(size))
}

// PlayVideo plays back MJPEG-encoded AVI video.
func (ge GE) PlayVideo(options uint32) {
	ge.aw32(CMD_PLAYVIDEO)
	ge.wr32(options)
}

// VideoStart initializes the AVI video decoder.
func (ge GE) VideoStart() {
	ge.aw32(CMD_VIDEOSTART)
}

// VideoFrame loads the next frame of video.
func (ge GE) VideoFrame(dst, ptr int) {
	ge.aw32(CMD_VIDEOFRAME)
	ge.wr32(uint32(dst))
	ge.wr32(uint32(ptr))
}

// MemCRC computes a CRC-32 for a block of EVE memory.
func (ge GE) MemCRC(addr, num int) {
	ge.aw32(CMD_MEMCRC)
	ge.wr32(uint32(addr))
	ge.wr32(uint32(num))
}

// MemZero writes zero to a block of memory.
func (ge GE) MemZero(addr, num int) {
	ge.aw32(CMD_MEMZERO)
	ge.wr32(uint32(addr))
	ge.wr32(uint32(num))
}

// MemSet fills memory with a byte value.
func (ge GE) MemSet(addr int, val byte, num int) {
	ge.aw32(CMD_MEMSET)
	ge.wr32(uint32(addr))
	ge.wr32(uint32(val))
	ge.wr32(uint32(num))
}

// MemCpy copies a block of memory.
func (ge GE) MemCpy(dst, src, num int) {
	ge.aw32(CMD_MEMCPY)
	ge.wr32(uint32(dst))
	ge.wr32(uint32(src))
	ge.wr32(uint32(num))
}

func (ge GE) ButtonRaw(x, y, w, h int, font, options uint16) {
	ge.aw32(CMD_BUTTON)
	ge.wr32(uint32(x)&0xFFFF | uint32(y)&0xFFFF<<16)
	ge.wr32(uint32(w)&0xFFFF | uint32(h)&0xFFFF<<16)
	ge.wr32(uint32(font) | uint32(options)<<16)
}

// Button draws a button.
func (ge GE) Button(x, y, w, h int, font, options uint16, s string) {
	ge.ButtonRaw(x, y, w, h, font, options)
	ge.wrs(s)
	ge.wr8(0)
}

// Close closes the write transaction and returns number of bytes written,
// rounded up to multiple of 4 (to avoid rounding use ge.Writer.Close).
func (ge GE) Close() int {
	return (ge.Writer.Close() + 3) &^ 3
}
