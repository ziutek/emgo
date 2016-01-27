// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl

package rtc

import (
	"mmio"

	"stm32/hal/raw/bkp"
	"stm32/hal/raw/rtc"
)

func waitForSync(RTC *rtc.RTC_Periph) {
	RTC.RSF().Clear()
	for RTC.RSF().Load() == 0 {
	}
}

func waitForWrite(RTC *rtc.RTC_Periph) {
	for RTC.RTOFF().Load() == 0 {
	}
}

type twoReg struct {
	high, low *mmio.U16
}

func (tr twoReg) Load() uint32 {
	return uint32(tr.high.Load())<<16 | uint32(tr.low.Load())
}

func (tr twoReg) Store(u uint32) {
	tr.high.Store(uint16(u >> 16))
	tr.low.Store(uint16(u))
}

type fourReg struct {
	hh, hm, ml, ll *mmio.U16
}

func (tr fourReg) Load() uint64 {
	h := uint32(tr.hh.Load())<<16 | uint32(tr.hm.Load())
	l := uint32(tr.ml.Load())<<16 | uint32(tr.ll.Load())
	return uint64(h)<<32 | uint64(l)
}

func (tr fourReg) Store(u uint64) {
	v := uint32(u >> 32)
	tr.hh.Store(uint16(v >> 16))
	tr.hm.Store(uint16(v))
	v = uint32(u)
	tr.ml.Store(uint16(v >> 16))
	tr.ll.Store(uint16(v))
}

type rtcBackup struct {
	p *bkp.BKP_Periph
}

func (b rtcBackup) Status() *mmio.U16 {
	return &b.p.DR1.U16
}

func (b rtcBackup) CntExt() *mmio.U16 {
	return &b.p.DR2.U16
}

func (b rtcBackup) LastISR() twoReg {
	return twoReg{&b.p.DR3.U16, &b.p.DR4.U16}
}

func (b rtcBackup) StartSec() fourReg {
	return fourReg{&b.p.DR5.U16, &b.p.DR6.U16, &b.p.DR7.U16, &b.p.DR8.U16}
}

func (b rtcBackup) StartNanosec() twoReg {
	return twoReg{&b.p.DR9.U16, &b.p.DR10.U16}
}

/*
const dbg = itm.Port(17)

func print64(s string, i int64) {
	dbg.WriteString(s)
	strconv.WriteInt64(dbg, i, 16, 0)
}

func println64(s string, i int64) {
	dbg.WriteString(s)
	strconv.WriteInt64(dbg, i, 16, 0)
	dbg.WriteString("\r\n")
}

func print32(s string, u uint32) {
	dbg.WriteString(s)
	strconv.WriteUint32(dbg, u, 16, 0)
}

func println32(s string, u uint32) {
	dbg.WriteString(s)
	strconv.WriteUint32(dbg, u, 16, 0)
	dbg.WriteString("\r\n")
}
*/
