package nrf24

// DCI represents the simplified nRF24L01(+) Data and Control Interface: only
// SPI part of the full DCI is need.
type DCI interface {
	// WriteRead performs full SPI transation: sets CSN low, writes and reads oi
	// data and sets CSN high.
	WriteRead(oi ...[]byte) int
	
	// Err returns and clears internal error status.
	Err(clear bool) error
}

// Radio provides interface to nRF24L01(+) transceiver. Radio has many methods
// that are mainly used to send commands that read or write its internal
// registers. Every such command, as side effect, always reads the value of
// STATUS regster as it was just before the command was executed. This status
// value is always returned as the last return value of the command method.
type Radio struct {
	DCI DCI
}

// NewRadio provides convenient way to create heap allocated Radio struct.
func NewRadio(dci DCI) *Radio {
	dev := new(Radio)
	dev.DCI = dci
	return dev
}

// Err returns the error value of the last executed command. You can freely
// invoke many commands before check an error. If one command have caused an
// error the subsequent commands will not be executed until Err(true) will be
// called.
func (r *Radio) Err(clear bool) error {
	return r.DCI.Err(clear)
}
