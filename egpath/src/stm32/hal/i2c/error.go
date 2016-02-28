package i2c

type Error int16

const (
	BusErr   Error = 1 << 0
	ArbLost  Error = 1 << 1
	AckFail  Error = 1 << 2
	Overrun  Error = 1 << 3
	PECErr   Error = 1 << 4
	Timeout  Error = 1 << 6
	SMBAlert Error = 1 << 7

	SoftTimeout Error = 1 << 8
	BelatedStop Error = 1 << 9
	ActiveRead  Error = 1 << 10 // Write when active read transaction.
	DMAErr      Error = 1 << 11
)

func (e Error) Error() string {
	return "I2C error"
}