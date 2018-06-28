package sdmc

import (
	"sdcard"
)

// ReadBlocks reads 512-byte blocks from card to buf. It reads no more than
// buf.NumBlock() blocks. It returns number of blocks read.
func (c *Card) ReadBlocks(buf sdcard.Data) (n int, err error) {
	return 0, nil
}
