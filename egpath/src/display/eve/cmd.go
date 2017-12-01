package eve

// GE allows to write Graphics Engine commands.
type GE Writer

func (ge GE) DL() DL {
	return DL(cmd)
}
