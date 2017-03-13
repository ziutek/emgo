package gpio

type Mode byte

const (
	In     Mode = 0
	InOut  Mode = 1
	Discon Mode = 2
	Out    Mode = 3
)

type Pull byte

const (
	NoPull   Pull = 0
	PullDown Pull = 1
	PullUp   Pull = 3
)

type Config struct {
	Mode Mode
	Pull Pull
}

func (p *Port) SetupPin(n int, cfg *Config) {
	p.pincnf[n].Store(uint32(cfg.Pull)<<2 | uint32(cfg.Mode))
}

// Setup configures pins.
func (p *Port) Setup(pins Pins, cfg *Config) {
	for n := 0; n < 32; n++ {
		if pins&(1<<uint(n)) != 0 {
			p.SetupPin(n, cfg)
		}
	}
}
