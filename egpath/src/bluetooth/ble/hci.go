package ble

// HCI defines Host Controller Interface for Bluetooth LE.
type HCI interface {
	Recv() (DataPDU, error)
	Send(DataPDU) error
}
