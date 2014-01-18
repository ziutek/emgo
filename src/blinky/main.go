package main

import (
	"delay"
	"stm32/clock"
	"stm32/flash"
	"stm32/gpio"
	"stm32/periph"
)

func stm32f4init() {
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

func setupLEDpins() {
	periph.AHB1ClockEnable(periph.GPIOD)
	periph.AHB1Reset(periph.GPIOD)

	gpio.D.SetMode(Green, gpio.Out)
	gpio.D.SetMode(Orange, gpio.Out)
	gpio.D.SetMode(Red, gpio.Out)
	gpio.D.SetMode(Blue, gpio.Out)
}

const (
	w1 = 1e6
	w2 = 1e7
)

func loop() {
	gpio.D.ResetBit(Green)
	gpio.D.SetBit(Orange)
	gpio.D.SetBit(Red)
	delay.Loop(w1)
	gpio.D.ResetBit(Red)
	gpio.D.ResetBit(Orange)
	gpio.D.SetBit(Blue)
	delay.Loop(w1)
	gpio.D.ResetBit(Blue)
	gpio.D.SetBit(Orange)
	gpio.D.SetBit(Red)
	delay.Loop(w1)
	gpio.D.ResetBit(Red)
	gpio.D.ResetBit(Orange)
	gpio.D.SetBit(Green)
	delay.Loop(w2)
}

func main() {
	stm32f4init()

	setupLEDpins()

	for {
		loop()
	}
}
