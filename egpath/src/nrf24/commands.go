package nrf24

func (r *Radio) exec(args ...[]byte) STATUS {
	if r.err != nil {
		return 0
	}
	_, r.err = r.DCI.WriteRead(args...)
	return STATUS(args[1][0])
}

func (r *Radio) R_REGISTER(addr byte, val []byte) STATUS {
	cmd := []byte{addr}
	return r.exec(cmd, cmd, nil, val)
}

func (r *Radio) W_REGISTER(addr byte, val ...byte) STATUS {
	cmd := []byte{addr | 0x20}
	return r.exec(cmd, cmd, val)
}

func checkPayLen(pay []byte) {
	if len(pay) > 32 {
		panic("plen>32")
	}
}

func (r *Radio) R_RX_PAYLOAD(pay []byte) STATUS {
	checkPayLen(pay)
	cmd := []byte{0x61}
	return r.exec(cmd, cmd, nil, pay)
}

func (r *Radio) W_TX_PAYLOAD(pay []byte) STATUS {
	checkPayLen(pay)
	cmd := []byte{0xa0}
	return r.exec(cmd, cmd, pay)
}

func (r *Radio) FLUSH_TX() STATUS {
	cmd := []byte{0xe1}
	return r.exec(cmd, cmd)
}

func (r *Radio) FLUSH_RX() STATUS {
	cmd := []byte{0xe2}
	return r.exec(cmd, cmd)
}

func (r *Radio) REUSE_TX_PL() STATUS {
	cmd := []byte{0xe3}
	return r.exec(cmd, cmd)
}

// ACTIVATE is nRF24L01 specific command.
func (r *Radio) ACTIVATE(b byte) STATUS {
	cmd := []byte{0x50}
	return r.exec(cmd, cmd)
}

func (d *Radio) R_RX_PL_WID() (int, STATUS) {
	cmd := []byte{0x60, 0xff}
	s := d.exec(cmd, cmd)
	return int(cmd[1]), s
}

func (r *Radio) W_ACK_PAYLOAD(pn int, pay []byte) STATUS {
	checkPayNum(pn)
	checkPayLen(pay)
	cmd := []byte{byte(0xa8 | pn)}
	return r.exec(cmd, cmd, pay)
}

func (r *Radio) W_TX_PAYLOAD_NOACK(pay []byte) STATUS {
	checkPayLen(pay)
	cmd := []byte{0xa0}
	return r.exec(cmd, cmd, pay)
}

func (r *Radio) NOP() STATUS {
	cmd := []byte{0xff}
	return r.exec(cmd, cmd)
}
