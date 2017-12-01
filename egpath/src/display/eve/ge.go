package eve

// GE provides convenient way to write Graphics Engine commands. Every command
// is a function call, so for better performance or lower RAM usage, use raw
// Writer with many Graphics Engine commands in array/slice.
type GE struct {
	DL
}
