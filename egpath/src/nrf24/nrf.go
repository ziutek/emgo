package nrf24

// DCI represents the simplified nRF24L01(+) Data and Control Interface: only SPI
// part of full DCI is need.
type DCI interface {
	// WriteRead perform SPI conversation: sets CSN low, writes and reads oi
	// data and sets CSN high.
	WriteRead(oi ...[]byte) (n int, err error)
}

// Radio provides interface to nRF24L01(+) transceiver. Radio has many methods
// that are mainly used to send commands that read or write its internal
// registers. Every such command, as side effect, always reads the value of
// STATUS regster as it was just before the command was executed. This status
// value is always returned as the last return value of the command method.
type Radio struct {
	DCI DCI

	// Status is the value of the STATUS register received by the last executed
	// command.
	//Status Status

	err error
}

// NewRadio provides convenient way to create heap allocated Radio struct.
func NewRadio(dci DCI) *Radio {
	dev := new(Radio)
	dev.DCI = dci
	return dev
}

// Err returns the error value of the last executed command. You can freely
// invoke many commands before check an error. If one command have caused an
// error the subsequent commands will not be executed until ClearErr will be
// called.
func (r *Radio) Err() error {
	return r.err
}

// ClearErr clears internal error variable.
func (r *Radio) ClearErr() {
	r.err = nil
}
