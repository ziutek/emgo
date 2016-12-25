### Christmas Tree Lights ###

![US-100 ultrasonic distance sensor](https://raw.githubusercontent.com/ziutek/emgo/devel/egpath/src/stm32/examples/minidev/treelights/images/controller.jpg

![US-100 ultrasonic distance sensor](https://raw.githubusercontent.com/ziutek/emgo/devel/egpath/src/stm32/examples/minidev/treelights/images/picture1.jpg


#### Components ####

1. String of 50 RGB LEDs, each with WS2811 controller:

![WS2811 based string of 50 LEDs](https://raw.githubusercontent.com/ziutek/emgo/devel/egpath/src/stm32/examples/minidev/treelights/images/ledstring.jpg)

2. US-100 ultrasonic distance sensor:

![US-100 ultrasonic distance sensor](https://raw.githubusercontent.com/ziutek/emgo/devel/egpath/src/stm32/examples/minidev/treelights/images/US-100.jpg)


3. STM32 Mini Development Board:

![STM32 Mini Development Board](https://raw.githubusercontent.com/ziutek/emgo/devel/egpath/src/stm32/examples/minidev/board.jpg)

4. A few electronic components:

- Some 8 - 16 Ohm speaker (I used one disassembled from the old Compaq PC),

- 2 x BD135 and 2 x BD136 transistors to drive the speaker (can be BD139/BD140 or other medium power complementary NPN/PNP pair), 

- SN74HCT04N (6 inverters in one case). One inverter used by UART based WS281x driver (it performs required signal inversion and expansion to 0 - 5 V). Two used to drive speaker transistors.

#### Schematic ####

![WS2811 based string of 50 LEDs](https://raw.githubusercontent.com/ziutek/emgo/devel/egpath/src/stm32/examples/minidev/treelights/images/schematic.jpg)

####  Description #####

![US-100 ultrasonic distance sensor](https://raw.githubusercontent.com/ziutek/emgo/devel/egpath/src/stm32/examples/minidev/treelights/images/picture2.jpg
