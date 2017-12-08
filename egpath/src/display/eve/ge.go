package eve

import (
	"bits"
)

// GE provides a convenient way to write Graphics Engine commands. Every command
// is a function call, so for better performance or lower RAM usage, use raw
// Writer with many Graphics Engine commands in array / slice.
type GE struct {
	DL
}

// Close closes the write transaction and returns number of bytes written,
// rounded up to multiple of 4 (to avoid rounding use ge.Writer.Close).
func (ge GE) Close() int {
	return (ge.Writer.Close() + 3) &^ 3
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

// ButtonHeader writes only header of CMD_BUTTON command (without label string).
// Use Write* methods to write button label. Label string must be terminated
// with zero byte.
func (ge GE) ButtonHeader(x, y, w, h int, font, options uint16) {
	ge.aw32(CMD_BUTTON)
	ge.wr32(uint32(x)&0xFFFF | uint32(y)&0xFFFF<<16)
	ge.wr32(uint32(w)&0xFFFF | uint32(h)&0xFFFF<<16)
	ge.wr32(uint32(font) | uint32(options)<<16)
}

// Button draws a button.
func (ge GE) Button(x, y, w, h int, font, options uint16, s string) {
	ge.ButtonHeader(x, y, w, h, font, options)
	ge.ws(s)
	ge.wr8(0)
}

// Clock draws an analog clock.
func (ge GE) Clock(x, y, r int, options uint16, h, m, s, ms int) {
	ge.aw32(CMD_CLOCK)
	ge.wr32(uint32(x)&0xFFFF | uint32(y)&0xFFFF<<16)
	ge.wr32(uint32(r)&0xFFFF | uint32(options)<<16)
	ge.wr32(uint32(h)&0xFFFF | uint32(m)&0xFFFF<<16)
	ge.wr32(uint32(s)&0xFFFF | uint32(ms)&0xFFFF<<16)
}

// FgColor sets the foreground color.
func (ge GE) FgColor(rgb uint32) {
	ge.aw32(CMD_FGCOLOR)
	ge.wr32(rgb)
}

// BgColor sets the background color.
func (ge GE) BgColor(rgb uint32) {
	ge.aw32(CMD_BGCOLOR)
	ge.wr32(rgb)
}

// GradColor sets the 3D button highlight color.
func (ge GE) GradColor(rgb uint32) {
	ge.aw32(CMD_GRADCOLOR)
	ge.wr32(rgb)
}

// Gauge draws a gauge.
func (ge GE) Gauge(x, y, r int, options uint16, major, minor, val, max int) {
	ge.aw32(CMD_GAUGE)
	ge.wr32(uint32(x)&0xFFFF | uint32(y)&0xFFFF<<16)
	ge.wr32(uint32(r)&0xFFFF | uint32(options)<<16)
	ge.wr32(uint32(major)&0xFFFF | uint32(minor)&0xFFFF<<16)
	ge.wr32(uint32(val)&0xFFFF | uint32(max)&0xFFFF<<16)
}

//Gradienta draws a smooth color gradient.
func (ge GE) Gradient(x0, y0 int, rgb0 uint32, x1, y1 int, rgb1 uint32) {
	ge.aw32(CMD_GRADIENT)
	ge.wr32(uint32(x0)&0xFFFF | uint32(y0)&0xFFFF<<16)
	ge.wr32(rgb0)
	ge.wr32(uint32(x1)&0xFFFF | uint32(y1)&0xFFFF<<16)
	ge.wr32(rgb1)
}

// KeysHeader writes only header of CMD_KEYS command (without key labels). Use
// Write* methods to write key labels. Labels string must be terminated with
// zero byte.
func (ge GE) KeysHeader(x, y, w, h int, font, options uint16) {
	ge.aw32(CMD_KEYS)
	ge.wr32(uint32(x)&0xFFFF | uint32(y)&0xFFFF<<16)
	ge.wr32(uint32(w)&0xFFFF | uint32(h)&0xFFFF<<16)
	ge.wr32(uint32(font) | uint32(options)<<16)
}

// Keys draws a row of keys.
func (ge GE) Keys(x, y, w, h int, font, options uint16, s string) {
	ge.KeysHeader(x, y, w, h, font, options)
	ge.ws(s)
	ge.wr8(0)
}

// Progress draws a progress bar.
func (ge GE) Progress(x, y, w, h int, options uint16, val, max int) {
	ge.aw32(CMD_PROGRESS)
	ge.wr32(uint32(x)&0xFFFF | uint32(y)&0xFFFF<<16)
	ge.wr32(uint32(w)&0xFFFF | uint32(h)&0xFFFF<<16)
	ge.wr32(uint32(options) | uint32(val)&0xFFFF<<16)
	ge.wr32(uint32(max) & 0xFFFF)
}

// Progress draws a scroll bar.
func (ge GE) Scrollbar(x, y, w, h int, options uint16, val, size, max int) {
	ge.aw32(CMD_SCROLLBAR)
	ge.wr32(uint32(x)&0xFFFF | uint32(y)&0xFFFF<<16)
	ge.wr32(uint32(w)&0xFFFF | uint32(h)&0xFFFF<<16)
	ge.wr32(uint32(options) | uint32(val)&0xFFFF<<16)
	ge.wr32(uint32(size) | uint32(max)&0xFFFF<<16)
}

// Slider draws a slider.
func (ge GE) Slider(x, y, w, h int, options uint16, val, max int) {
	ge.aw32(CMD_SLIDER)
	ge.wr32(uint32(x)&0xFFFF | uint32(y)&0xFFFF<<16)
	ge.wr32(uint32(w)&0xFFFF | uint32(h)&0xFFFF<<16)
	ge.wr32(uint32(options) | uint32(val)&0xFFFF<<16)
	ge.wr32(uint32(max) & 0xFFFF)
}

// Dial draws a rotary dial control.
func (ge GE) Dial(x, y, r int, options uint16, val int) {
	ge.aw32(CMD_DIAL)
	ge.wr32(uint32(x)&0xFFFF | uint32(y)&0xFFFF<<16)
	ge.wr32(uint32(r)&0xFFFF | uint32(options)<<16)
	ge.wr32(uint32(val))
}

// ToggleHeader writes only header of CMD_TOGGLE command (without label string).
// Use Write* methods to write toggle label. Label string must be terminated
// with zero byte.
func (ge GE) ToggleHeader(x, y, w int, font, options uint16, state bool) {
	ge.aw32(CMD_TOGGLE)
	ge.wr32(uint32(x)&0xFFFF | uint32(y)&0xFFFF<<16)
	ge.wr32(uint32(w)&0xFFFF | uint32(font)<<16)
	ge.wr32(uint32(options) | uint32(bits.One(!state)-1)<<16)
}

// Toggle draws a toggle switch.
func (ge GE) Toggle(x, y, w int, font, options uint16, state bool, s string) {
	ge.ToggleHeader(x, y, w, font, options, state)
	ge.ws(s)
	ge.wr8(0)
}

// TextHeader writes only header of CMD_TEXT command (without text string). Use
// Write* methods to write text. Text string must be terminated with zero byte.
//  ge.TextHeader(40, 40, 18, 0)
//  fmt.Fprintf(ge, "Weight: %.1f kg\000", weight)
func (ge GE) TextHeader(x, y int, font, options uint16) {
	ge.aw32(CMD_TEXT)
	ge.wr32(uint32(x)&0xFFFF | uint32(y)&0xFFFF<<16)
	ge.wr32(uint32(font) | uint32(options)<<16)
}

// Text draws text.
func (ge GE) Text(x, y int, font, options uint16, s string) {
	ge.TextHeader(x, y, font, options)
	ge.ws(s)
	ge.wr8(0)
}

// SetBase sets the base for number output.
func (ge GE) SetBase(base int) {
	ge.aw32(CMD_TEXT)
	ge.wr32(uint32(base))
}

// Number draws number.
func (ge GE) Number(x, y int, font, options uint16, n int) {
	ge.aw32(CMD_NUMBER)
	ge.wr32(uint32(x)&0xFFFF | uint32(y)&0xFFFF<<16)
	ge.wr32(uint32(font) | uint32(options)<<16)
	ge.wr32(uint32(n))
}

// LoadIdentity instructs the graphics engine to set the current matrix to the
// identity matrix, so it is able to form the new matrix as requested by Scale,
// Rotate, Translate command.
func (ge GE) LoadIdentity() {
	ge.aw32(CMD_LOADIDENTITY)
}

// SetMatrix assigns the value of the current matrix to the bitmap transform
// matrix of the graphics engine by generating display list commands.
func (ge GE) SetMatrix(a, b, c, d, e, f int) {
	ge.aw32(CMD_SETMATRIX)
	ge.wr32(uint32(a))
	ge.wr32(uint32(b))
	ge.wr32(uint32(c))
	ge.wr32(uint32(d))
	ge.wr32(uint32(e))
	ge.wr32(uint32(f))
}

// GetMatrix retrieves the current matrix within the context of the graphics
// engine.
func (ge GE) GetMatrix() {
	ge.aw32(CMD_GETMATRIX)
}

// GetPtr gets the end memory address of data inflated by Inflate command.
func (ge GE) GetPtr() {
	ge.aw32(CMD_GETPTR)
}

// GetProps gets the image properties decompressed by LoadImage.
func (ge GE) GetProps() {
	ge.aw32(CMD_GETPROPS)
}

// Scale applies a scale to the current matrix.
func (ge GE) Scale(sx, sy int) {
	ge.aw32(CMD_SCALE)
	ge.wr32(uint32(sx))
	ge.wr32(uint32(sy))
}

// Rotate applies a rotation to the current matrix.
func (ge GE) Rotate(a int) {
	ge.aw32(CMD_ROTATE)
	ge.wr32(uint32(a))
}

// Translate applies a translation to the current matrix.
func (ge GE) Translate(tx, ty int) {
	ge.aw32(CMD_TRANSLATE)
	ge.wr32(uint32(tx))
	ge.wr32(uint32(ty))
}

// Calibrate execute the touch screen calibration routine.
func (ge GE) Calibrate() {
	ge.aw32(CMD_CALIBRATE)
}

// Sketch starts a continuous sketch update. It does not display anything, only
// draws to the bitmap located in RAM_G, at address addr.
func (ge GE) Sketch(x, y, w, h, addr int, format uint16) {
	ge.aw32(CMD_SKETCH)
	ge.wr32(uint32(x)&0xFFFF | uint32(y)&0xFFFF<<16)
	ge.wr32(uint32(w)&0xFFFF | uint32(h)&0xFFFF<<16)
	ge.wr32(uint32(addr))
	ge.wr32(uint32(format))
}
