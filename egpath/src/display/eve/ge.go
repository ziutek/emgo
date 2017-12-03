package eve

// GE provides a convenient way to write Graphics Engine commands. Every command
// is a function call, so for better performance or lower RAM usage, use raw
// Writer with many Graphics Engine commands in array / slice.
type GE struct {
	DL
}

// Start starts a new display list.
func (ge GE) DLStart() {
	ge.w.wr32(CMD_DLSTART)
}

// Swap swapx the current display list.
func (ge GE) Swap() {
	ge.w.wr32(CMD_SWAP)
}

// ColdStart sets co-processor engine state to default values.
func (ge GE) ColdStart() {
	ge.w.wr32(CMD_COLDSTART)
}

// Interrupt triggers interrupt INT_CMDFLAG.
func (ge GE) Interrupt() {
	ge.w.wr32(CMD_INTERRUPT)
}

// Append appends more commands to current display list.
func (ge GE) Append(addr, num int) {
	ge.w.wr32(CMD_APPEND)
	ge.w.wr32(uint32(addr))
	ge.w.wr32(uint32(num))
}

// RegRead reads a register value.
func (ge GE) RegRead(addr int) {
	ge.w.wr32(CMD_REGREAD)
	ge.w.wr32(uint32(addr))
}

// MemWrite writes bytes into memory.
func (ge GE) MemWrite(addr, num int) {
	ge.w.wr32(CMD_MEMWRITE)
	ge.w.wr32(uint32(addr))
	ge.w.wr32(uint32(num))
}
