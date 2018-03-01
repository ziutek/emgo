package main

import (
	"stm32/hal/exti"
	"stm32/hal/gpio"
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
	t     *tim.TIM_Periph
	btnIn gpio.Pin
	btnEL exti.Lines
	lastCnt uint32
	lastBtn int

	State chan EncoderState
}

func (e *Encoder) Init(t *tim.TIM_Periph, btn gpio.Pin) {
	e.t = t
	e.btnIn = btn
	e.btnEL = exti.LineIndex(btn.Index())
	e.State = make(chan EncoderState, 1)

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
	t.CR1.Store(2<<tim.CKDn | tim.CEN)

	e.btnEL.Connect(btn.Port())
	e.btnEL.EnableFallTrig()
	e.btnEL.EnableRisiTrig()
	e.btnEL.EnableIRQ()
}

func (e *Encoder) SetCnt(cnt int) {
	e.btnEL.DisableIRQ()
	e.t.CEN().Clear()
	if e.t.CNT.U32.Load() != uint32(cnt) {
		e.t.CNT.U32.Store(uint32(cnt))
		// Remove possible invalid state from channel.
		select {
		case <-e.State:
		default:
		}
		e.ISR()
	}
	e.t.CEN().Set()
	e.btnEL.EnableIRQ()
}

func (e *Encoder) ISR() {
	for {
		e.btnEL.ClearPending()
		e.t.SR.Store(0)
		cnt := e.t.CNT.U32.Load()
		btn := e.btnIn.Load()
		if cnt == e.lastCnt && btn == e.lastBtn {
			return
		}
		if cnt != e.lastCnt {
			e.t.CCR3.Store(tim.CCR3(cnt - 1))
			e.t.CCR4.Store(tim.CCR4(cnt + 1))
		}
		e.lastCnt = cnt
		e.lastBtn = btn
		select {
		case e.State <- EncoderState(int(int16(cnt))<<1 + 1 - btn):
		default:
		}
	}
}

var encoder Encoder

func encoderISR() {
	encoder.ISR()
}
