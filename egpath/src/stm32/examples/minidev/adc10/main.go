package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"rtos"
	"strconv"

	"stm32/hal/adc"
	"stm32/hal/dma"
	"stm32/hal/gpio"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
	"stm32/hal/usart"

	"stm32/hal/raw/rcc"
	"stm32/hal/raw/tim"
)

var (
	adcd *adc.Driver
	adct *tim.TIM_Periph
	tts  *usart.Driver
)

const (
	in0 = 0
	in1 = 1
	in3 = 2
	in4 = 3
	vcc = 4
	in6 = 5
	in7 = 6
	in8 = 7
	in9 = 8
	nin = 9
)

var instr = [nin]string{
	in0: "in0",
	in1: "in1",
	in3: "in3",
	in4: "in4",
	vcc: "vcc",
	in6: "in6",
	in7: "in7",
	in8: "in8",
	in9: "in9",
}

func init() {
	system.SetupPLL(8, 1, 36/8)
	systick.Setup(2e6)

	// GPIO

	var apins [nin]gpio.Pin

	gpio.A.EnableClock(true)
	apins[in0] = gpio.A.Pin(0)
	apins[in1] = gpio.A.Pin(1)
	opin := gpio.A.Pin(2) // USART2_TX for one-wire bus
	apins[in3] = gpio.A.Pin(3)
	apins[in4] = gpio.A.Pin(4)
	apins[vcc] = gpio.A.Pin(5)
	apins[in6] = gpio.A.Pin(6)
	apins[in7] = gpio.A.Pin(7)
	tx := gpio.A.Pin(9)
	rx := gpio.A.Pin(10)

	gpio.B.EnableClock(true)
	apins[in8] = gpio.B.Pin(0)
	apins[in9] = gpio.B.Pin(1)

	// DMA
	dma1 := dma.DMA1
	dma1.EnableClock(true)

	// USART

	tx.Setup(&gpio.Config{Mode: gpio.Alt})
	rx.Setup(&gpio.Config{Mode: gpio.AltIn, Pull: gpio.PullUp})
	tts = usart.NewDriver(
		usart.USART1, dma1.Channel(4, 0), dma1.Channel(5, 0), make([]byte, 80),
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


	const n = 300 // number of samples per input

	buf := make([]uint16, n*nin)
	var base [nin]int32

	for {
		skipLine()

		if _, err := adcd.Read16(buf); err != nil {
			fmt.Printf("error: %s\r\n", err.Error())
			continue
		}

		yvcc := int32(0)
		for i := vcc; i < len(buf); i += nin {
			yvcc += int32(buf[i])
		}
		yvcc = yvcc / n

		fmt.Printf("vcc: %d\r\n", yvcc)

		if err := readBase(&base); err != nil {
			fmt.Printf("error: %s\r\n", err.Error())
			continue
		}

		fmt.Printf("base: %d\r\n", base[:])

		for i := 0; i < nin; i++ {
			y0 := base[i]
			var ysum int32
			var ysum2 uint64
			for k := i; k < len(buf); k += nin {
				y := int32(buf[k]) - y0
				ysum += y
				ysum2 += uint64(y * y)
			}
			fmt.Printf(
				"%s: %d %d\r\n",
				instr[i], ysum/n, sqrt(uint32(ysum2/n)),
			)
		}
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

func nextInt32(s []byte) ([]byte, int32, error) {
	i := 0
	for i < len(s) {
		b := s[i]
		if b == ' ' || b == '\t' {
			break
		}
		i++
	}
	v, err := strconv.ParseInt32(s[:i], 0)
	return s[i:], v, err
}

func skipLine() error {
	for {
		b, err := tts.ReadByte()
		if err != nil {
			return err
		}
		if b == '\n' || b == '\r' {
			return nil
		}
	}
}

func readBase(base *[nin]int32) error {
	var buf [80]byte
	i := -1
	for {
		i++
		b, err := tts.ReadByte()
		if err != nil {
			return err
		}
		if b == '\n' || b == '\r' {
			break
		}
		if i == len(buf) {
			return errors.New("line too long")
		}
		buf[i] = b
	}
	s := buf[:i]
	for i := 0; i < nin; i++ {
		s = bytes.TrimSpace(s)
		if len(s) == 0 {
			break
		}
		var (
			v   int32
			err error
		)
		s, v, err = nextInt32(s)
		if err != nil {
			return err
		}
		base[i] = v
	}
	return nil
}

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
