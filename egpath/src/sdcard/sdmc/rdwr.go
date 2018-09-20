package sdmc

import (
	"sdcard"
)

// ReadBlocks reads buf.NumBlocks() 512-byte blocks from the card to buf
// starting from block number addr.
func (c *Card) ReadBlocks(addr int64, buf sdcard.Data) error {
	nblocks := buf.NumBlocks()
	if nblocks == 0 {
		return nil
	}
	if uint64(addr) >= uint64(c.cap) {
		return ErrBadAddr
	}
	if c.oca&sdcard.HCXC == 0 {
		addr *= 512
	}
	h := c.host
	h.SetupData(sdcard.Recv|sdcard.Block512, buf.Words(), nblocks*512)
	var err error
	if nblocks == 1 {
		err = c.statusCmd(sdcard.CMD17(uint(addr)))
	} else {
		err = c.statusCmd(sdcard.CMD18(uint(addr)))
		if err != nil {
			return err
		}
		err = c.statusCmd(sdcard.CMD12())
	}
	return err
}

// WriteBlocks writes buf.NumBlocks() 512-byte blocks from buf to the card
// starting at block number addr.
func (c *Card) WriteBlocks(addr int64, buf sdcard.Data) error {
	nblocks := buf.NumBlocks()
	if nblocks == 0 {
		return nil
	}
	if uint64(addr) >= uint64(c.cap) {
		return ErrBadAddr
	}
	if c.oca&sdcard.HCXC == 0 {
		addr *= 512
	}
	h := c.host
	h.SetupData(sdcard.Send|sdcard.Block512, buf.Words(), nblocks*512)
	var err error
	if nblocks == 1 {
		err = c.statusCmd(sdcard.CMD24(uint(addr)))
	} else {
		err = c.statusCmd(sdcard.CMD25(uint(addr)))
		if err != nil {
			return err
		}
		err = c.statusCmd(sdcard.CMD12())
	}
	return err
}
