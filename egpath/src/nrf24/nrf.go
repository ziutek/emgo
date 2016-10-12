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

// Radio provides interface to nRF24L01(+) transceiver.
type Radio struct {
	DCI DCI

	// Status is the value of the STATUS register received by the last executed
	// command.
	Status Status

	// Err is the error value of the last executed command. You can freely
	// invoke many commands before check an error. If one command have returned
	// an error the subsequent commands will not be executed as long as
	// Err != nil.
	Err error
}

// NewRadio provides convenient way to create heap allocated Radio struct.
func NewRadio(dci DCI) *Radio {
	dev := new(Radio)
	dev.DCI = dci
	return dev
}

func (d *Radio) SetCE(v int) error {
	d.DCI.SetCE(v)
	return nil
}
