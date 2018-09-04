package sdmc

import (
	"rtos"

	"sdcard"
)

const busyTimeout = 500e6 // Timeout for SDXC: 500 ms (SDSC, SDHC: 250 ms).

func waitDataReady(h sdcard.Host) bool {
	return h.Wait(rtos.Nanosec() + busyTimeout)
}

// ReadBlocks reads buf.NumBlocks() 512-byte blocks from the card to buf
// starting from block number addr.
func (c *Card) ReadBlocks(addr int64, buf sdcard.Data) error {
	if buf.NumBlocks() == 0 {
		return nil
	}
	if uint64(addr) >= uint64(c.cap) {
		return ErrBadAddr
	}
	if c.oca&sdcard.HCXC == 0 {
		addr *= 512
	}
	h := c.host
	if !waitDataReady(h) {
		return ErrBusyTimeout
	}
	h.SetupData(sdcard.Recv|sdcard.Block512, buf.Words(), len(buf.Words())*8)
	var err error
	if buf.NumBlocks() == 1 {
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
	if buf.NumBlocks() == 0 {
		return nil
	}
	if uint64(addr) >= uint64(c.cap) {
		return ErrBadAddr
	}
	if c.oca&sdcard.HCXC == 0 {
		addr *= 512
	}
	h := c.host
	if !waitDataReady(h) {
		return ErrBusyTimeout
	}
	h.SetupData(sdcard.Send|sdcard.Block512, buf.Words(), len(buf.Words())*8)
	var err error
	if buf.NumBlocks() == 1 {
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
