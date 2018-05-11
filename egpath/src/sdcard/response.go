package sdcard

type Response [4]uint32

func (r Response) R1() CardStatus {
	return CardStatus(r[0])
}

func (r Response) R2CID() CID {
	return CID(r)
}

func (r Response) R2CSD() CSD {
	return CSD(r)
}

func (r Response) R3() OCR {
	return OCR(r[0])
}

func (r Response) R7() (vhs VHS, pattern byte) {
	return VHS(r[0] >> 8 & 0xF), byte(r[0])
}

type CardStatus uint32

const (
	AKE_SEQ_ERROR      CardStatus = 1 << 3
	APP_CMD            CardStatus = 1 << 5
	FX_EVENT           CardStatus = 1 << 6
	READY_FOR_DATA     CardStatus = 1 << 8
	CURRENT_STATE      CardStatus = 15 << 9
	StateIdle          CardStatus = 0 << 9
	StateReady         CardStatus = 1 << 9
	StateIdent         CardStatus = 2 << 9
	StateStby          CardStatus = 3 << 9
	StateTran          CardStatus = 4 << 9
	StateData          CardStatus = 5 << 9
	StateRcv           CardStatus = 6 << 9
	StatePrg           CardStatus = 7 << 9
	StateDis           CardStatus = 8 << 9
	StateIOOnly        CardStatus = 15 << 9
	ERASE_RESET        CardStatus = 1 << 13
	ARD_ECC_DISABLED   CardStatus = 1 << 14
	WP_ERASE_SKIP      CardStatus = 1 << 15
	CSD_OVERWRITE      CardStatus = 1 << 16
	ERROR              CardStatus = 1 << 19
	CC_ERROR           CardStatus = 1 << 20
	CARD_ECC_FAILED    CardStatus = 1 << 21
	ILLEGAL_COMMAND    CardStatus = 1 << 22
	COM_CRC_ERROR      CardStatus = 1 << 23
	LOCK_UNLOCK_FAILED CardStatus = 1 << 24
	CARD_IS_LOCKED     CardStatus = 1 << 25
	WP_VIOLATION       CardStatus = 1 << 26
	ERASE_PARAM        CardStatus = 1 << 27
	ERASE_SEQ_ERROR    CardStatus = 1 << 28
	BLOCK_LEN_ERROR    CardStatus = 1 << 29
	ADDRESS_ERROR      CardStatus = 1 << 30
	OUT_OF_RANGE       CardStatus = 1 << 31
)

type OCR uint32

const (
	DVC   OCR = 1 << 7  // Dual Voltage Card
	V28   OCR = 1 << 15 // 2.7-2.8 V
	V29   OCR = 1 << 16 // 2.8-2.9 V
	V30   OCR = 1 << 17 // 2.9-3.0 V
	V31   OCR = 1 << 18 // 3.0-3.1 V
	V32   OCR = 1 << 19 // 3.1-3.2 V
	V33   OCR = 1 << 20 // 3.2-3.3 V
	V34   OCR = 1 << 21 // 3.3-3.4 V
	V35   OCR = 1 << 22 // 3.4-3.5 V
	V36   OCR = 1 << 23 // 3.5-3.6 V
	S18   OCR = 1 << 24 // Switching to 1.8V
	XPC   OCR = 1 << 18 // SDXC maximum performance
	UHSII OCR = 1 << 29 // UHS-II Card Status
	HCXC  OCR = 1 << 30 // Card Capacity Status (set fot SDHC, SDXC).
	PWUP  OCR = 1 << 31 // Card in power up state (^Busy).
)

type CID [4]uint32

// MID returns the manufacturer ID.
func (cid CID) MID() byte {
	return byte(cid[0] >> 24)
}

// OID returns the OEM/Application ID.
func (cid CID) OID() [2]byte {
	return [2]byte{byte(cid[0] >> 16), byte(cid[0] >> 8)}
}

// PNM returns the product name.
func (cid CID) PNM() [5]byte {
	return [5]byte{
		byte(cid[0]), byte(cid[1] >> 24), byte(cid[1] >> 16), byte(cid[1] >> 8),
		byte(cid[1]),
	}
}

// PRV returns the product revision.
func (cid CID) PRV() byte {
	return byte(cid[2] >> 24)
}

// PSN returns the product serial number.
func (cid CID) PSN() uint32 {
	return cid[2]<<8 | cid[3]>>24
}

// MDT returns the manufacturing date.
func (cid CID) MDT() (year, month int) {
	mdt := int(cid[3] >> 8 & 0xFFF)
	return mdt>>4 + 2000, mdt & 15
}

// CRC returns the 7-bit CRC field.
func (cid CID) CRC() byte {
	return byte(cid[3] >> 1)
}

type CSD [4]uint32
