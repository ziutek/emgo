package main

import (
	"unsafe"

	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/rtcst"
)

func init() {
	system.Setup(8, 1, 72/8)
	rtcst.Setup(32768)
}

func HostIO(cmd int, p unsafe.Pointer) int

func main() {

	tt := [...]byte{':', 't', 't', 0}

	type openArgs struct {
		path *byte
		mode int
		plen int
	}
	oa := openArgs{&tt[0], 4, len(tt) - 1}
	ret := HostIO(0x01, unsafe.Pointer(&oa))
	if ret == -1 {
		return
	}

	buf := [...]byte{'H', 'e', 'l', 'l', 'o', '!', '\n'}

	type writeArgs struct {
		fd   int
		data *byte
		n    int
	}
	wa := writeArgs{ret, &buf[0], len(buf)}
	ret = HostIO(0x05, unsafe.Pointer(&wa))
	if ret != 0 {
		return
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTCAlarm: rtcst.ISR,
}
