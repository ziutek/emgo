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
