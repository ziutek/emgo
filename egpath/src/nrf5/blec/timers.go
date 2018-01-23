package blec

import (
	"nrf5/hal/rtc"
	"nrf5/hal/te"
	"nrf5/hal/timer"
)

func timerInit(t *timer.Periph) {
	t.Task(timer.STOP).Trigger()
	t.StoreMODE(timer.TIMER)
	t.StoreBITMODE(timer.Bit32)
	t.StorePRESCALER(4) // 1 MHz (allows switch HFCLK to low power mode)
	t.DisableIRQ(te.EvAll)
	
	/*irq := rtos.IRQ(t.IRQ())
	irq.SetPrio(rtos.IRQPrioHighest)
	irq.Enable()*/
}

func rtcInit(rt *rtc.Periph) {
	rt.Task(rtc.STOP).Trigger()
	rt.StorePRESCALER(0) // 32768 Hz
	rt.DisableIRQ(te.EvAll)
	rt.DisablePPI(te.EvAll)
	
	/*irq := rtos.IRQ(rt.IRQ())
	irq.SetPrio(rtos.IRQPrioHighest)
	irq.Enable()*/
}
