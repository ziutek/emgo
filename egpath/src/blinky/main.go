package main

import (
	"delay"
	"stm32/clock"
	"stm32/flash"
	"stm32/gpio"
	"stm32/periph"
)

func init() {
	const flashLatency = 5 // Need for 2.7-3.6V and 150-168MHz

	flash.SetLatency(flashLatency)
	flash.SetPrefetch(true)
	flash.SetICache(true)
	flash.SetDCache(true)

	// Be sure that flash latency is set before incrase frequency.
	for flash.Latency() != flashLatency {
	}

	// Reset clock subsystem
	clock.ResetCR()
	clock.ResetPLLCFGR()
	clock.ResetCFGR()
	clock.ResetCIR()

	// Enable HSE clock
	clock.EnableHSE()
	for !clock.HSEReady() {
	}

	// Configure clocks for AHB, APB1, APB2 bus.
	clock.SetPrescalerAHB(clock.AHBDiv1)
	clock.SetPrescalerAPB1(clock.APBDiv4) // SysFreq / div <= 42 MHz
	clock.SetPrescalerAPB2(clock.APBDiv2) // SysFreq / div <= 84 MHz

	// Enable main PLL
	clock.SetPLLSrc(clock.PLLSrcHSE) // 8 MHz external oscilator
	clock.SetPLLInputDiv(4)          // 2 MHz
	clock.SetMainPLLMul(168)         // 336 MHz
	clock.SetMainPLLSysDiv(2)        // 168 MHz
	clock.SetMainPLLPeriphDiv(7)     // 48 MHz
	clock.EnableMainPLL()
	for !clock.MainPLLReady() {
	}

	// Set PLL as system clock source
	clock.SetSysClock(clock.PLL)
	for clock.SysClock() != clock.PLL {
	}
}

const (
	Green = 12 + iota
	Orange
	Red
	Blue
)

var LED = gpio.D

func init() {
	periph.AHB1ClockEnable(periph.GPIOD)
	periph.AHB1Reset(periph.GPIOD)

	LED.SetMode(Green, gpio.Out)
	LED.SetMode(Orange, gpio.Out)
	LED.SetMode(Red, gpio.Out)
	LED.SetMode(Blue, gpio.Out)
}

func Exported(p gpio.Port) {
	notExported(p)
}

func notExported(p gpio.Port) {
	p.SetBit(12)
}

func loop() {
	const (
		W1 = 1e6
		W2 = 1e7
	)
	var LED = LED

	LED.ResetBit(Green)
	LED.SetBit(Orange)
	delay.Loop(W1)
	LED.SetBit(Red)
	delay.Loop(W1)
	LED.ResetBit(Red)
	delay.Loop(W1)
	LED.ResetBit(Orange)
	LED.SetBit(Blue)
	delay.Loop(W1)
	LED.ResetBit(Blue)
	LED.SetBit(Orange)
	delay.Loop(W1)
	LED.SetBit(Red)
	delay.Loop(W1)
	LED.ResetBit(Red)
	delay.Loop(W1)
	LED.ResetBit(Orange)
	LED.SetBit(Green)
	delay.Loop(W2)
}

const (
	A = -1/2.1
	B = float32(1./-2)
	C = float64(1/2)
	D = 2
	E = uint64(2)
	F = uint32(3)
	G = int32(5)
)

func main() {
	a := float32(34)
	b := a + 2.0/3
	c := complex(1, a)
	d := complex(float32(1), -2)
	var e complex64 = 2i
	for {
		a = b + imag(c) + real(d) + imag(e)
		loop()
	}
}
