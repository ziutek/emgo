package ble

// HCI defines Host Controller Interface for Bluetooth LE.
type HCI interface {
	// Recv returns next received PDU.
	Recv() (DataPDU, error)

	// GetSend returns current transmit PDU. This PDU is not zeroed (can
	// contain random data in payload and header).
	GetSend() DataPDU

	// Send sends current transmit PDU. 
	Send() error
}
