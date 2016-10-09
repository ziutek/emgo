package nrf24

// DCI represents simplified nRF24L01(+) Data and Control Interface.
// It allows to perform two basic operations: communicate with nRF24L01(+) over
// SPI and enable/disable its RF part.
type DCI interface {
	// WriteRead perform SPI conversation: sets CSN low, writes and reads oi
	// data, sets CSN high.
	WriteRead(oi ...[]byte) (n int, err error)

	// Set CE line. v==0 sets CE low, v==1 sets CE high, v==2 pulses CE high for
	// 10 µs and leaves it low.
	SetCE(v int) error
}

// Device provides interface to nRF24L01(+) transceiver.
type Device struct {
	DCI DCI

	// Status is value of status register just before the last executed command.
	Status Status

	// Err is error value of last executed command. You can freely invoke many
	// commands before check an error. If one command have returned an error
	// the subsequent commands will not be executed.
	Err error
}

func (d *Device) SetCE(v int) error {
	d.DCI.SetCE(v)
	return nil
}
