package nrf24

func (r *Radio) exec(args ...[]byte) Status {
	if r.err != nil {
		return 0
	}
	_, r.err = r.DCI.WriteRead(args...)
	return Status(args[1][0])
}

// Reg invokes R_REGISTER command.
func (r *Radio) Reg(addr byte, val []byte) Status {
	cmd := []byte{addr}
	return r.exec(cmd, cmd, nil, val)
}

// SetReg invokes W_REGISTER command.
func (r *Radio) SetReg(addr byte, val ...byte) Status {
	cmd := []byte{addr | 0x20}
	return r.exec(cmd, cmd, val)
}

func checkPayLen(pay []byte) {
	if len(pay) > 32 {
		panic("plen>32")
	}
}

// ReadRxP invokes R_RX_PAYLOAD command.
func (r *Radio) ReadRx(pay []byte) Status {
	checkPayLen(pay)
	cmd := []byte{0x61}
	return r.exec(cmd, cmd, nil, pay)
}

// WriteTxP invokes W_TX_PAYLOAD command.
func (r *Radio) WriteTx(pay []byte) Status {
	checkPayLen(pay)
	cmd := []byte{0xa0}
	return r.exec(cmd, cmd, pay)
}

// FlushTx invokes FLUSH_TX command.
func (r *Radio) FlushTx() Status {
	cmd := []byte{0xe1}
	return r.exec(cmd, cmd)
}

// FlushRx invokes FLUSH_RX command.
func (r *Radio) FlushRx() Status {
	cmd := []byte{0xe2}
	return r.exec(cmd, cmd)
}

// ReuseTx invokes REUSE_TX_PL command.
func (r *Radio) ReuseTx() Status {
	cmd := []byte{0xe3}
	return r.exec(cmd, cmd)
}

// Activate invokes nRF24L01 ACTIVATE command.
func (r *Radio) Activate(b byte) Status {
	cmd := []byte{0x50}
	return r.exec(cmd, cmd)
}

// RxLen invokes R_RX_PL_WID command.
func (d *Radio) RxLen() (int, Status) {
	cmd := []byte{0x60, 0xff}
	s := d.exec(cmd, cmd)
	return int(cmd[1]), s
}

// WriteAck invokes W_ACK_PAYLOAD command.
func (r *Radio) WriteAck(pn int, pay []byte) Status {
	checkPayNum(pn)
	checkPayLen(pay)
	cmd := []byte{byte(0xa8 | pn)}
	return r.exec(cmd, cmd, pay)
}

// WriteTxNoAck invokes W_TX_PAYLOAD_NOACK command.
func (r *Radio) WriteTxNoAck(pay []byte) Status {
	checkPayLen(pay)
	cmd := []byte{0xa0}
	return r.exec(cmd, cmd, pay)
}

// NOP invokes NOP command.
func (r *Radio) NOP() Status {
	cmd := []byte{0xff}
	return r.exec(cmd, cmd)
}
