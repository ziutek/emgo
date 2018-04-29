### Christmas Tree Lights ###

[![Video](https://raw.githubusercontent.com/ziutek/emgo/devel/egpath/src/stm32/examples/minidev/treelights/images/video.jpg)](http://www.youtube.com/watch?v=7qUz77a7IhU)
![Controller](https://raw.githubusercontent.com/ziutek/emgo/devel/egpath/src/stm32/examples/minidev/treelights/images/controller.jpg)
![Picture 1](https://raw.githubusercontent.com/ziutek/emgo/devel/egpath/src/stm32/examples/minidev/treelights/images/picture1.jpg)
![Picture 2](https://raw.githubusercontent.com/ziutek/emgo/devel/egpath/src/stm32/examples/minidev/treelights/images/picture2.jpg)

#### Components ####

1. String of 50 RGB LEDs, each with WS2811 controller:
	
	![WS2811 based string of 50 LEDs](https://raw.githubusercontent.com/ziutek/emgo/devel/egpath/src/stm32/examples/minidev/treelights/images/ledstring.jpg)
	
2. US-100 ultrasonic distance sensor:
	
	![US-100 ultrasonic distance sensor](https://raw.githubusercontent.com/ziutek/emgo/devel/egpath/src/stm32/examples/minidev/treelights/images/US-100.jpg)
	
3. STM32 Mini Development Board:
	
	![US-100 ultrasonic distance sensor](https://raw.githubusercontent.com/ziutek/emgo/devel/egpath/src/stm32/examples/minidev/treelights/images/minidev.jpg)
	
4. A few electronic components:
	
- 8 - 16 &#8486; speaker (I used one disassembled from the old Compaq PC),
- 2 x BD135 and 2 x BD136 transistors to drive the speaker (can be BD139/BD140 or other medium power complementary NPN/PNP pair), 
- SN74HCT04N (6 inverters in one case). One inverter used by UART based WS281x driver (it performs required signal inversion and expansion to 0 - 5 V). Two used to drive speaker transistors.

#### Schematic ####

![Schematic](https://raw.githubusercontent.com/ziutek/emgo/devel/egpath/src/stm32/examples/minidev/treelights/images/schematic.png)

####  Description #####

##### Lights #####

This project uses STM32 UART peripheral to send data to WS2811 controllers. One byte transmitted by UART represents three WS2811 bits (8 bytes per 24 bit RGB pixel).

The raw STM32 UART Tx signal can not be used to form valid WS2811 bitstream, because:

- the generated start and stop bits have bad polarity,
- MCUs is powered by 3.3 V linear regulator but WS281x requires at least 3.5 V for logical high signal when powered from 5 V source.

Some way is required to invert UART signal and convert it to 5 V logic. This project uses SN74HCT04N chip (hex inverters in one case). As HCT device it accepts TTL logic at input, and when powered form 5 V, generates inverted signal with correct logic levels at output.

##### Audio #####

STM32103C8T6 has no internal DAC peripheral. Instead of use external DAC the internal general purpose timer in PWM mode is used to generate sound. To avoid big output capacitor and increase power, the speaker is driven by symetrical output. The software controls two channels of STM32 timer in the way shown in the figure below: 

![PWM](https://raw.githubusercontent.com/ziutek/emgo/devel/egpath/src/stm32/examples/minidev/treelights/images/pwm.png)

According to the datasheet and the Flash Size Register the STM32F103R8T6 MCU has 64 KB of Flash. Audio code heavily exploits the "property" that 128 KB of Flash seems to be available on any STM32F103R8T6 chip.

Use ../load-oocd-128k.sh to flash this program.
