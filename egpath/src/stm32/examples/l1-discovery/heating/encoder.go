package main

import (
	"sync/atomic"

	"arch/cortexm/bitband"

	"stm32/hal/exti"
	"stm32/hal/raw/tim"
)

type EncoderState int

func (es EncoderState) Cnt() int {
	return int(es >> 1)
}

func (es EncoderState) ModCnt(m int) int {
	return (m + es.Cnt()%m) % m
}

func (es EncoderState) Btn() bool {
	return es&1 != 0
}

type Encoder struct {
	t       *tim.TIM_Periph
	b       bitband.Bit
	base    uint32
	lastCnt uint32
	lastBtn int

	State chan EncoderState
}

func (e *Encoder) Init(t *tim.TIM_Periph, b bitband.Bit, l exti.Lines) {
	e.t = t
	e.b = b
	e.State = make(chan EncoderState, 3)

	t.CKD().Store(2 << tim.CKDn)
	t.SMCR.StoreBits(tim.SMS, 1<<tim.SMSn)
	t.CCMR1.StoreBits(
		tim.CC1S|tim.CC2S|tim.IC1F|tim.IC2F,
		1<<tim.CC1Sn|1<<tim.CC2Sn|0xf<<tim.IC1Fn|0xf<<tim.IC2Fn,
	)
	t.CCER.SetBits(tim.CC1P | tim.CC2P)
	t.CNT.Store(0)
	t.CCR3.Store(0xffffffff)
	t.CCR4.Store(1)
	t.DIER.Store(tim.CC3IE | tim.CC4IE)
	t.CEN().Set()

	l.EnableFallTrig()
	l.EnableRisiTrig()
	l.EnableIRQ()
}

func (e *Encoder) SetCnt(cnt int) {
	atomic.StoreUint32(&e.base, e.t.CNT.U32.Load()-uint32(cnt))
}

func (e *Encoder) ISR() {
	for {
		exti.L4.ClearPending()
		e.t.SR.Store(0)
		cnt := e.t.CNT.U32.Load()
		btn := e.b.Load()
		if cnt == e.lastCnt && btn == e.lastBtn {
			return
		}
		if cnt != e.lastCnt {
			e.t.CCR3.Store(tim.CCR3_Bits(cnt - 1))
			e.t.CCR4.Store(tim.CCR4_Bits(cnt + 1))
		}
		e.lastCnt = cnt
		e.lastBtn = btn
		cnt -= atomic.LoadUint32(&e.base)
		select {
		case e.State <- EncoderState(int(int16(cnt))<<1 + 1 - btn):
			//ledGreen.Set()
			//delay.Loop(1e4)
			//ledGreen.Clear()
		default:
		}
	}
}

var encoder Encoder

func encoderISR() {
	encoder.ISR()
}
