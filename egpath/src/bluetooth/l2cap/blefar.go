package l2cap

import (
	"io"

	"bluetooth/ble"
)

// BLEFAR performs Fragmentation And Recombination of L2CAP PDUs using BLE HCI.
type BLEFAR struct {
	hci ble.HCI
}

func (far *BLEFAR) SetHCI(hci ble.HCI) {
	far.hci = hci
}

func (far *BLEFAR) ReadPDU() (r io.Reader, cid, length int) {
	return nil, 0, 0
}

func (far *BLEFAR) WritePDU(cid, length int) io.Writer {
	return nil
}
