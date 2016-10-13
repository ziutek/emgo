package nrf24

// DCI represents the simplified nRF24L01(+) Data and Control Interface: only SPI
// part of full DCI is need.
type DCI interface {
	// WriteRead perform SPI conversation: sets CSN low, writes and reads oi
	// data and sets CSN high.
	WriteRead(oi ...[]byte) (n int, err error)
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
