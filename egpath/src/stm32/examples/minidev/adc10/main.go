package main

import (
	"fmt"
	"io"
	"rtos"

	"stm32/hal/adc"
	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
	"stm32/hal/usart"

	"stm32/hal/raw/rcc"
	"stm32/hal/raw/tim"

	"github.com/ziutek/emgo/egroot/src/delay"
)

var (
	adcd *adc.Driver
	adct *tim.TIM_Periph
	tts  *usart.Driver
)

const (
	IN0   = 0
	IN1   = 1
	IN3   = 2
	IN4   = 3
	INVCC = 4
	IN6   = 5
	IN7   = 6
	IN8   = 7
	IN9   = 8
	NIN   = 9
)

func init() {
	system.SetupPLL(8, 1, 36/8)
	systick.Setup(2e6)

	// GPIO

	var apins [NIN]gpio.Pin

	gpio.A.EnableClock(true)
	apins[IN0] = gpio.A.Pin(0)
	apins[IN1] = gpio.A.Pin(1)
	opin := gpio.A.Pin(2) // USART2_TX for one-wire bus
	apins[IN3] = gpio.A.Pin(3)
	apins[IN4] = gpio.A.Pin(4)
	apins[INVCC] = gpio.A.Pin(5)
	apins[IN6] = gpio.A.Pin(6)
	apins[IN7] = gpio.A.Pin(7)
	tx := gpio.A.Pin(9)
	rx := gpio.A.Pin(10)

	gpio.B.EnableClock(true)
	apins[IN8] = gpio.B.Pin(0)
	apins[IN9] = gpio.B.Pin(1)

	// DMA
	dma1 := dma.DMA1
	dma1.EnableClock(true)

	// USART

	tx.Setup(&gpio.Config{Mode: gpio.Alt})
	rx.Setup(&gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	tts = usart.NewDriver(
		usart.USART1, dma1.Channel(4, 0), dma1.Channel(5, 0), make([]byte, 40),
	)
	tts.Periph().EnableClock(true)
	tts.Periph().SetBaudRate(115200)
	tts.Periph().Enable()
	tts.EnableRx()
	tts.EnableTx()
	fmt.DefaultWriter = tts

	rtos.IRQ(irq.USART1).Enable()
	rtos.IRQ(irq.DMA1_Channel4).Enable()
	rtos.IRQ(irq.DMA1_Channel5).Enable()

	// ADC

	for _, pin := range apins {
		pin.Setup(&gpio.Config{Mode: gpio.Ana})
	}
	adcd = adc.NewDriver(adc.ADC1, dma1.Channel(1, 0))
	adcd.P.EnableClock(true)
	rcc.RCC.ADCPRE().Store(2 << rcc.ADCPREn) // ADCclk = APB2clk / 6 = 12 MHz

	rtos.IRQ(irq.ADC1_2).Enable()
	rtos.IRQ(irq.DMA1_Channel1).Enable()

	// ADC timer.

	rcc.RCC.TIM3EN().Set()
	adct = tim.TIM3
	adct.CR2.Store(2 << tim.MMSn) // Update event as TRGO.
	adct.CR1.Store(tim.CEN)

	// One-wire
	_ = opin
}

func main() {
	adcd.P.SetSamplTime(1, adc.MaxSamplTime(55.5*2)) // 55.5 + 12.5 = 68
	adcd.P.SetSequence(0, 1, 3, 4, 5, 6, 7, 8, 9)
	adcd.P.SetTrigSrc(adc.ADC12_TIM3_TRGO)
	adcd.P.SetTrigEdge(adc.EdgeRising)
	//adcd.P.SetAlignLeft(true)
	//adcd.SetReadMSB(true)

	adcd.Enable(true)

	// Max. SR = 36 MHz / 6 / 68 ≈ 88235 Hz

	div1, div2 := 72, 100 // ADC SR = 36 MHz / (div1 * div2) * NIN
	adct.PSC.Store(tim.PSC(div1 - 1))
	adct.ARR.Store(tim.ARR(div2 - 1))
	adct.EGR.Store(tim.UG)

	const n = 256 // number of samples per input
	buf := make([]uint16, n*NIN)

	const (
		avcc   = 5000.0 / 28.0 // A/VCC
		offset = 10
	)

	for {
		_, err := adcd.Read16(buf)
		checkErr(err)

		for i := 0; i < len(buf); i += NIN {
			for k := 0; k < NIN; k++ {
				fmt.Printf("%4d ", buf[i+k])
			}
			fmt.Printf("\r\n")
		}

		yvcc := int32(0)
		for i := INVCC; i < len(buf); i += NIN {
			yvcc += int32(buf[i])
		}
		yvcc /= n

		y0 := yvcc/2 - offset

		y6avg := int32(0)
		y6rms := uint32(0)
		for i := IN6; i < len(buf); i += NIN {
			dy := int32(buf[i]) - y0
			y6avg += dy
			y6rms += uint32(dy * dy)
		}
		y6avg /= n
		y6rms = sqrt(y6rms / n)

		scale := avcc / float64(yvcc)

		fmt.Printf(
			"yavg = %d  iavg = %.1f A  yrms = %d irms = %.1f A\r\n",
			y6avg, float64(y6avg)*scale,
			y6rms, float64(y6rms)*scale,
		)
		delay.Millisec(1e3)
	}
}

func ttsISR() {
	tts.ISR()
}

func ttsRxDMAISR() {
	tts.RxDMAISR()
}

func ttsTxDMAISR() {
	tts.TxDMAISR()
}

func adcISR() {
	adcd.ISR()
}

func adcDMAISR() {
	adcd.DMAISR()
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.USART1:        ttsISR,
	irq.DMA1_Channel4: ttsTxDMAISR,
	irq.DMA1_Channel5: ttsRxDMAISR,

	irq.ADC1_2:        adcISR,
	irq.DMA1_Channel1: adcDMAISR,
}

//// utils

func draw(w io.Writer, x uint16) {
	const s = "                                                                                                                                                                                                                                                                "
	fmt.Fprintf(w, "%-5d %s+\r\n", x, s[len(s)-int(x>>8):])
}

func checkErr(err error) {
	if err == nil {
		return
	}
	fmt.Println(err.Error())
	for {
	}
}

func sqrt(num uint32) uint32 {
	op := num
	res := uint32(0)
	one := uint32(1) << 30
	for one > op {
		one >>= 2
	}
	for one != 0 {
		if op >= res+one {
			op -= res + one
			res = (res >> 1) + one
		} else {
			res >>= 1
		}
		one >>= 2
	}
	return res
}
