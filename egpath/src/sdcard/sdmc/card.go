// Package sdmc implements access to Secure Digital Memory Cards.
package sdmc

import (
	"errors"

	"sdcard"
)

var (
	ErrInitIC = errors.New("sdmc: init IC")
	ErrInitOC = errors.New("sdmc: init OC")
	ErrStatus = errors.New("sdmc: status")
)

type Card struct {
	host   sdcard.Host
	oca    sdcard.OCR
	status sdcard.CardStatus
	rca    uint16
}

func MakeCard(host sdcard.Host) Card {
	return Card{host: host}
}

func NewCard(host sdcard.Host) *Card {
	c := new(Card)
	*c = MakeCard(host)
	return c
}

// Status returns last received card status.
func (c *Card) Status() sdcard.CardStatus {
	return c.status
}

func (c *Card) statusCmd(cmd sdcard.Command, arg uint32) {
	c.status = c.host.SendCmd(cmd, arg).R1()
}

func (c *Card) checkErr() error {
	if err := c.host.Err(true); err != nil {
		c.host.SetClock(0)
		return err
	}
	errFlags := sdcard.OUT_OF_RANGE |
		sdcard.ADDRESS_ERROR |
		sdcard.BLOCK_LEN_ERROR |
		sdcard.ERASE_SEQ_ERROR |
		sdcard.ERASE_PARAM |
		sdcard.WP_VIOLATION |
		sdcard.LOCK_UNLOCK_FAILED |
		sdcard.COM_CRC_ERROR |
		sdcard.ILLEGAL_COMMAND |
		sdcard.CARD_ECC_FAILED |
		sdcard.CC_ERROR |
		sdcard.ERROR |
		sdcard.CSD_OVERWRITE |
		sdcard.WP_ERASE_SKIP |
		sdcard.CARD_ECC_DISABLED |
		sdcard.ERASE_RESET |
		sdcard.AKE_SEQ_ERROR
	if c.status&errFlags != 0 {
		return ErrStatus
	}
	return nil
}
